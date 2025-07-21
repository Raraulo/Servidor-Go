package token

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// SecretKey debe estar definida en tu configuración principal
// Ejemplo: var SecretKey = []byte(os.Getenv("clave-secreta"))

// VerificaAdmin permite acceso solo al usuario admin1
func VerificaAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, `{"error":"Falta el token", "ok":false}`, http.StatusForbidden)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
			}
			return SecretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"error":"Token no válido", "ok":false}`, http.StatusForbidden)
			return
		}

		username := fmt.Sprintf("%v", claims["Username"])
		fmt.Println("🔐 VerificaAdmin - Usuario autenticado:", username)

		if username != "admin1" {
			http.Error(w, `{"error":"No autorizado (solo admin)", "ok":false}`, http.StatusForbidden)
			return
		}

		// Opcional: pasar información útil
		r.Header.Set("Usuario", username)
		r.Header.Set("PersonaID", fmt.Sprintf("%v", claims["Id"]))

		next.ServeHTTP(w, r)
	}
}

// TokenMiddleware permite acceso a cualquier usuario autenticado (no solo admin)
func TokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, `{"error":"Falta el token", "ok":false}`, http.StatusForbidden)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
			}
			return SecretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"error":"Token inválido", "ok":false}`, http.StatusUnauthorized)
			return
		}

		// Guardamos en el header temporalmente el usuario y persona_id
		r.Header.Set("Usuario", fmt.Sprintf("%v", claims["Username"]))
		r.Header.Set("PersonaID", fmt.Sprintf("%v", claims["Id"]))

		next.ServeHTTP(w, r)
	}
}
