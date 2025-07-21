package seguridad

import (
	"encoding/json"
	"net/http"
	"srvr/db"
	"srvr/dominios"
)

func ObtenerAnimales(w http.ResponseWriter, r *http.Request) {
	var animales dominios.Animales
	if err := db.BaseDeDatos.Find(&animales).Error; err != nil {
		http.Error(w, "No se pudieron obtener los animales", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animales)
}
