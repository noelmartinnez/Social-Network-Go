# Social Network with Go

Esta es la aplicación perteneciente a la asignatura de "Seguridad en el Diseño del Software" del 4º año del Grado de Ingeniería Informática en la Universidad de Alicante.  

## Características

La definición y las características de la aplicación son las siguientes:

* Arquitectura cliente/servidor, realizándose ambos en Go.
* Mecanismo de autentificación seguro (gestión de contraseñas, identidades y sesión).
* Transporte de red seguro entre cliente y servidor (se puede emplear HTTPS o TLS).
* Almacenamiento seguro (cifrado en descanso).
* Sistema de gestión de contenido general (público).
* Sistema de comunicación privado (cifrado) entre usuarios.
* Gestión de categorías de contenido o grupos de usuarios (puede incluir seguridad adicional).
* Gestión de diferentes roles de usuarios (administradores, moderadores, etc.).
* Sistema de registro de eventos (logging), para mejorar la trazabilidad (remoto).
* Sistema de copia de seguridad (backup), para recuperarse ante incidentes (remoto).

## Ejecución

Para poder poner en marcha la aplicación hay que ejecutar los siguientes comandos:
- go run .\main.go server
- go run .\main.go client

Antes hay que haber creado la base de datos en el entorno que se prefiera y debe estar activa.

---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

# Social Network with Go

This is the application for the subject "Security in Software Design" in the 4th year of the Computer Engineering degree at the University of Alicante.

## Features

The definition and features of the application are as follows:

* Client/server architecture, both implemented in Go.
* Secure authentication mechanism (management of passwords, identities, and sessions).
* Secure network transport between client and server (HTTPS or TLS can be used).
* Secure storage (encryption at rest).
* General content management system (public).
* Private communication system (encrypted) between users.
* Management of content categories or user groups (may include additional security).
* Management of different user roles (administrators, moderators, etc.).
* Event logging system, to improve traceability (remote).
* Backup system, to recover from incidents (remote).

## Execution

To start the application, execute the following commands:
- `go run .\main.go server`
- `go run .\main.go client`

Before that, you need to create the database in the preferred environment and ensure it is active.

