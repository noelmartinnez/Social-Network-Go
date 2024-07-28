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
