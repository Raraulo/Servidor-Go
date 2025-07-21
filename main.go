package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"srvr/dominios"
	"srvr/rutas"
	"srvr/token"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// middleware general para proteger rutas y extraer token
func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("→ %s %s\n", r.Method, r.URL.Path)

		// Rutas públicas que no requieren autenticación
		noAuthRequired := []struct {
			Method string
			Path   string
		}{
			{http.MethodPost, "/api/login"},
			{http.MethodPost, "/api/user"},
			{http.MethodGet, "/api/producto"},
		}

		isPublic := func() bool {
			for _, route := range noAuthRequired {
				if r.Method == route.Method && r.URL.Path == route.Path {
					return true
				}
			}
			// Permitir GET /api/producto/{id}
			if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/producto/") {
				return true
			}
			// Permitir endpoints públicos de catálogo
			if r.Method == http.MethodGet && (r.URL.Path == "/api/catalogo/mas-vendidos") {
				return true
			}
			return false
		}()

		if !isPublic {
			// Valida token y agrega cabeceras Usuario y PersonaID
			if _, err := token.TokenOk(w, r); err != nil {
				w.WriteHeader(http.StatusForbidden)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"ok":    strconv.FormatBool(false),
					"error": "Token no válido",
				})
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// init corre antes de main: ejecuta migraciones
func init() {
	dominios.MigrarPersona()
	dominios.MigrarAnimal()
	dominios.MigrarCategoria()
	dominios.MigrarProducto()
	dominios.MigrarFactura() // incluye factura_detalle
	dominios.MigrarVenta()
}

func main() {
	router := mux.NewRouter()

	// Endpoints de autenticación
	router.HandleFunc("/api/login", token.ChequeaUsuario).Methods(http.MethodPost)
	router.HandleFunc("/api/token", token.RevisaTokenHandler).Methods(http.MethodPost)

	// Endpoints del sistema
	rutas.Rutas_usro(router)
	rutas.Rutas_producto(router)
	rutas.Rutas_factura(router)
	rutas.Rutas_animal(router)

	// Endpoints de reportes (dashboard admin)
	rutas.RegistrarRutasReportes(router)

	// Endpoints públicos de catálogo (los ve cualquier usuario)
	rutas.RegistrarRutasPublicas(router)

	// Configuración CORS para Angular
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Token"},
		AllowCredentials: true,
		Debug:            true,
	})

	fmt.Println("🟢 Servidor iniciado en http://localhost:4000")
	log.Fatal(http.ListenAndServe(":4000", c.Handler(middleware(router))))
}
