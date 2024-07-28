package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

// Realizar una copia de seguridad de la base de datos
func backupDatabase(w http.ResponseWriter, r *http.Request) {
	err := backupAndEncrypt()
	if err != nil {
		http.Error(w, "Backup failed: "+err.Error(), http.StatusInternalServerError)
		logEvent(0, "backup_failed", "Backup failed", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	w.Write([]byte("Backup successful"))
}

// Restaurar la base de datos
func restoreDatabase(w http.ResponseWriter, r *http.Request) {
	err := decryptAndRestore()
	if err != nil {
		http.Error(w, "Restore failed: "+err.Error(), http.StatusInternalServerError)
		logEvent(0, "restore_failed", "Restore failed", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
		return
	}

	w.Write([]byte("Restore successful"))
}

// Encriptar los datos usando AES y una clave de 32 bytes
func backupAndEncrypt() error {
	const (
		databaseUser      = "root"
		databasePassword  = "root"
		databaseName      = "sds"
		databaseHost      = "localhost" // Dirección del servidor remoto
		databasePort      = "3306"      // Puerto del servidor MySQL
		backupFileName    = "backup.sql"
		encryptedFileName = "backup.enc"
		decryptionKey     = "thisis32bitlongpassphraseimusing"
	)

	cmd := exec.Command("mysqldump", "-h"+databaseHost, "-P"+databasePort, "-u"+databaseUser, "-p"+databasePassword, databaseName)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		logEvent(0, "backup_failed", "mysqldump failed", "", map[string]interface{}{"error": err.Error()})
		return fmt.Errorf("mysqldump failed: %v", err)

	}

	encryptedData, err := encrypt(string(out.Bytes()), []byte(decryptionKey))
	if err != nil {
		logEvent(0, "backup_failed", "encryption failed", "", map[string]interface{}{"error": err.Error()})
		return fmt.Errorf("encryption failed: %v", err)

	}

	encryptedBytes := []byte(encryptedData)

	err = ioutil.WriteFile(encryptedFileName, encryptedBytes, 0644)
	if err != nil {
		logEvent(0, "backup_failed", "failed to write encrypted file", "", map[string]interface{}{"error": err.Error()})
		return fmt.Errorf("failed to write encrypted file: %v", err)
	}

	fmt.Println("Database backup and encryption successful")
	return nil
}

// Desencriptar los datos y restaurar la base de datos
func decryptAndRestore() error {
	const (
		databaseUser      = "root"
		databasePassword  = "root"
		databaseName      = "sds2"
		databaseHost      = "localhost"
		databasePort      = "3306"
		encryptedFileName = "backup.enc"
		decryptionKey     = "thisis32bitlongpassphraseimusing"
	)

	encryptedData, err := ioutil.ReadFile(encryptedFileName)
	if err != nil {
		logEvent(0, "restore_failed", "failed to read encrypted file", "", map[string]interface{}{"error": err.Error()})
		return fmt.Errorf("failed to read encrypted file: %v", err)
	}

	encryptedString := string(encryptedData)

	decryptedData, err := decrypt((encryptedString), []byte(decryptionKey))
	if err != nil {
		logEvent(0, "restore_failed", "decryption failed", "", map[string]interface{}{"error": err.Error()})
		return fmt.Errorf("decryption failed: %v", err)
	}

	tmpFileName := "decrypted.sql"
	err = ioutil.WriteFile(tmpFileName, decryptedData, 0644)
	if err != nil {
		logEvent(0, "restore_failed", "failed to write decrypted file", "", map[string]interface{}{"error": err.Error()})
		return fmt.Errorf("failed to write decrypted file: %v", err)
	}
	defer os.Remove(tmpFileName)

	cmd := exec.Command("mysql", "-h"+databaseHost, "-P"+databasePort, "-u"+databaseUser, "-p"+databasePassword, databaseName)
	cmd.Stdin = bytes.NewReader(decryptedData)
	err = cmd.Run()
	if err != nil {
		logEvent(0, "restore_failed", "mysql restore failed", "", map[string]interface{}{"error": err.Error()})
		return fmt.Errorf("mysql restore failed: %v", err)
	}

	fmt.Println("Database decryption and restoration successful")
	return nil
}
