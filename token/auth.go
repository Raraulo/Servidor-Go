package token

import (
	"encoding/json"
	"net/http"
	"time"

	"srvr/db"
	"srvr/dominios"

	"github.com/golang-jwt/jwt/v5"
)

// ====================== LOGIN =========================

func ChequeaUsuario(w http.ResponseWriter, r *http.Request) {
	type Usuario struct {
		Login string
		Pass  string
	}
	pr := Usuario{}

	if err := json.NewDecoder(r.Body).Decode(&pr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	persona := dominios.Persona{}
	db.BaseDeDatos.Where("prsnlogn = ? AND prsnpass = ?", pr.Login, pr.Pass).First(&persona)

	if persona.Id == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ok":  false,
			"msg": "Login o password incorrectos",
		})
		return
	}

	GeneraToken(w, r, persona.Id, persona.Login)
}

// ====================== GENERADOR DE TOKEN =========================

func GeneraToken(w http.ResponseWriter, r *http.Request, id int64, login string) {
	now := time.Now()
	claims := jwt.MapClaims{
		"Id":        id,
		"Username":  login,
		"ExpiresAt": now.Add(30 * 24 * time.Hour).Unix(),
		"IssuedAt":  now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	var persona dominios.Persona
	if err := db.BaseDeDatos.First(&persona, id).Error; err != nil {
		http.Error(w, `{"ok":false, "msg":"No se pudo obtener los datos del usuario"}`, http.StatusInternalServerError)
		return
	}

	// respuesta final esperada por el frontend
	response := map[string]interface{}{
		"ok":    true,
		"token": tokenString,
		"Registro": map[string]interface{}{
			"id":     persona.Id,
			"nombre": persona.Nombre,
			"login":  persona.Login,
			"correo": persona.Mail,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ====================== GET NUEVO TOKEN (REFRESH) =========================

func RevisaTokenHandler(w http.ResponseWriter, r *http.Request) {
	token, err := TokenOk(w, r)
	if err != nil {
		http.Error(w, `{"ok":false, "msg":"Token inválido"}`, http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		http.Error(w, `{"ok":false}`, http.StatusUnauthorized)
		return
	}

	idRaw := claims["Id"]
	idFloat, ok := idRaw.(float64)
	if !ok {
		http.Error(w, `{"ok":false, "msg":"Id inválido"}`, http.StatusUnauthorized)
		return
	}

	var persona dominios.Persona
	db.BaseDeDatos.First(&persona, int64(idFloat))

	if persona.Id == 0 {
		http.Error(w, `{"ok":false, "msg":"No existe usuario"}`, http.StatusUnauthorized)
		return
	}

	GeneraToken(w, r, persona.Id, persona.Login)
}
