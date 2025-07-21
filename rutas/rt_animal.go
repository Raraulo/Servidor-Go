package rutas

import (
	"net/http"
	"srvr/handlers/seguridad"

	"github.com/gorilla/mux"
)

// Rutas_animal registra los endpoints para animal (tipo de mascota)
func Rutas_animal(r *mux.Router) {
	r.HandleFunc("/api/animal", seguridad.ObtenerAnimales).Methods(http.MethodGet)
}
