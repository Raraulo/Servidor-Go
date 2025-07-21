package rutas

import (
	"srvr/handlers/seguridad"

	"github.com/gorilla/mux"
)

func Rutas_usro(mx *mux.Router) {
	// EndPoints de usuario
	mx.HandleFunc("/api/user", seguridad.GetPersonas).Methods("GET")
	mx.HandleFunc("/api/user/{id:[0-9]+}", seguridad.GetPersona).Methods("GET")
	mx.HandleFunc("/api/user", seguridad.CreaPersona).Methods("POST")
	mx.HandleFunc("/api/user/{id:[0-9]+}", seguridad.UpdatePersona).Methods("PUT")
	mx.HandleFunc("/api/user/{id:[0-9]+}", seguridad.DeletePersona).Methods("DELETE")

	// ✅ Endpoint de login
	mx.HandleFunc("/api/login", seguridad.Login).Methods("POST")
}
