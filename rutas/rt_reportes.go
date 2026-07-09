package rutas

import (
	"srvr/handlers/seguridad"
	"srvr/token"

	"github.com/gorilla/mux"
)

// RegistrarRutasReportes define los endpoints protegidos para administración
func RegistrarRutasReportes(r *mux.Router) {
	// Reportes de productos y clientes
	r.HandleFunc("/api/reportes/top-productos", token.VerificaAdmin(seguridad.TopProductos)).Methods("GET")
	r.HandleFunc("/api/reportes/top-clientes", token.VerificaAdmin(seguridad.TopClientes)).Methods("GET")
	r.HandleFunc("/api/reportes/stock-disponible", token.VerificaAdmin(seguridad.StockDisponible)).Methods("GET")
	r.HandleFunc("/api/reportes/bajo-stock", seguridad.ProductosBajoStock).Methods("GET")

	// Reportes financieros y métricas generales
	r.HandleFunc("/api/reportes/valor-stock", token.VerificaAdmin(seguridad.ValorStock)).Methods("GET")
	r.HandleFunc("/api/reportes/total-recaudado", token.VerificaAdmin(seguridad.TotalRecaudadoHandler)).Methods("GET")
	r.HandleFunc("/api/reportes/total-vendido", token.VerificaAdmin(seguridad.TotalProductosVendidosHandler)).Methods("GET")

	// Corregido: Asegúrate de que en seguridad.go exista la función VentasPorDia2026
	r.HandleFunc("/api/reportes/ventas-por-dia-2026", token.VerificaAdmin(seguridad.VentasPorDia2026)).Methods("GET")
}

// RegistrarRutasPublicas define endpoints de acceso libre
func RegistrarRutasPublicas(r *mux.Router) {
	r.HandleFunc("/api/catalogo/mas-vendidos", seguridad.TopProductos).Methods("GET")
}
