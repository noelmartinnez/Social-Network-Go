package server

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/argon2"
)

// Función que registra un nuevo usuario
func registerUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		logEvent(0, "register_user_failed", "Only POST method is allowed", r.RemoteAddr, nil)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Error decoding user data", http.StatusBadRequest)
		logEvent(0, "register_user_failed", "Error decoding user data", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	var adminExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM USUARIO U INNER JOIN ROLES R ON U.rol_id = R.id WHERE R.nombre = 'administrador')").Scan(&adminExists)
	if err != nil {
		http.Error(w, "Error checking for admin existence", http.StatusInternalServerError)
		logEvent(0, "register_user_failed", "Error checking for admin existence", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	var rolID int
	if !adminExists {
		err = db.QueryRow("SELECT id FROM ROLES WHERE nombre = 'administrador'").Scan(&rolID)
	} else {
		err = db.QueryRow("SELECT id FROM ROLES WHERE nombre = 'usuario'").Scan(&rolID)
	}

	if err != nil {
		http.Error(w, "Error retrieving role ID", http.StatusInternalServerError)
		logEvent(0, "register_user_failed", "Error retrieving role ID", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	salt := make([]byte, 16)
	if _, err = rand.Read(salt); err != nil {
		http.Error(w, "Error generating salt", http.StatusInternalServerError)
		logEvent(0, "register_user_failed", "Error generating salt", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	hash := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)
	hashedPassword := fmt.Sprintf("%x:%x", salt, hash)

	encryptedEmail, err := encrypt(user.Email, aesKey)
	if err != nil {
		http.Error(w, "Error encrypting email", http.StatusInternalServerError)
		logEvent(0, "register_user_failed", "Error encrypting email", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	encryptedUsername, err := encrypt(user.Username, aesKey)
	if err != nil {
		http.Error(w, "Error encrypting username", http.StatusInternalServerError)
		logEvent(0, "register_user_failed", "Error encrypting username", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	encryptedNombre, err := encrypt(user.Nombre, aesKey)
	if err != nil {
		http.Error(w, "Error encrypting nombre", http.StatusInternalServerError)
		logEvent(0, "register_user_failed", "Error encrypting nombre", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	encryptedApellidos, err := encrypt(user.Apellidos, aesKey)
	if err != nil {
		http.Error(w, "Error encrypting apellidos", http.StatusInternalServerError)
		logEvent(0, "register_user_failed", "Error encrypting apellidos", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	hashedUsername := hashValue(user.Username + claveHash)
	fmt.Println(hashedUsername)

	_, err = db.Exec("INSERT INTO USUARIO (email, username, hashed_username, password, nombre, apellidos, rol_id) VALUES (?, ?, ?, ?, ?, ?, ?)",
		encryptedEmail, encryptedUsername, hashedUsername, hashedPassword, encryptedNombre, encryptedApellidos, rolID)
	if err != nil {
		http.Error(w, "Error inserting user into database", http.StatusInternalServerError)
		logEvent(0, "register_user_failed", "Error inserting user into database", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "User registered successfully")
	logEvent(0, "register_user_success", "User registered successfully", r.RemoteAddr, map[string]interface{}{"username": user.Username})
}

// Iniciar sesión de usuario
func loginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Error decoding user data", http.StatusBadRequest)
		return
	}

	// Hashear el username para la búsqueda
	hashedUsername := hashValue(user.Username + claveHash)
	fmt.Println(hashedUsername)

	var storedHash, userRole, encryptedEmail, encryptedNombre, encryptedApellidos, encryptedUsername string
	var userID int
	err = db.QueryRow("SELECT U.id, U.password, R.nombre, U.email, U.username, U.nombre, U.apellidos FROM USUARIO U INNER JOIN ROLES R ON U.rol_id = R.id WHERE U.hashed_username = ?", hashedUsername).Scan(&userID, &storedHash, &userRole, &encryptedEmail, &encryptedUsername, &encryptedNombre, &encryptedApellidos)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		logEvent(0, "login_failed", "User not found", r.RemoteAddr, map[string]interface{}{"username": user.Username})
		return
	}

	// Separar la sal y el hash almacenados
	parts := strings.Split(storedHash, ":")
	if len(parts) != 2 {
		http.Error(w, "Stored password format is incorrect", http.StatusInternalServerError)
		logEvent(userID, "login_failed", "Stored password format is incorrect", r.RemoteAddr, nil)
		return
	}

	salt, err := hex.DecodeString(parts[0])
	if err != nil {
		http.Error(w, "Error decoding salt", http.StatusInternalServerError)
		logEvent(userID, "login_failed", "Error decoding salt", r.RemoteAddr, nil)
		return
	}

	storedPasswordHash, err := hex.DecodeString(parts[1])
	if err != nil {
		http.Error(w, "Error decoding hash", http.StatusInternalServerError)
		logEvent(userID, "login_failed", "Error decoding hash", r.RemoteAddr, nil)
		return
	}

	// Generar el hash de la contraseña proporcionada usando la misma sal
	testHash := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)

	// Comparar los hashes
	if !bytes.Equal(storedPasswordHash, testHash) {
		http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
		logEvent(userID, "login_failed", "Invalid login credentials", r.RemoteAddr, nil)
		return
	}

	// Descifrar los datos del usuario
	decryptedEmail, err := decrypt(encryptedEmail, (aesKey))
	if err != nil {
		http.Error(w, "Error decrypting email", http.StatusInternalServerError)
		logEvent(userID, "login_failed", "Error decrypting email", r.RemoteAddr, nil)
		return
	}

	decryptedNombre, err := decrypt(encryptedNombre, (aesKey))
	if err != nil {
		http.Error(w, "Error decrypting nombre", http.StatusInternalServerError)
		logEvent(userID, "login_failed", "Error decrypting nombre", r.RemoteAddr, nil)
		return
	}

	decryptedApellidos, err := decrypt(encryptedApellidos, (aesKey))
	if err != nil {
		http.Error(w, "Error decrypting apellidos", http.StatusInternalServerError)
		logEvent(userID, "login_failed", "Error decrypting apellidos", r.RemoteAddr, nil)
		return
	}

	// Crear el JWT
	claims := &Claims{
		Username: user.Username,
		Rol:      userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Error creating the token", http.StatusInternalServerError)
		logEvent(userID, "login_failed", "Error creating the token", r.RemoteAddr, nil)
		return
	}

	logEvent(userID, "login_success", "User logged in successfully", r.RemoteAddr, nil)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":     tokenString,
		"rol":       userRole,
		"userID":    userID,
		"email":     string(decryptedEmail),
		"nombre":    string(decryptedNombre),
		"apellidos": string(decryptedApellidos),
	})
}

// Obtener la lista de usuarios disponibles
func getUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("estoy aqui sabes")
	if r.Method != "GET" {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		logEvent(0, "get_users_failed", "Only GET method is allowed", r.RemoteAddr, nil)
		return
	}

	claims, ok := r.Context().Value(claimsKey).(*Claims)
	if !ok {
		http.Error(w, "Error processing user claims", http.StatusInternalServerError)
		logEvent(0, "get_users_failed", "Error processing user claims", r.RemoteAddr, nil)
		return
	}
	senderUsername := claims.Username
	fmt.Print(senderUsername + " soy yo")

	hashedSenderUsername := hashValue(senderUsername + claveHash)

	rows, err := db.Query("SELECT username FROM USUARIO WHERE hashed_username != ?", hashedSenderUsername)
	if err != nil {
		http.Error(w, "Error querying users", http.StatusInternalServerError)
		logEvent(0, "get_users_failed", "Error querying users", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var encryptedUsername string
		err := rows.Scan(&encryptedUsername)
		if err != nil {
			http.Error(w, "Error scanning users", http.StatusInternalServerError)
			logEvent(0, "get_users_failed", "Error scanning users", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
			return
		}

		decryptedUsername, err := decrypt(encryptedUsername, (aesKey))
		if err != nil {
			http.Error(w, "Error decrypting username", http.StatusInternalServerError)
			logEvent(0, "get_users_failed", "Error decrypting username", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
			return
		}

		users = append(users, string(decryptedUsername))
	}
	fmt.Println(users)
	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating through users", http.StatusInternalServerError)
		logEvent(0, "get_users_failed", "Error iterating through users", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, "Error encoding users to JSON", http.StatusInternalServerError)
		logEvent(0, "get_users_failed", "Error encoding users to JSON", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	logEvent(0, "get_users_success", "Users retrieved successfully", r.RemoteAddr, nil)
}

// Promover a un usuario a moderador
func promoteToModerator(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		logEvent(0, "promote_to_moderator_failed", "Only POST method is allowed", r.RemoteAddr, nil)
		return
	}

	claims, ok := r.Context().Value(claimsKey).(*Claims)
	if !ok || claims.Rol != "administrador" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		logEvent(0, "promote_to_moderator_failed", "Unauthorized access", r.RemoteAddr, map[string]interface{}{"role": claims.Rol})
		return
	}

	var reqData struct {
		UserID int `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		logEvent(0, "promote_to_moderator_failed", "Error decoding request body", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	var currentRol string
	err := db.QueryRow("SELECT R.nombre FROM USUARIO U INNER JOIN ROLES R ON U.rol_id = R.id WHERE U.id = ?", reqData.UserID).Scan(&currentRol)
	if err != nil {
		http.Error(w, "Error retrieving user role", http.StatusInternalServerError)
		logEvent(0, "promote_to_moderator_failed", "Error retrieving user role", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	if currentRol != "usuario" {
		http.Error(w, "Only normal users can be promoted", http.StatusBadRequest)
		logEvent(0, "promote_to_moderator_failed", "Only normal users can be promoted", r.RemoteAddr, map[string]interface{}{"current_role": currentRol})
		return
	}

	_, err = db.Exec("UPDATE USUARIO SET rol_id = (SELECT id FROM ROLES WHERE nombre = 'moderador') WHERE id = ?", reqData.UserID)
	if err != nil {
		http.Error(w, "Error updating user role", http.StatusInternalServerError)
		logEvent(0, "promote_to_moderator_failed", "Error updating user role", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "User promoted to moderator successfully")
	logEvent(0, "promote_to_moderator_success", "User promoted to moderator successfully", r.RemoteAddr, map[string]interface{}{"promoted_user_id": reqData.UserID})
}

// Obtener la lista de todos los usuarios
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		logEvent(0, "get_all_users_failed", "Only GET method is allowed", r.RemoteAddr, nil)
		return
	}

	rows, err := db.Query("SELECT id, email, username, nombre, apellidos, rol_id FROM USUARIO")
	if err != nil {
		http.Error(w, "Error querying users", http.StatusInternalServerError)
		logEvent(0, "get_all_users_failed", "Error querying users", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []struct {
		ID        int    `json:"id"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		Nombre    string `json:"nombre,omitempty"`
		Apellidos string `json:"apellidos,omitempty"`
		RolID     int    `json:"rol_id"`
	}

	for rows.Next() {
		var user struct {
			ID        int    `json:"id"`
			Email     string `json:"email"`
			Username  string `json:"username"`
			Nombre    string `json:"nombre,omitempty"`
			Apellidos string `json:"apellidos,omitempty"`
			RolID     int    `json:"rol_id"`
		}
		err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.Nombre, &user.Apellidos, &user.RolID)
		if err != nil {
			http.Error(w, "Error scanning users", http.StatusInternalServerError)
			logEvent(0, "get_all_users_failed", "Error scanning users", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
			return
		}

		decryptedEmail, err := decrypt(user.Email, aesKey)
		if err != nil {
			http.Error(w, "Error decrypting email", http.StatusInternalServerError)
			logEvent(user.ID, "get_all_users_failed", "Error decrypting email", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
			return
		}
		user.Email = string(decryptedEmail)

		decryptedUsername, err := decrypt(user.Username, aesKey)
		if err != nil {
			http.Error(w, "Error decrypting username", http.StatusInternalServerError)
			logEvent(user.ID, "get_all_users_failed", "Error decrypting username", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
			return
		}
		user.Username = string(decryptedUsername)
		decryptedNombre, err := decrypt(user.Nombre, aesKey)
		if err != nil {
			http.Error(w, "Error decrypting name", http.StatusInternalServerError)
			logEvent(user.ID, "get_all_users_failed", "Error decrypting name", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
			return
		}
		user.Nombre = string(decryptedNombre)
		decryptedApellidos, err := decrypt(user.Apellidos, aesKey)
		if err != nil {
			http.Error(w, "Error decrypting last name", http.StatusInternalServerError)
			logEvent(user.ID, "get_all_users_failed", "Error decrypting last name", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
			return
		}
		user.Apellidos = string(decryptedApellidos)
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating through users", http.StatusInternalServerError)
		logEvent(0, "get_all_users_failed", "Error iterating through users", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
	logEvent(0, "get_all_users_success", "Users retrieved successfully", r.RemoteAddr, nil)
}

// Obtener la lista de roles disponibles
func getRoles(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		logEvent(0, "get_roles_failed", "Only GET method is allowed", r.RemoteAddr, nil)
		return
	}

	rows, err := db.Query("SELECT id, nombre FROM ROLES")
	if err != nil {
		http.Error(w, "Error querying roles", http.StatusInternalServerError)
		logEvent(0, "get_roles_failed", "Error querying roles", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	defer rows.Close()

	var roles []Rol

	for rows.Next() {
		var rol Rol
		err := rows.Scan(&rol.ID, &rol.Nombre)

		if err != nil {
			http.Error(w, "Error scanning roles", http.StatusInternalServerError)
			logEvent(0, "get_roles_failed", "Error scanning roles", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
			return
		}

		roles = append(roles, rol)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating through roles", http.StatusInternalServerError)
		logEvent(0, "get_roles_failed", "Error iterating through roles", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
	logEvent(0, "get_roles_success", "Roles retrieved successfully", r.RemoteAddr, nil)
}

// Eliminar un usuario
func deleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Only DELETE method is allowed", http.StatusMethodNotAllowed)
		logEvent(0, "delete_user_failed", "Only DELETE method is allowed", r.RemoteAddr, nil)
		return
	}

	claims, ok := r.Context().Value(claimsKey).(*Claims)
	if !ok || claims.Rol != "administrador" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		logEvent(0, "delete_user_failed", "Unauthorized access attempt", r.RemoteAddr, nil)
		return
	}

	userID := strings.TrimPrefix(r.URL.Path, "/user/")

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Error starting transaction", http.StatusInternalServerError)
		logEvent(0, "delete_user_failed", "Error starting transaction", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	rows, err := db.Query("SELECT id FROM GRUPO WHERE administrador_id = ?", userID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error querying user's groups", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var groupID int
		if err := rows.Scan(&groupID); err != nil {
			tx.Rollback()
			http.Error(w, "Error scanning user's groups", http.StatusInternalServerError)
			return
		}

		// Eliminar el grupo administrado
		if _, err := tx.Exec("DELETE FROM POST WHERE grupo_id = ?", groupID); err != nil {
			tx.Rollback()
			http.Error(w, "Error deleting group's posts", http.StatusInternalServerError)
			return
		}

		if _, err := tx.Exec("UPDATE USUARIO SET grupo_id = NULL WHERE grupo_id = ?", groupID); err != nil {
			tx.Rollback()
			http.Error(w, "Error updating users' group_id", http.StatusInternalServerError)
			return
		}

		if _, err := tx.Exec("DELETE FROM GRUPO WHERE id = ?", groupID); err != nil {
			tx.Rollback()
			http.Error(w, "Error deleting group", http.StatusInternalServerError)
			return
		}
	}

	if _, err := tx.Exec("DELETE FROM POST WHERE usuario_id = ?", userID); err != nil {
		tx.Rollback()
		http.Error(w, "Error deleting user's posts", http.StatusInternalServerError)
		logEvent(0, "delete_user_failed", "Error deleting user's posts", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	if _, err := tx.Exec("DELETE FROM MENSAJE WHERE remitente_id = ? OR destinatario_id = ?", userID, userID); err != nil {
		tx.Rollback()
		http.Error(w, "Error deleting user's messages", http.StatusInternalServerError)
		logEvent(0, "delete_user_failed", "Error deleting user's messages", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	if _, err := tx.Exec("DELETE FROM USUARIO WHERE id = ?", userID); err != nil {
		tx.Rollback()
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		logEvent(0, "delete_user_failed", "Error deleting user", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		logEvent(0, "delete_user_failed", "Error committing transaction", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "User and user's posts deleted successfully")
	logEvent(0, "delete_user_success", "User and user's posts deleted successfully", r.RemoteAddr, map[string]interface{}{"user_id": userID})
}
