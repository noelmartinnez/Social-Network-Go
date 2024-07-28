package client

import (
	"net/http"
)

// Función para restaurar la base de datos
func restoreDatabase() {
	if authToken == "" {
		printMessage("Por favor, inicie sesión primero.")
		return
	}

	req, err := http.NewRequest("GET", serviceURL+"/restore", nil)
	if err != nil {
		printMessage("Error al crear la solicitud.")
		return
	}

	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		printMessage("Error al restaurar la base de datos.")
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		printMessage("Base de datos restaurada correctamente.")
	} else {
		printMessage("Error al restaurar la base de datos.")
	}
}

// Función para crear un backup
func createBackup() {
	if authToken == "" {
		printMessage("Por favor, inicie sesión primero.")
		return
	}

	req, err := http.NewRequest("GET", serviceURL+"/backup", nil)
	if err != nil {
		printMessage("Error al crear la solicitud.")
		return
	}

	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		printMessage("Error al crear el backup.")
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		printMessage("Backup creado correctamente.")
	} else {
		printMessage("Error al crear el backup.")
	}
}
