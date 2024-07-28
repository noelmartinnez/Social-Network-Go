package server

import (
	"encoding/json"
	"log"
)

// Loguea un evento en la base de datos
func logEvent(userID int, eventType, eventDescription, ipAddress string, additionalData map[string]interface{}) {
	jsonData, err := json.Marshal(additionalData)
	if err != nil {
		log.Println("Error marshaling additional data:", err)
		return
	}

	_, err = db.Exec(
		"INSERT INTO event_logs (user_id, event_type, event_description, ip_address, additional_data) VALUES (?, ?, ?, ?, ?)",
		userID, eventType, eventDescription, ipAddress, jsonData,
	)
	if err != nil {
		log.Println("Error logging event:", err)
	}
}
