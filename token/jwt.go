package token

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go/request"
	"github.com/golang-jwt/jwt/v5"
)

// ✅ Validación general del token desde headers Authorization o token
func TokenOk(w http.ResponseWriter, r *http.Request) (*jwt.Token, error) {
	var tokenStr string

	// 1. Authorization: Bearer ...
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) > 7 && strings.HasPrefix(authHeader, "Bearer ") {
		tokenStr = authHeader[7:]
	} else {
		// 2. Header personalizado: token
		var err error
		tokenStr, err = request.HeaderExtractor{"token"}.ExtractToken(r)
		if err != nil || tokenStr == "undefined" {
			return nil, errors.New("sin token válido en headers")
		}
	}

	// 3. Parsear JWT
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// 4. Validar claims correctamente
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return token, nil
	}

	return nil, errors.New("token inválido o no autorizado")
}

// ✅ Middleware para cualquier usuario autenticado
func VerificaUsuario(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tok, err := TokenOk(w, r)
		if err != nil || !tok.Valid {
			http.Error(w, `{"error":"Usuario no autenticado", "ok":false}`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}
