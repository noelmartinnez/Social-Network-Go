package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

func printMessage(message string) {
	println()
	fmt.Println(message)
	println()
}

// Función para validar un email
func isValidEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
	return regex.MatchString(email)
}

// Función helper para comprobar si un slice contiene un valor
func contains(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// Función para obtener los IDs de los posts
func getValidPostIDs() []int {
	if authToken == "" {
		printMessage("Por favor, inicie sesión primero.")
		return nil
	}

	req, err := http.NewRequest("GET", serviceURL+"/posts", nil)
	if err != nil {
		printMessage("Error al crear la solicitud.")
		return nil
	}

	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		printMessage("Error al obtener los posts.")
		return nil
	}

	defer resp.Body.Close()

	var posts []struct {
		ID        int    `json:"id"`
		Titulo    string `json:"titulo"`
		Texto     string `json:"texto"`
		UsuarioID int    `json:"usuario_id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		printMessage("Error al decodificar la respuesta.")
		return nil
	}

	if len(posts) == 0 {
		printMessage("No hay posts disponibles.")
		return nil
	}

	var validPostIDs []int
	for _, post := range posts {
		validPostIDs = append(validPostIDs, post.ID)
	}
	return validPostIDs
}

// Función para obtener los IDs de los usuarios
func getValidUserIDs() []int {
	if authToken == "" {
		printMessage("Por favor, inicie sesión primero.")
		return nil
	}

	req, err := http.NewRequest("GET", serviceURL+"/all-users", nil)
	if err != nil {
		printMessage("Error al crear la solicitud.")
		return nil
	}

	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		printMessage("Error al obtener los usuarios.")
		return nil
	}

	defer resp.Body.Close()

	var users []struct {
		ID        int    `json:"id"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		Nombre    string `json:"nombre,omitempty"`
		Apellidos string `json:"apellidos,omitempty"`
		RolID     int    `json:"rol_id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		printMessage("Error al decodificar la respuesta.")
		return nil
	}

	if len(users) == 0 {
		printMessage("No hay usuarios disponibles.")
		return nil
	}

	var validUserIDs []int
	for _, user := range users {
		validUserIDs = append(validUserIDs, user.ID)
	}
	return validUserIDs
}

// Implementar función de cierre de sesión
func logout() {
	printMessage("Cerrando sesión...")
	authToken = "" // Limpiar el token de autenticación
	rol = ""       // Limpiar el rol
}
