package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

// Función para crear un grupo
func createGroup() {
	if authToken == "" {
		printMessage("Por favor, inicie sesión primero.")
		return
	}

	var group Grupo
	scanner := bufio.NewScanner(os.Stdin)

	println()
	fmt.Println("Ingrese la descripción del grupo:")
	if scanner.Scan() {
		group.Descripcion = scanner.Text()
	}

	groupBytes, err := json.Marshal(group)
	if err != nil {
		printMessage("Error al codificar el grupo.")
		return
	}

	req, err := http.NewRequest("POST", serviceURL+"/createGroup", bytes.NewBuffer(groupBytes))
	if err != nil {
		fmt.Println("Error al crear la solicitud.")
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		printMessage("Error al crear el grupo.")
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		printMessage("Grupo creado correctamente.")
		authenticatedOptions()
	} else {
		printMessage("Error al crear el grupo.")
	}
}

// Función para ver los grupos
func viewGroups(allowJoin bool) {
	if authToken == "" {
		printMessage("Por favor, inicie sesión primero.")
		return
	}

	req, err := http.NewRequest("GET", serviceURL+"/groups", nil)
	if err != nil {
		printMessage("Error al crear la solicitud.")
		return
	}

	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		printMessage("Error al obtener los grupos.")
		return
	}

	defer resp.Body.Close()

	var groups []struct {
		ID              int    `json:"id"`
		Descripcion     string `json:"descripcion"`
		AdministradorID int    `json:"administrador_id"`
		NumUsers        int    `json:"num_users"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&groups); err != nil {
		printMessage("Error al decodificar la respuesta.")
		return
	}

	if len(groups) == 0 {
		printMessage("No hay grupos disponibles.")
		return
	}

	println()
	fmt.Println("Grupos disponibles:")
	for _, group := range groups {
		fmt.Printf("ID: %d\nDescripción: %s\nNúmero de usuarios: %d\n\n", group.ID, group.Descripcion, group.NumUsers)
	}

	if allowJoin {
		fmt.Println("Ingrese el ID del grupo al que desea unirse (o 'q' para salir):")
		var input string
		fmt.Scanln(&input)
		if input == "q" {
			return
		}

		groupID, err := strconv.Atoi(input)
		if err != nil {
			printMessage("ID de grupo inválido.")
			return
		}

		joinGroup(groupID)
	}
}

// Función para ver los posts del grupo
func viewGroupPosts() {
	if authToken == "" {
		printMessage("Por favor, inicie sesión primero.")
		return
	}

	req, err := http.NewRequest("GET", serviceURL+"/groupPosts", nil)
	if err != nil {
		printMessage("Error al crear la solicitud.")
		return
	}

	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		printMessage("Error al obtener los posts del grupo.")
		return
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
		return
	}

	if len(posts) == 0 {
		printMessage("No hay posts del grupo disponibles.")
		return
	}

	println()
	fmt.Println("Posts del grupo:")
	for _, post := range posts {
		fmt.Printf("ID: %d\nTítulo: %s\nTexto: %s\n\n", post.ID, post.Titulo, post.Texto)
	}
}

// Función para unirse a un grupo
func joinGroup(groupID int) {
	if authToken == "" {
		printMessage("Por favor, inicie sesión primero.")
		return
	}

	reqBody := struct {
		GroupID int `json:"group_id"`
	}{
		GroupID: groupID,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		printMessage("Error al codificar la solicitud.")
		return
	}

	req, err := http.NewRequest("POST", serviceURL+"/joinGroup", bytes.NewBuffer(reqBytes))
	if err != nil {
		fmt.Println("Error al crear la solicitud.")
		return
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		printMessage("Error al unirse al grupo.")
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		printMessage("Se unió al grupo correctamente.")
		authenticatedOptions()
	} else {
		printMessage("Error al unirse al grupo.")
	}
}

// Función para abandonar un grupo
func leaveGroup() {
	if authToken == "" {
		printMessage("Por favor, inicie sesión primero.")
		return
	}

	req, err := http.NewRequest("POST", serviceURL+"/leaveGroup", nil)
	if err != nil {
		printMessage("Error al crear la solicitud.")
		return
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		printMessage("Error al abandonar el grupo.")
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		printMessage("Abandonó el grupo correctamente.")
		authenticatedOptions()
	} else {
		printMessage("Error al abandonar el grupo.")
	}
}

// Función para la administración del grupo
func manageGroup() {
	println()
	for {
		fmt.Println("Seleccione una opción de administración de grupo:")
		fmt.Println("1. Ver miembros del grupo")
		fmt.Println("2. Eliminar miembro del grupo")
		fmt.Println("3. Editar descripción del grupo")
		fmt.Println("4. Eliminar grupo")
		fmt.Println("5. Regresar al menú anterior")

		var option int
		fmt.Scanln(&option)

		switch option {
		case 1:
			viewGroupMembers()
		case 2:
			removeGroupMember()
		case 3:
			editGroupDescription()
		case 4:
			deleteGroup()
		case 5:
			println()
			return
		default:
			printMessage("Opción no válida, intente nuevamente.")
		}
	}
}

// Función para ver los miembros del grupo
func viewGroupMembers() ([]struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}, error) {
	if authToken == "" {
		printMessage("Por favor, inicie sesión primero.")
		return nil, fmt.Errorf("no authentication token provided")
	}

	req, err := http.NewRequest("GET", serviceURL+"/groupMembers", nil)
	if err != nil {
		printMessage("Error al crear la solicitud.")
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		printMessage("Error al obtener los miembros del grupo.")
		return nil, err
	}

	defer resp.Body.Close()

	var members []struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&members); err != nil {
		printMessage("Error al decodificar la respuesta.")
		return nil, err
	}

	if len(members) == 0 {
		printMessage("No hay miembros en el grupo.")
		return nil, nil
	}

	println()
	fmt.Println("Miembros del grupo:")
	for _, member := range members {
		fmt.Printf("ID: %d - Username: %s\n", member.ID, member.Username)
	}

	println()

	return members, nil
}

// Función para eliminar un miembro del grupo
func removeGroupMember() {
	members, err := viewGroupMembers()
	if err != nil || members == nil {
		return
	}

	fmt.Println("Ingrese el ID del miembro a eliminar (o 'q' para salir):")
	var input string
	fmt.Scanln(&input)
	if input == "q" {
		return
	}

	memberID, err := strconv.Atoi(input)
	if err != nil || memberID < 1 {
		printMessage("ID de miembro inválido.")
		return
	}

	reqBody := struct {
		MemberID int `json:"member_id"`
	}{
		MemberID: memberID,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		printMessage("Error al codificar la solicitud.")
		return
	}

	req, err := http.NewRequest("POST", serviceURL+"/removeGroupMember", bytes.NewBuffer(reqBytes))
	if err != nil {
		fmt.Println("Error al crear la solicitud.")
		return
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		printMessage("Error al eliminar el miembro del grupo.")
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		printMessage("Miembro del grupo eliminado correctamente.")
	} else {
		printMessage("Error al eliminar el miembro del grupo.")
	}
}

// Función para editar la descripción del grupo
func editGroupDescription() {
	if authToken == "" {
		printMessage("Por favor, inicie sesión primero.")
		return
	}

	fmt.Println("Ingrese la nueva descripción del grupo:")
	scanner := bufio.NewScanner(os.Stdin)
	var newDescription string
	if scanner.Scan() {
		newDescription = scanner.Text()
	}

	reqBody := struct {
		Descripcion string `json:"descripcion"`
	}{
		Descripcion: newDescription,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		printMessage("Error al codificar la solicitud.")
		return
	}

	req, err := http.NewRequest("PUT", serviceURL+"/groupDescription", bytes.NewBuffer(reqBytes))
	if err != nil {
		fmt.Println("Error al crear la solicitud.")
		return
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		printMessage("Error al editar la descripción del grupo.")
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		printMessage("Descripción del grupo actualizada correctamente.")
	} else {
		printMessage("Error al actualizar la descripción del grupo.")
	}
}

// Función para eliminar un grupo
func deleteGroup() {
	if authToken == "" {
		printMessage("Por favor, inicie sesión primero.")
		return
	}

	req, err := http.NewRequest("DELETE", serviceURL+"/deleteGroup", nil)
	if err != nil {
		printMessage("Error al crear la solicitud.")
		return
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		printMessage("Error al eliminar el grupo.")
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		printMessage("Grupo eliminado correctamente.")
		userGroupInfo()
		authenticatedOptions()
	} else {
		printMessage("Error al eliminar el grupo.")
	}
}
