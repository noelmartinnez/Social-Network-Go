package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// Crear un nuevo grupo
func createGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var group Grupo
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		http.Error(w, "Error decoding group data", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(claimsKey).(*Claims)
	if !ok {
		http.Error(w, "Error processing user claims", http.StatusInternalServerError)
		return
	}

	hashedUsername := hashValue(claims.Username + claveHash)

	var userID int
	err = db.QueryRow("SELECT id FROM USUARIO WHERE hashed_username = ?", hashedUsername).Scan(&userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}
	encryptedDescription, err := encrypt(group.Descripcion, (aesKey))
	if err != nil {
		http.Error(w, "Error encrypting group description", http.StatusInternalServerError)

		return
	}

	result, err := db.Exec("INSERT INTO GRUPO (descripcion, administrador_id) VALUES (?, ?)", encryptedDescription, userID)
	if err != nil {
		http.Error(w, "Error creating group", http.StatusInternalServerError)
		return
	}

	groupID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Error getting group ID", http.StatusInternalServerError)
		return
	}

	// Actualizar el usuario para ser administrador del nuevo grupo
	_, err = db.Exec("UPDATE USUARIO SET grupo_id = ? WHERE id = ?", groupID, userID)
	if err != nil {
		http.Error(w, "Error updating user group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Group created successfully")
}

