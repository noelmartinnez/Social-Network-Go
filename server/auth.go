package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// Middleware para autenticar las solicitudes
func authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tknStr := r.Header.Get("Authorization")
		if tknStr == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			fmt.Println("Falta el header de Authorization")
			logEvent(0, "authentication_failed", "Missing Authorization header", r.RemoteAddr, nil)
			return
		}

		tknStrParts := strings.Fields(tknStr)
		if len(tknStrParts) != 2 || tknStrParts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			fmt.Println("Formato del header de Authorization inválido")
			logEvent(0, "authentication_failed", "Invalid Authorization header format", r.RemoteAddr, map[string]interface{}{"header": tknStr})
			return
		}

		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tknStrParts[1], claims, func(token *jwt.Token) (interface{}, error) {
			// Check the signing method
			fmt.Println("entro a parse with claims")
			return jwtKey, nil
		})
		if err != nil {
			fmt.Println("Error parsing token:", err)
			http.Error(w, "Error parsing token", http.StatusUnauthorized)
			logEvent(0, "authentication_failed", "Error parsing token", r.RemoteAddr, map[string]interface{}{"error": err.Error()})
			return
		}
		if !tkn.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			fmt.Println("Token inválido")
			logEvent(0, "authentication_failed", "Invalid token", r.RemoteAddr, nil)
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		logEvent(0, "authentication_success", "User authenticated successfully", r.RemoteAddr, map[string]interface{}{"username": claims.Username})
		next(w, r.WithContext(ctx))
	}
}
