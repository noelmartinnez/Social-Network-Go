package server

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

// Función que verifica si el usuario es administrador
func adminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(claimsKey).(*Claims)
		if !ok || claims.Rol != "administrador" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

// Función que verifica si el usuario está autenticado
func hashValue(value string) string {
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}
