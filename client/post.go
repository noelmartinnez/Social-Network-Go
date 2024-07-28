package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// Función para publicar un post
func postTexto() {
	if authToken == "" {
		printMessage("Por favor, inicie sesión primero.")
		return
	}

	var post struct {
		Post
		EsPublico bool `json:"es_publico"`
	}

	scanner := bufio.NewScanner(os.Stdin)

	println()
	fmt.Println("Ingrese el título del post:")
	if scanner.Scan() {
		post.Titulo = scanner.Text()
	}

	println()
	fmt.Println("Ingrese el texto del post:")
	if scanner.Scan() {
		post.Texto = scanner.Text()
	}

	if currentGroupID > 0 {
		println()
		fmt.Println("¿El post es público? (sí/no):")
		if scanner.Scan() {
			respuesta := scanner.Text()
			post.EsPublico = strings.ToLower(respuesta) == "sí"
		}
	} else {
		post.EsPublico = true
	}

	postBytes, err := json.Marshal(post)
	if err != nil {
		printMessage("Error al codificar el post.")
		return
	}

	req, err := http.NewRequest("POST", serviceURL+"/createPost", bytes.NewBuffer(postBytes))
	if err != nil {
		fmt.Println("Error al crear la solicitud.")
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		printMessage("Error al publicar el post.")
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		printMessage("Post publicado correctamente.")
	} else {
		printMessage("Error al publicar el post.")
	}
}

// Función para ver los posts
func viewPosts() {
	if authToken == "" {
		printMessage("Por favor, inicie sesión primero.")
		return
	}

	req, err := http.NewRequest("GET", serviceURL+"/posts", nil)
	if err != nil {
		printMessage("Error al crear la solicitud.")
		return
	}

	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		printMessage("Error al obtener los posts.")
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
		printMessage("No hay posts disponibles.")
		return
	}

	println()
	fmt.Println("Posts disponibles:")
	for _, post := range posts {
		fmt.Printf("ID: %d\nTítulo: %s\nTexto: %s\n\n", post.ID, post.Titulo, post.Texto)
	}
}

// Función para eliminar un post
func deletePost() {
	if authToken == "" {
		printMessage("Por favor, inicie sesión primero.")
		return
	}

	validPostIDs := getValidPostIDs()

	fmt.Println("Seleccione el ID del post a eliminar:")
	viewPosts()

	var postID int
	fmt.Scanln(&postID)

	if !contains(validPostIDs, postID) {
		printMessage("ID inválido o fuera de rango. Operación cancelada.")
		return
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/post/%d", serviceURL, postID), nil)
	if err != nil {
		printMessage("Error al crear la solicitud.")
		return
	}

	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)

	if err != nil {
		printMessage("Error al eliminar el post.")
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		printMessage("Post eliminado correctamente.")
	} else {
		printMessage("Error al eliminar el post.")
	}
}
