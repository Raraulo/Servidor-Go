package rutas

import (
	"srvr/handlers/seguridad"
	"srvr/token"

	"github.com/gorilla/mux"
)

func Rutas_producto(r *mux.Router) {
	// Rutas públicas
	r.HandleFunc("/api/producto", seguridad.GetProductos).Methods("GET")
	r.HandleFunc("/api/producto/{id:[0-9]+}", seguridad.GetProducto).Methods("GET")

	// Solo admin puede crear, actualizar y eliminar
	r.HandleFunc("/api/producto", token.VerificaAdmin(seguridad.CreateProducto)).Methods("POST")
	r.HandleFunc("/api/producto/{id:[0-9]+}", token.VerificaAdmin(seguridad.UpdateProducto)).Methods("PUT")
	r.HandleFunc("/api/producto/{id:[0-9]+}", token.VerificaAdmin(seguridad.DeleteProducto)).Methods("DELETE")
	r.HandleFunc("/api/categoria", seguridad.GetCategorias).Methods("GET")

}
