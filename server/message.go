package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Manejar los mensajes
func handleMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		claims, ok := r.Context().Value(claimsKey).(*Claims)
		if !ok {
			http.Error(w, "Error processing user claims", http.StatusInternalServerError)
			fmt.Println("Error processing user claims")
			logEvent(0, "message_retrieval_failed", "Error processing user claims", r.RemoteAddr, nil)
			return
		}

		var userID int
		err := db.QueryRow("SELECT id FROM USUARIO WHERE hashed_username = ?", hashValue(claims.Username+claveHash)).Scan(&userID)
		if err != nil {
			http.Error(w, "User not found", http.StatusBadRequest)
			fmt.Println("User not found:", claims.Username)
			logEvent(0, "message_retrieval_failed", "User not found", r.RemoteAddr, map[string]interface{}{"username": claims.Username})
			return
		}
		var recipientUsername = r.URL.Query().Get("recipient")
		var userIDOtro int
		err2 := db.QueryRow("SELECT id FROM USUARIO WHERE hashed_username = ?", hashValue(recipientUsername+claveHash)).Scan(&userIDOtro)
		if err2 != nil {
			http.Error(w, "Recipient not found", http.StatusBadRequest)
			fmt.Println("Recipient not found:", recipientUsername)
			logEvent(userID, "message_retrieval_failed", "Recipient not found", r.RemoteAddr, map[string]interface{}{"recipient": recipientUsername})
			return
		}
		rows, err := db.Query("SELECT id, remitente_id, destinatario_id, contenido, DATE_FORMAT(fecha_envio, '%Y-%m-%d %H:%i:%s') AS fecha_envio, remitente_id = ? AS sent_by_user FROM MENSAJE WHERE (remitente_id = ? AND destinatario_id = ?) OR (remitente_id = ? AND destinatario_id = ?) ORDER BY fecha_envio ASC", userID, userID, userIDOtro, userIDOtro, userID)
		if err != nil {
			http.Error(w, "Error querying messages", http.StatusInternalServerError)
			fmt.Println("Error querying messages")
			logEvent(userID, "message_retrieval_failed", "Error querying messages", r.RemoteAddr, nil)
			return
		}
		defer rows.Close()
		var messages []Message
		for rows.Next() {
			var message Message
			var fechaString string
			var sentByUser int
			err := rows.Scan(&message.Id, &message.Remitente_id, &message.Destinatario_id, &message.Contenido, &fechaString, &sentByUser)
			if err != nil {
				http.Error(w, "Error scanning messages", http.StatusInternalServerError)
				fmt.Println("Error scanning messages")
				logEvent(userID, "message_retrieval_failed", "Error scanning messages", r.RemoteAddr, nil)
				return
			}
			message.SentByUser = sentByUser == 1
			fecha, err := time.Parse("2006-01-02 15:04:05", fechaString)
			if err != nil {
				http.Error(w, "Error parsing date", http.StatusInternalServerError)
				fmt.Println("Error parsing date")
				logEvent(userID, "message_retrieval_failed", "Error parsing date", r.RemoteAddr, nil)
				return
			}
			message.Fecha_envio = fecha

			if message.Contenido != "" {
				decryptedContent, err := decrypt(message.Contenido, aesKey)
				if err != nil {
					http.Error(w, "Error decrypting message", http.StatusInternalServerError)
					fmt.Println("Error decrypting message")
					logEvent(userID, "message_retrieval_failed", "Error decrypting message", r.RemoteAddr, nil)
					return
				}
				message.Contenido = string(decryptedContent)
			}

			messages = append(messages, message)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Error iterating through messages", http.StatusInternalServerError)
			fmt.Println("Error iterating through messages")
			logEvent(userID, "message_retrieval_failed", "Error iterating through messages", r.RemoteAddr, nil)
			return
		}

		logEvent(userID, "message_retrieval_success", "Messages retrieved successfully", r.RemoteAddr, map[string]interface{}{"recipient": recipientUsername})

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(messages)
		if err != nil {
			http.Error(w, "Error encoding messages to JSON", http.StatusInternalServerError)
			fmt.Println("Error encoding messages to JSON")
			logEvent(userID, "message_retrieval_failed", "Error encoding messages to JSON", r.RemoteAddr, nil)
		}
	} else if r.Method == "POST" {
		claims, ok := r.Context().Value(claimsKey).(*Claims)
		if !ok {
			http.Error(w, "Error processing user claims", http.StatusInternalServerError)
			fmt.Println("Error processing user claims")
			logEvent(0, "message_sending_failed", "Error processing user claims", r.RemoteAddr, nil)
			return
		}
		var newMessage struct {
			Recipient string `json:"recipient"`
			Content   string `json:"content"`
		}
		err := json.NewDecoder(r.Body).Decode(&newMessage)
		if err != nil {
			http.Error(w, "Error decoding message data", http.StatusBadRequest)
			fmt.Println("Error decoding message data")
			logEvent(0, "message_sending_failed", "Error decoding message data", r.RemoteAddr, nil)
			return
		}
		var senderID int
		err = db.QueryRow("SELECT id FROM USUARIO WHERE hashed_username = ?", hashValue(claims.Username+claveHash)).Scan(&senderID)
		if err != nil {
			http.Error(w, "User not found", http.StatusBadRequest)
			fmt.Println("User not found:", claims.Username)
			logEvent(0, "message_sending_failed", "User not found", r.RemoteAddr, map[string]interface{}{"username": claims.Username})
			return
		}
		var recipientID int
		err = db.QueryRow("SELECT id FROM USUARIO WHERE hashed_username = ?", hashValue(newMessage.Recipient+claveHash)).Scan(&recipientID)
		if err != nil {
			http.Error(w, "Recipient not found", http.StatusBadRequest)
			fmt.Println("Recipient not found:", newMessage.Recipient)
			logEvent(senderID, "message_sending_failed", "Recipient not found", r.RemoteAddr, map[string]interface{}{"recipient": newMessage.Recipient})
			return
		}
		encryptedContent, err := encrypt(newMessage.Content, aesKey)
		if err != nil {
			http.Error(w, "Error encrypting message", http.StatusInternalServerError)
			fmt.Println("Error encrypting message")
			logEvent(senderID, "message_sending_failed", "Error encrypting message", r.RemoteAddr, nil)
			return
		}
		_, err = db.Exec("INSERT INTO MENSAJE (remitente_id, destinatario_id, contenido) VALUES (?, ?, ?)", senderID, recipientID, encryptedContent)
		if err != nil {
			http.Error(w, "Error inserting message into database", http.StatusInternalServerError)
			fmt.Println("Error inserting message into database")
			logEvent(senderID, "message_sending_failed", "Error inserting message into database", r.RemoteAddr, nil)
			return
		}

		logEvent(senderID, "message_sent", "Message sent successfully", r.RemoteAddr, map[string]interface{}{"recipient": newMessage.Recipient})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Message sent successfully"})
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		fmt.Println("Method not allowed")
	}
}
