package rutas

import (
	"srvr/handlers/seguridad"
	"srvr/token"

	"github.com/gorilla/mux"
)

// Rutas_factura registra todos los endpoints relacionados con facturas
func Rutas_factura(mx *mux.Router) {


	// Rutas para administrador
	mx.HandleFunc("/api/factura", token.VerificaAdmin(seguridad.ListarFacturas)).Methods("GET")
	mx.HandleFunc("/api/factura/{id:[0-9]+}", token.VerificaAdmin(seguridad.ObtenerFactura)).Methods("GET")
	mx.HandleFunc("/api/factura", token.VerificaAdmin(seguridad.CrearFactura)).Methods("POST")
	mx.HandleFunc("/api/factura/{id:[0-9]+}", token.VerificaAdmin(seguridad.ActualizarFactura)).Methods("PUT")
	mx.HandleFunc("/api/factura/{id:[0-9]+}", token.VerificaAdmin(seguridad.EliminarFactura)).Methods("DELETE")

	// Ver mis facturas con detalle (usuario autenticado)
	mx.HandleFunc("/api/factura/mias",
		token.TokenMiddleware(seguridad.VerFacturasConDetallePorUsuario)).
		Methods("GET")

	// Ver facturas de cualquier persona (opcional), basado en persona_id en URL
	mx.HandleFunc("/api/facturas/persona/{persona_id:[0-9]+}",
		token.TokenMiddleware(seguridad.VerFacturasPorUsuario)).
		Methods("GET")

	// Crear factura con detalle (compra)
	mx.HandleFunc("/api/factura/compra",
		token.TokenMiddleware(seguridad.CrearFacturaConDetalle)).
		Methods("POST")
}


