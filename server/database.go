package server

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Inicializar la conexión a la base de datos
func init() {
	var err error
	db, err = sql.Open("mysql", "root:noel@tcp(localhost:3306)/sds")
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to database successfully.")
}
