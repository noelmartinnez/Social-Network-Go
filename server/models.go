package server

import "time"

type User struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Nombre    string `json:"nombre"`
	Apellidos string `json:"apellidos"`
	Rol       string `json:"rol"` // Nombre del rol del usuario
}

type Post struct {
	Titulo string `json:"titulo"`
	Texto  string `json:"texto"`
}
type Message struct {
	Id              int       `json:"id"`
	Remitente_id    int       `json:"remitente_id"`
	Destinatario_id int       `json:"destinatario_id"`
	Contenido       string    `json:"contenido"`
	Fecha_envio     time.Time `json:"fecha_envio"`
	SentByUser      bool      `json:"sent_by_user"`
}

type Rol struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type Grupo struct {
	ID              int    `json:"id"`
	Descripcion     string `json:"descripcion"`
	AdministradorID int    `json:"administrador_id"`
}
