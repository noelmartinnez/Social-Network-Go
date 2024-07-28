package client

import "time"

type User struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Nombre    string `json:"nombre"`
	Apellidos string `json:"apellidos"`
	Rol       string `json:"rol"`
}

type Post struct {
	Titulo string `json:"titulo"`
	Texto  string `json:"texto"`
}

type Rol struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type Message struct {
	ID          int       `json:"id"`
	SenderID    int       `json:"remitente_id"`
	RecipientID int       `json:"destinatario_id"`
	Content     string    `json:"contenido"`
	SentAt      time.Time `json:"fecha_envio"`
	SentByUser  bool      `json:"sent_by_user"`
}

type Grupo struct {
	ID              int    `json:"id"`
	Descripcion     string `json:"descripcion"`
	AdministradorID int    `json:"administrador_id"`
}