// Obtener todos los grupos
func getGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	query := `
		SELECT g.id, g.descripcion, g.administrador_id, COUNT(u.id) as num_users
		FROM GRUPO g
		LEFT JOIN USUARIO u ON g.id = u.grupo_id
		GROUP BY g.id, g.descripcion, g.administrador_id
	`
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Error querying groups", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var groups []struct {
		ID              int    `json:"id"`
		Descripcion     string `json:"descripcion"`
		AdministradorID int    `json:"administrador_id"`
		NumUsers        int    `json:"num_users"`
	}

	for rows.Next() {
		var group struct {
			ID              int    `json:"id"`
			Descripcion     string `json:"descripcion"`
			AdministradorID int    `json:"administrador_id"`
			NumUsers        int    `json:"num_users"`
		}
		descripcionDesencrypted, err2 := decrypt(group.Descripcion, aesKey)
		if err2 != nil {
			http.Error(w, "Error decrypting group description", http.StatusInternalServerError)
			return
		}
		group.Descripcion = string(descripcionDesencrypted)

		err := rows.Scan(&group.ID, &group.Descripcion, &group.AdministradorID, &group.NumUsers)
		if err != nil {
			http.Error(w, "Error scanning groups", http.StatusInternalServerError)
			return
		}
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating through groups", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

// Unirse a un grupo
func joinGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqData struct {
		GroupID int `json:"group_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(claimsKey).(*Claims)
	if !ok {
		http.Error(w, "Error processing user claims", http.StatusInternalServerError)
		return
	}

	hashedUsername := hashValue(claims.Username + claveHash)

	var userID int
	err := db.QueryRow("SELECT id FROM USUARIO WHERE hashed_username = ?", hashedUsername).Scan(&userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	// Verificar si el grupo tiene menos de 5 miembros
	var memberCount int
	err = db.QueryRow("SELECT COUNT(*) FROM USUARIO WHERE grupo_id = ?", reqData.GroupID).Scan(&memberCount)
	if err != nil {
		http.Error(w, "Error querying group members", http.StatusInternalServerError)
		return
	}

	if memberCount >= 5 {
		http.Error(w, "Group is full", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE USUARIO SET grupo_id = ? WHERE id = ?", reqData.GroupID, userID)
	if err != nil {
		http.Error(w, "Error joining group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Joined group successfully")
}

// Obtener los posts privados del grupo
func getGroupPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value(claimsKey).(*Claims)
	if !ok {
		http.Error(w, "Error processing user claims", http.StatusInternalServerError)
		return
	}

	hashedUsername := hashValue(claims.Username + claveHash)

	var userID int
	err := db.QueryRow("SELECT id FROM USUARIO WHERE hashed_username = ?", hashedUsername).Scan(&userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	var groupID int
	err = db.QueryRow("SELECT grupo_id FROM USUARIO WHERE id = ?", userID).Scan(&groupID)
	if err != nil {
		http.Error(w, "User is not in a group", http.StatusBadRequest)
		return
	}

	rows, err := db.Query("SELECT id, titulo, texto, usuario_id FROM POST WHERE grupo_id = ? AND es_publico = FALSE", groupID)
	if err != nil {
		http.Error(w, "Error querying posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []struct {
		ID        int    `json:"id"`
		Titulo    string `json:"titulo"`
		Texto     string `json:"texto"`
		UsuarioID int    `json:"usuario_id"`
	}

	for rows.Next() {
		var post struct {
			ID        int    `json:"id"`
			Titulo    string `json:"titulo"`
			Texto     string `json:"texto"`
			UsuarioID int    `json:"usuario_id"`
		}
		err := rows.Scan(&post.ID, &post.Titulo, &post.Texto, &post.UsuarioID)
		if err != nil {
			http.Error(w, "Error scanning posts", http.StatusInternalServerError)
			return
		}

		decryptedText, err := decrypt(post.Texto, (aesKey))
		if err != nil {
			http.Error(w, "Error decrypting post text", http.StatusInternalServerError)
			return
		}
		decryptedTitulo, err := decrypt(post.Titulo, (aesKey))
		if err != nil {
			http.Error(w, "Error decrypting post text", http.StatusInternalServerError)
			return
		}
		post.Texto = string(decryptedText)
		post.Titulo = string(decryptedTitulo)

		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating through posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

/*
func manageGroup(w http.ResponseWriter, r *http.Request) {
	groupID := strings.TrimPrefix(r.URL.Path, "/group/")

	claims, ok := r.Context().Value(claimsKey).(*Claims)
	if !ok {
		http.Error(w, "Error processing user claims", http.StatusInternalServerError)
		return
	}

	hashedUsername := hashValue(claims.Username + claveHash)

	var userID int
	err := db.QueryRow("SELECT id FROM USUARIO WHERE hashed_username = ?", hashedUsername).Scan(&userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	// Verificar si el usuario es administrador del grupo
	var adminID int
	err = db.QueryRow("SELECT administrador_id FROM GRUPO WHERE id = ?", groupID).Scan(&adminID)
	if err != nil {
		http.Error(w, "Group not found", http.StatusBadRequest)
		return
	}

	if adminID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case "DELETE":
		// Eliminar el grupo
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "Error starting transaction", http.StatusInternalServerError)
			return
		}

		if _, err := tx.Exec("DELETE FROM POST WHERE grupo_id = ?", groupID); err != nil {
			tx.Rollback()
			http.Error(w, "Error deleting group's posts", http.StatusInternalServerError)
			return
		}

		if _, err := tx.Exec("DELETE FROM USUARIO WHERE grupo_id = ?", groupID); err != nil {
			tx.Rollback()
			http.Error(w, "Error deleting group members", http.StatusInternalServerError)
			return
		}

		if _, err := tx.Exec("DELETE FROM GRUPO WHERE id = ?", groupID); err != nil {
			tx.Rollback()
			http.Error(w, "Error deleting group", http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(); err != nil {
			http.Error(w, "Error committing transaction", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Group deleted successfully")

	case "PUT":
		// Editar la descripción del grupo
		var reqData struct {
			Descripcion string `json:"descripcion"`
		}
		if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("UPDATE GRUPO SET descripcion = ? WHERE id = ?", reqData.Descripcion, groupID)
		if err != nil {
			http.Error(w, "Error updating group description", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Group description updated successfully")

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
*/

// Obtener la información del grupo del usuario
func userGroupInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value(claimsKey).(*Claims)
	if !ok {
		http.Error(w, "Error processing user claims", http.StatusInternalServerError)
		return
	}

	hashedUsername := hashValue(claims.Username + claveHash)

	var userID int
	var groupID sql.NullInt32
	err := db.QueryRow("SELECT id, grupo_id FROM USUARIO WHERE hashed_username = ?", hashedUsername).Scan(&userID, &groupID)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	var adminID int
	var isAdmin bool
	if groupID.Valid {
		err = db.QueryRow("SELECT administrador_id FROM GRUPO WHERE id = ?", groupID.Int32).Scan(&adminID)
		if err != nil {
			http.Error(w, "Group not found", http.StatusBadRequest)
			return
		}
		isAdmin = (userID == adminID)
	} else {
		isAdmin = false
	}

	groupInfo := struct {
		GroupID int  `json:"group_id"`
		IsAdmin bool `json:"is_admin"`
	}{
		GroupID: int(groupID.Int32),
		IsAdmin: isAdmin,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groupInfo)
}

// Función para abandonar un grupo
func leaveGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value(claimsKey).(*Claims)
	if !ok {
		http.Error(w, "Error processing user claims", http.StatusInternalServerError)
		return
	}

	hashedUsername := hashValue(claims.Username + claveHash)

	var userID int
	var groupID sql.NullInt32
	err := db.QueryRow("SELECT id, grupo_id FROM USUARIO WHERE hashed_username = ?", hashedUsername).Scan(&userID, &groupID)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	if !groupID.Valid {
		http.Error(w, "User is not in a group", http.StatusBadRequest)
		return
	}

	var adminID int
	err = db.QueryRow("SELECT administrador_id FROM GRUPO WHERE id = ?", groupID).Scan(&adminID)
	if err != nil {
		http.Error(w, "Group not found", http.StatusBadRequest)
		return
	}

	if adminID == userID {
		http.Error(w, "Admin cannot leave the group. Delete the group or transfer admin rights.", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE USUARIO SET grupo_id = NULL WHERE id = ?", userID)
	if err != nil {
		http.Error(w, "Error leaving group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Left group successfully")
}

// Función para obtener los miembros de un grupo
func getGroupMembers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value(claimsKey).(*Claims)
	if !ok {
		http.Error(w, "Error processing user claims", http.StatusInternalServerError)
		return
	}

	hashedUsername := hashValue(claims.Username + claveHash)

	var groupID sql.NullInt32
	err := db.QueryRow("SELECT grupo_id FROM USUARIO WHERE hashed_username = ?", hashedUsername).Scan(&groupID)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	if !groupID.Valid {
		http.Error(w, "User is not in a group", http.StatusBadRequest)
		return
	}

	rows, err := db.Query("SELECT id, username FROM USUARIO WHERE grupo_id = ? ORDER BY id ASC", groupID)
	if err != nil {
		http.Error(w, "Error querying group members", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var members []struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
	}
	for rows.Next() {
		var member struct {
			ID       int    `json:"id"`
			Username string `json:"username"`
		}

		var encryptedUsername string
		err := rows.Scan(&member.ID, &encryptedUsername)

		if err != nil {
			http.Error(w, "Error scanning members", http.StatusInternalServerError)
			return
		}

		decryptedUsername, err := decrypt(encryptedUsername, (aesKey))

		if err != nil {
			http.Error(w, "Error decrypting username", http.StatusInternalServerError)
			return
		}

		member.Username = string(decryptedUsername)
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating through members", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

// Función para eliminar un miembro del grupo
func removeGroupMember(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value(claimsKey).(*Claims)
	if !ok {
		http.Error(w, "Error processing user claims", http.StatusInternalServerError)
		return
	}

	hashedUsername := hashValue(claims.Username + claveHash)

	var userID int
	var groupID sql.NullInt32
	err := db.QueryRow("SELECT id, grupo_id FROM USUARIO WHERE hashed_username = ?", hashedUsername).Scan(&userID, &groupID)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	if !groupID.Valid {
		http.Error(w, "User is not in a group", http.StatusBadRequest)
		return
	}

	var adminID int
	err = db.QueryRow("SELECT administrador_id FROM GRUPO WHERE id = ?", groupID).Scan(&adminID)
	if err != nil {
		http.Error(w, "Group not found", http.StatusBadRequest)
		return
	}

	if adminID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var reqData struct {
		MemberID int `json:"member_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	if reqData.MemberID == adminID {
		http.Error(w, "Cannot remove the group administrator", http.StatusBadRequest)
		return
	}

	var memberExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM USUARIO WHERE id = ? AND grupo_id = ?)", reqData.MemberID, groupID).Scan(&memberExists)
	if err != nil {
		http.Error(w, "Error verifying member", http.StatusInternalServerError)
		return
	}

	if !memberExists {
		http.Error(w, "Member not found in the group", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE USUARIO SET grupo_id = NULL WHERE id = ?", reqData.MemberID)
	if err != nil {
		http.Error(w, "Error removing group member", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Group member removed successfully")
}

// Función para editar la descripción del grupo
func editGroupDescription(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Only PUT method is allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value(claimsKey).(*Claims)
	if !ok {
		http.Error(w, "Error processing user claims", http.StatusInternalServerError)
		return
	}

	hashedUsername := hashValue(claims.Username + claveHash)

	var userID int
	var groupID sql.NullInt32
	err := db.QueryRow("SELECT id, grupo_id FROM USUARIO WHERE hashed_username = ?", hashedUsername).Scan(&userID, &groupID)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	if !groupID.Valid {
		http.Error(w, "User is not in a group", http.StatusBadRequest)
		return
	}

	var adminID int
	err = db.QueryRow("SELECT administrador_id FROM GRUPO WHERE id = ?", groupID).Scan(&adminID)
	if err != nil {
		http.Error(w, "Group not found", http.StatusBadRequest)
		return
	}

	if adminID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var reqData struct {
		Descripcion string `json:"descripcion"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}
	encryptedDescription, err := encrypt(reqData.Descripcion, (aesKey))
	if err != nil {
		http.Error(w, "Error encrypting group description", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("UPDATE GRUPO SET descripcion = ? WHERE id = ?", encryptedDescription, groupID)
	if err != nil {
		http.Error(w, "Error updating group description", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Group description updated successfully")
}

// Función para eliminar un grupo
func deleteGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Only DELETE method is allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value(claimsKey).(*Claims)
	if !ok {
		http.Error(w, "Error processing user claims", http.StatusInternalServerError)
		return
	}

	hashedUsername := hashValue(claims.Username + claveHash)

	var userID int
	var groupID sql.NullInt32
	err := db.QueryRow("SELECT id, grupo_id FROM USUARIO WHERE hashed_username = ?", hashedUsername).Scan(&userID, &groupID)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	if !groupID.Valid {
		http.Error(w, "User is not in a group", http.StatusBadRequest)
		return
	}

	var adminID int
	err = db.QueryRow("SELECT administrador_id FROM GRUPO WHERE id = ?", groupID).Scan(&adminID)
	if err != nil {
		http.Error(w, "Group not found", http.StatusBadRequest)
		return
	}

	if adminID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Error starting transaction", http.StatusInternalServerError)
		return
	}

	if _, err := tx.Exec("DELETE FROM POST WHERE grupo_id = ?", groupID); err != nil {
		tx.Rollback()
		http.Error(w, "Error deleting group's posts", http.StatusInternalServerError)
		return
	}

	if _, err := tx.Exec("UPDATE USUARIO SET grupo_id = NULL WHERE grupo_id = ?", groupID); err != nil {
		tx.Rollback()
		http.Error(w, "Error updating group members", http.StatusInternalServerError)
		return
	}

	if _, err := tx.Exec("DELETE FROM GRUPO WHERE id = ?", groupID); err != nil {
		tx.Rollback()
		http.Error(w, "Error deleting group", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Group deleted successfully")
}
