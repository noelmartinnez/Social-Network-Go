package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// Función para registrar un usuario
func registerUser() {
	var user User

	scanner := bufio.NewScanner(os.Stdin)

	println()
	fmt.Println("Ingrese el email:")
	if scanner.Scan() {
		user.Email = scanner.Text()
		if !isValidEmail(user.Email) {
			printMessage("Email inválido.")
			return
		}
	}

	println()
	fmt.Println("Ingrese el nombre de usuario:")
	if scanner.Scan() {
		user.Username = scanner.Text()
	}

	println()
	fmt.Println("Ingrese la contraseña:")
	if scanner.Scan() {
		user.Password = scanner.Text()
	}

	println()
	fmt.Println("Ingrese el nombre:")
	if scanner.Scan() {
		user.Nombre = scanner.Text()
	}

	println()
	fmt.Println("Ingrese los apellidos:")
	if scanner.Scan() {
		user.Apellidos = scanner.Text()
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		printMessage("Error al codificar el usuario.")
		return
	}

	resp, err := client.Post(serviceURL+"/register", "application/json", bytes.NewBuffer(userBytes))
	if err != nil {
		printMessage("Error al registrar el usuario.")
		return
	}

	defer resp.Body.Close()

	printMessage("Usuario registrado correctamente.")
}

// Función para iniciar sesión
func loginUser() {
	var user User

	scanner := bufio.NewScanner(os.Stdin)

	println()
	fmt.Println("Ingrese el nombre de usuario:")
	if scanner.Scan() {
		user.Username = scanner.Text()
	}

	println()
	fmt.Println("Ingrese la constraseña:")
	if scanner.Scan() {
		user.Password = scanner.Text()
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		printMessage("Error al codificar el usuario.")
		return
	}

	resp, err := client.Post(serviceURL+"/login", "application/json", bytes.NewBuffer(userBytes))
	if err != nil {
		printMessage("Error al iniciar sesión.")
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		printMessage("Credenciales incorrectas o error al iniciar sesión.")
		return
	}

	var responseMap map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		printMessage("Error al decodificar la respuesta del servidor.")
		return
	}

	authToken = responseMap["token"].(string)
	rol = responseMap["rol"].(string)
	currentUserID = int(responseMap["userID"].(float64))

	printMessage("Inicio de sesión correcto.")
}

// Función par ver todos los usuarios
func viewAllUsers() {
	if authToken == "" {
		printMessage("Por favor, inicie sesión primero.")
		return
	}

	req, err := http.NewRequest("GET", serviceURL+"/all-users", nil)
	if err != nil {
		printMessage("Error al crear la solicitud.")
		return
	}

	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		printMessage("Error al obtener los usuarios.")
		return
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
		return
	}

	if len(users) == 0 {
		printMessage("No hay usuarios disponibles.")
		return
	}

	println()
	fmt.Println("Usuarios disponibles:")
	for _, user := range users {
		fmt.Printf("ID: %d\nEmail: %s\nUsername: %s\nNombre: %s\nApellidos: %s\nRol ID: %d\n\n", user.ID, user.Email, user.Username, user.Nombre, user.Apellidos, user.RolID)
	}
}

// Función para ver los usuarios
func viewUsers() {
	if authToken == "" {
		fmt.Println("Por favor, inicie sesión primero.")
		return
	}

	reqURL := serviceURL + "/users"
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		fmt.Println("Error al crear la solicitud:", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error al obtener la lista de usuarios:", err)
		return
	}
	defer resp.Body.Close()

	var users []string
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		fmt.Println("Error al decodificar la lista de usuarios:", err)
		return
	}

	fmt.Println("Usuarios disponibles para enviar mensajes:")
	for i, user := range users {
		fmt.Printf("%d. %s\n", i+1, user)
	}

	fmt.Println("Seleccione el número de usuario al que desea enviar un mensaje (o presione 'q' para salir):")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "q" {
		fmt.Println("Saliendo del chat...")
		return
	}

	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(users) {
		fmt.Println("Entrada inválida. Por favor, seleccione un número de usuario válido.")
		return
	}

	recipient := users[index-1]
	viewChat(recipient)
}

