// Comunicación entre C/S asegurada mediante HTTPS/TLS.
// Utilizando certificado.crt y llave.key para el servidor y el cliente.
// Donde el certificado.crt es el certificado del servidor y la llave.key es la llave privada del servidor.

// Se usa JWT para manejar sesiones de usuario y gestionar la autenticación.
// Argon2 + Salting aleatorio para almacenar contraseñas de usuario de forma segura.

// Para crear el certificado y la clave privada se ha usado:
// openssl req -new -newkey rsa:4096 -x509 -sha256 -days 365 -nodes -out certificado.crt -keyout llave.key -subj "/CN=localhost" -addext "subjectAltName=DNS:localhost"

// Se ha implementado el sistema de roles de los usuarios, donde se distingue entre usuario normal, moderador y administrador.
// El administrador puede eliminar usuarios y posts, y el moderador puede eliminar posts. También pudiendo hacer las misma acciones que un usuario normal.
// El primer usuario que se registra en la aplicación se convierte en administrador y los demás en usuarios normales.
// Entonces el administrador puede promocionar a un usuario normal a moderador siempre que no lo sea ya.

// Hay un sistema de mensajería en el cual los usuarios pueden mandarse mensajes entre ellos y además tienen cifrado simétrico en cada extremo, es decir, los
// usuarios cifran y descifran los mensajes mientras que el servidor no sabe lo que hay en estos para tener más seguridad los usuarios

package main

import (
	"fmt"
	"os"
	"sds/client"
	"sds/server"
)

func main() {
	fmt.Println("Práctica de desarrollo de software - Red social")
	s := "Introduce server para funcionalidad de servidor y client para funcionalidad de cliente"

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "server":
			fmt.Println("Entrando en modo servidor...")
			server.Run()
		case "client":
			fmt.Println("Entrando en modo cliente...")
			client.Run()
		default:
			fmt.Println("Parámetro '", os.Args[1], "' desconocido. ", s)
		}
	} else {
		fmt.Println(s)
	}
}
