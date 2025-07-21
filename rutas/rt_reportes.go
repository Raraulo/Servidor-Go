package rutas

import (
	"srvr/handlers/seguridad"
	"srvr/token"

	"github.com/gorilla/mux"
)

// Endpoints para reportes (solo admin, protegidos)
func RegistrarRutasReportes(r *mux.Router) {
	r.HandleFunc("/api/reportes/top-productos", token.VerificaAdmin(seguridad.TopProductos)).Methods("GET")

	r.HandleFunc("/api/reportes/top-clientes", token.VerificaAdmin(seguridad.TopClientes)).Methods("GET")
	r.HandleFunc("/api/reportes/stock-disponible", token.VerificaAdmin(seguridad.StockDisponible)).Methods("GET")
	r.HandleFunc("/api/reportes/bajo-stock", seguridad.ProductosBajoStock).Methods("GET")

	// Nuevos endpoints para reportes generales:
	r.HandleFunc("/api/reportes/valor-stock", token.VerificaAdmin(seguridad.ValorStock)).Methods("GET")
	r.HandleFunc("/api/reportes/total-recaudado", token.VerificaAdmin(seguridad.TotalRecaudadoHandler)).Methods("GET")
	r.HandleFunc("/api/reportes/total-vendido", token.VerificaAdmin(seguridad.TotalProductosVendidosHandler)).Methods("GET")
	//  ventas por día de 2025
	r.HandleFunc("/api/reportes/ventas-por-dia-2025", token.VerificaAdmin(seguridad.VentasPorDia2025)).Methods("GET")
}

// Endpoints públicos (acceso libre, para catálogo y web de usuario)
func RegistrarRutasPublicas(r *mux.Router) {
	r.HandleFunc("/api/catalogo/mas-vendidos", seguridad.TopProductos).Methods("GET")
}