func viewChat(recipient string) {
	if authToken == "" {
		fmt.Println("Por favor, inicie sesión primero.")
		return
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/message?recipient=%s", serviceURL, recipient), nil)
	if err != nil {
		fmt.Println("Error al crear la solicitud:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error al obtener el chat:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: estado HTTP", resp.StatusCode)
		var bodyBytes []byte
		if resp.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(resp.Body)
		}
		bodyString := string(bodyBytes)
		fmt.Println("Respuesta del servidor:", bodyString)
		return
	}

	var messages []Message
	err = json.NewDecoder(resp.Body).Decode(&messages)
	if err != nil {
		fmt.Println("Error al decodificar los mensajes:", err)
		return
	}

	fmt.Println("Chat con", recipient, ":")
	for _, msg := range messages {
		var prefix string
		if msg.SentByUser {
			prefix = "[Yo] "
		} else {
			prefix = "[Otro] "
		}

		msgTime := msg.SentAt.Format("2006-01-02 15:04:05")
		if msg.SentByUser {
			color.Green("%s%s: %s\n", prefix, msgTime, msg.Content)
		} else {
			color.Magenta("%s%s: %s\n", prefix, msgTime, msg.Content)
		}
	}

	fmt.Println("Ingrese el mensaje: (o salga con 'q')")
	reader := bufio.NewReader(os.Stdin)
	message, _ := reader.ReadString('\n')
	message = strings.TrimSpace(message)

	if message == "q" {
		fmt.Println("Saliendo del chat...")
		return
	}

	reqBody := struct {
		Recipient string `json:"recipient"`
		Content   string `json:"content"`
	}{
		Recipient: recipient,
		Content:   message,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println("Error al serializar el mensaje:", err)
		return
	}

	reqs, errs := http.NewRequest("POST", fmt.Sprintf("%s/message", serviceURL), bytes.NewBuffer(reqBytes))
	if errs != nil {
		fmt.Println("Error al crear la solicitud:", errs)
		return
	}

	reqs.Header.Set("Authorization", "Bearer "+authToken)
	reqs.Header.Set("Content-Type", "application/json")

	respu, err := client.Do(reqs)
	if err != nil {
		fmt.Println("Error al enviar el mensaje:", err)
		return
	}
	defer respu.Body.Close()

	if respu.StatusCode == http.StatusOK {
		fmt.Println("Mensaje enviado exitosamente.")
	} else {
		fmt.Println("Error al enviar el mensaje. Código de estado:", respu.StatusCode)
	}
}

// Función para promover a un usuario a moderador
func promoteUserToModerator() {
	if authToken == "" {
		printMessage("Por favor, inicie sesión primero.")
		return
	}

	viewAllUsers()

	fmt.Println("Ingrese el ID del usuario a promover a moderador (o presione 'q' para salir):")
	var input string
	fmt.Scanln(&input)
	if input == "q" {
		printMessage("Operación cancelada.")
		return
	}

	userID, err := strconv.Atoi(input)
	if err != nil {
		printMessage("Entrada inválida. Por favor, seleccione un número de ID válido.")
		return
	}

	reqBody, err := json.Marshal(map[string]int{"user_id": userID})
	if err != nil {
		fmt.Println("Error al codificar la solicitud:", err)
		return
	}

	req, err := http.NewRequest("POST", serviceURL+"/promoteToModerator", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Error al crear la solicitud:", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		printMessage("Error al promover al usuario.")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		printMessage("Usuario promovido a moderador correctamente.")
	} else {
		printMessage("Error al promover al usuario.")
	}
}

// Función para eliminar un usuario
func deleteUser() {
	if authToken == "" {
		printMessage("Por favor, inicie sesión primero.")
		return
	}

	validUserIDs := getValidUserIDs()

	fmt.Println("Seleccione el ID del usuario a eliminar:")
	viewAllUsers()

	var userID int
	fmt.Scanln(&userID)

	if userID == currentUserID {
		printMessage("No se puede eliminar el usuario con la sesión iniciada. Operación cancelada.")
		return
	}

	// Comprobación del ID seleccionado
	if !contains(validUserIDs, userID) {
		printMessage("ID inválido o fuera de rango. Operación cancelada.")
		return
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/user/%d", serviceURL, userID), nil)
	if err != nil {
		printMessage("Error al crear la solicitud.")
		return
	}

	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)

	if err != nil {
		printMessage("Error al eliminar el usuario.")
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		printMessage("Usuario eliminado correctamente.")
	} else {
		printMessage("Error al eliminar el usuario.")
	}
}
