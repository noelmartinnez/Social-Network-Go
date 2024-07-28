package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Crear un nuevo post
func createPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		logEvent(0, "create_post_failed", "Only POST method is allowed", r.RemoteAddr, nil)
		return
	}

	var post struct {
		Post
		EsPublico bool `json:"es_publico"`
	}

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Error decoding post data", http.StatusBadRequest)
		logEvent(0, "create_post_failed", "Error decoding post data", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	claims, ok := r.Context().Value(claimsKey).(*Claims)
	if !ok {
		http.Error(w, "Error processing user claims", http.StatusInternalServerError)
		logEvent(0, "create_post_failed", "Error processing user claims", r.RemoteAddr, nil)
		return
	}
	hashedUsername := hashValue(claims.Username + claveHash)

	var userID int
	err = db.QueryRow("SELECT id FROM USUARIO WHERE hashed_username = ?", hashedUsername).Scan(&userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		logEvent(0, "create_post_failed", "User not found", r.RemoteAddr, map[string]interface{}{"error": err.Error(), "username": claims.Username})
		return
	}

	var groupID sql.NullInt32
	if !post.EsPublico {
		err = db.QueryRow("SELECT grupo_id FROM USUARIO WHERE id = ?", userID).Scan(&groupID)
		if err != nil {
			http.Error(w, "User is not in a group", http.StatusBadRequest)
			return
		}
	}

	encryptedText, err := encrypt(post.Texto, aesKey)
	if err != nil {
		http.Error(w, "Error encrypting post text", http.StatusInternalServerError)
		logEvent(userID, "create_post_failed", "Error encrypting post text", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}
	encryptedTitulo, err := encrypt(post.Titulo, aesKey)
	if err != nil {
		http.Error(w, "Error encrypting post text", http.StatusInternalServerError)
		logEvent(userID, "create_post_failed", "Error encrypting post title", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	_, err = db.Exec("INSERT INTO POST (titulo, texto, usuario_id, grupo_id, es_publico) VALUES (?, ?, ?, ?, ?)", encryptedTitulo, encryptedText, userID, groupID, post.EsPublico)
	if err != nil {
		http.Error(w, "Error inserting post into database", http.StatusInternalServerError)
		logEvent(userID, "create_post_failed", "Error inserting post into database", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Post created successfully")
	logEvent(userID, "create_post_success", "Post created successfully", r.RemoteAddr, map[string]interface{}{"post_title": post.Titulo})
}

// Obtener todos los posts
func getPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		logEvent(0, "get_posts_failed", "Only GET method is allowed", r.RemoteAddr, nil)
		return
	}

	rows, err := db.Query("SELECT id, titulo, texto, usuario_id, es_publico, grupo_id FROM POST WHERE es_publico = TRUE OR grupo_id IS NULL")
	if err != nil {
		http.Error(w, "Error querying posts", http.StatusInternalServerError)
		logEvent(0, "get_posts_failed", "Error querying posts", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}
	defer rows.Close()

	var posts []struct {
		ID        int           `json:"id"`
		Titulo    string        `json:"titulo"`
		Texto     string        `json:"texto"`
		UsuarioID int           `json:"usuario_id"`
		EsPublico bool          `json:"es_publico"`
		GrupoID   sql.NullInt32 `json:"grupo_id"`
	}

	for rows.Next() {
		var post struct {
			ID        int           `json:"id"`
			Titulo    string        `json:"titulo"`
			Texto     string        `json:"texto"`
			UsuarioID int           `json:"usuario_id"`
			EsPublico bool          `json:"es_publico"`
			GrupoID   sql.NullInt32 `json:"grupo_id"`
		}
		err := rows.Scan(&post.ID, &post.Titulo, &post.Texto, &post.UsuarioID, &post.EsPublico, &post.GrupoID)
		if err != nil {
			http.Error(w, "Error scanning posts", http.StatusInternalServerError)
			logEvent(0, "get_posts_failed", "Error scanning posts", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
			return
		}

		decryptedText, err := decrypt(post.Texto, aesKey)
		if err != nil {
			http.Error(w, "Error decrypting post text", http.StatusInternalServerError)
			logEvent(0, "get_posts_failed", "Error decrypting post text", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
			return
		}
		decryptedTitulo, err := decrypt(post.Titulo, aesKey)
		if err != nil {
			http.Error(w, "Error decrypting post text", http.StatusInternalServerError)
			logEvent(0, "get_posts_failed", "Error decrypting post title", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
			return
		}
		post.Texto = string(decryptedText)
		post.Titulo = string(decryptedTitulo)

		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating through posts", http.StatusInternalServerError)
		logEvent(0, "get_posts_failed", "Error iterating through posts", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
	logEvent(0, "get_posts_success", "Posts retrieved successfully", r.RemoteAddr, nil)
}

// Eliminar un post
func deletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Only DELETE method is allowed", http.StatusMethodNotAllowed)
		logEvent(0, "delete_post_failed", "Only DELETE method is allowed", r.RemoteAddr, nil)
		return
	}

	claims, ok := r.Context().Value(claimsKey).(*Claims)
	if !ok || (claims.Rol != "administrador" && claims.Rol != "moderador") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		logEvent(0, "delete_post_failed", "Unauthorized access attempt", r.RemoteAddr, nil)
		return
	}

	postID := strings.TrimPrefix(r.URL.Path, "/post/")

	_, err := db.Exec("DELETE FROM POST WHERE id = ?", postID)
	if err != nil {
		http.Error(w, "Error deleting post from database", http.StatusInternalServerError)
		logEvent(0, "delete_post_failed", "Error deleting post from database", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Post deleted successfully")
	logEvent(0, "delete_post_success", "Post deleted successfully", r.RemoteAddr, map[string]interface{}{"post_id": postID})
}
