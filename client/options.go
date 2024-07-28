package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Función para obtener información del grupo del usuario
func userGroupInfo() {
	if authToken == "" {
		return
	}

	req, err := http.NewRequest("GET", serviceURL+"/userGroup", nil)
	if err != nil {
		return
	}

	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	var groupInfo struct {
		GroupID int  `json:"group_id"`
		IsAdmin bool `json:"is_admin"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&groupInfo); err != nil {
		return
	}

	currentGroupID = groupInfo.GroupID
	isAdmin = groupInfo.IsAdmin
}

// Función para mostrar las opciones no autenticadas
func unauthenticatedOptions() {
	for {
		fmt.Println("Seleccione una opción:")
		fmt.Println("1. Registro")
		fmt.Println("2. Iniciar sesión")
		fmt.Println("3. Salir")

		var option int
		fmt.Scanln(&option)

		switch option {
		case 1:
			registerUser()
		case 2:
			loginUser()
			if authToken != "" {
				authenticatedOptions()
			}
		case 3:
			fmt.Println("Saliendo...")
			os.Exit(0)
		default:
			printMessage("Opción no válida, intente nuevamente.")
		}
	}
}

// Función para mostrar opciones basadas en el rol del usuario
func authenticatedOptions() {
	userGroupInfo()
	switch rol {
	case "administrador":
		adminOptions()
	case "moderador":
		moderatorOptions()
	case "usuario":
		if isAdmin {
			adminGroupOptions()
		} else if currentGroupID > 0 {
			memberGroupOptions()
		} else {
			userOptions()
		}
	default:
		fmt.Println("Rol no reconocido.")
		return
	}
}

// Función para las opciones de administrador
func adminOptions() {
	for {
		fmt.Println("Seleccione una opción de administrador:")
		fmt.Println("1. Publicar contenido")
		fmt.Println("2. Ver publicaciones")
		fmt.Println("3. Ver chats ")
		fmt.Println("4. Eliminar post")
		fmt.Println("5. Eliminar usuario")
		fmt.Println("6. Promover a moderador")
		fmt.Println("7. Crear backup")
		fmt.Println("8. Restaurar base de datos")
		fmt.Println("9. Cerrar sesión")

		var option int
		fmt.Scanln(&option)

		switch option {
		case 1:
			postTexto()
		case 2:
			viewPosts()
		case 3:
			viewUsers()
		case 4:
			deletePost()
		case 5:
			deleteUser()
		case 6:
			promoteUserToModerator()
		case 7:
			createBackup()
		case 8:
			restoreDatabase()
		case 9:
			logout()
			return
		default:
			printMessage("Opción no válida, intente nuevamente.")
		}
	}
}

// Función para las opciones de moderador
func moderatorOptions() {
	for {
		fmt.Println("Seleccione una opción de moderador:")
		fmt.Println("1. Publicar contenido")
		fmt.Println("2. Ver publicaciones")
		fmt.Println("3. Ver chats")
		fmt.Println("4. Eliminar post")
		fmt.Println("5. Cerrar sesión")

		var option int
		fmt.Scanln(&option)

		switch option {
		case 1:
			postTexto()
		case 2:
			viewPosts()
		case 3:
			viewUsers()
		case 4:
			deletePost()
		case 5:
			logout()
			return
		default:
			printMessage("Opción no válida, intente nuevamente.")
		}
	}
}

// Función para las opciones de usuario sin grupo
func userOptions() {
	for {
		fmt.Println("Seleccione una opción:")
		fmt.Println("1. Publicar contenido")
		fmt.Println("2. Ver publicaciones")
		fmt.Println("3. Ver chats")
		fmt.Println("4. Crear grupo")
		fmt.Println("5. Ver grupos")
		fmt.Println("6. Cerrar sesión")

		var option int
		fmt.Scanln(&option)

		switch option {
		case 1:
			postTexto()
		case 2:
			viewPosts()
		case 3:
			viewUsers()
		case 4:
			createGroup()
		case 5:
			viewGroups(true)
		case 6:
			logout()
			unauthenticatedOptions()
		default:
			printMessage("Opción no válida, intente nuevamente.")
		}
	}
}

// Función para las opciones de administrador de grupo
func adminGroupOptions() {
	for {
		fmt.Println("Seleccione una opción:")
		fmt.Println("1. Publicar contenido")
		fmt.Println("2. Ver publicaciones")
		fmt.Println("3. Ver chats")
		fmt.Println("4. Administración de grupo")
		fmt.Println("5. Ver posts del grupo")
		fmt.Println("6. Cerrar sesión")

		var option int
		fmt.Scanln(&option)

		switch option {
		case 1:
			postTexto()
		case 2:
			viewPosts()
		case 3:
			viewUsers()
		case 4:
			manageGroup()
		case 5:
			viewGroupPosts()
		case 6:
			logout()
			unauthenticatedOptions()
		default:
			printMessage("Opción no válida, intente nuevamente.")
		}
	}
}

// Función para las opciones de usuario miembro de un grupo
func memberGroupOptions() {
	for {
		fmt.Println("Seleccione una opción:")
		fmt.Println("1. Publicar contenido")
		fmt.Println("2. Ver publicaciones")
		fmt.Println("3. Ver chats")
		fmt.Println("4. Ver posts del grupo")
		fmt.Println("5. Abandonar grupo")
		fmt.Println("6. Cerrar sesión")

		var option int
		fmt.Scanln(&option)

		switch option {
		case 1:
			postTexto()
		case 2:
			viewPosts()
		case 3:
			viewUsers()
		case 4:
			viewGroupPosts()
		case 5:
			leaveGroup()
		case 6:
			logout()
			unauthenticatedOptions()
		default:
			printMessage("Opción no válida, intente nuevamente.")
		}
	}
}
