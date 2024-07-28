package client

import (
	_ "github.com/go-sql-driver/mysql"
)

var currentGroupID int
var isAdmin bool

var (
	client        = newTLSClient("certificado.crt")
	authToken     = ""
	rol           = ""
	currentUserID int
	serviceURL    = "https://localhost:443"
)

func Run() {
	unauthenticatedOptions()
}
