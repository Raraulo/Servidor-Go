package seguridad

import (
	"encoding/json"
	"net/http"
	"srvr/db"
)

type TopProducto struct {
	ID           int64   `json:"id"`
	Nombre       string  `json:"nombre"`
	Descripcion  string  `json:"descripcion"`
	Precio       float64 `json:"precio"`
	Stock        int     `json:"stock"`
	ImagenURL    string  `json:"imagen_url"`
	AnimalID     int64   `json:"animal_id"`
	Disponible   bool    `json:"disponible"`
	FechaIngreso string  `json:"fecha_ingreso"`
	TotalVendido int     `json:"total_vendido"`
}

type VentasPorDia struct {
	Dia         string  `json:"dia"`
	TotalVentas float64 `json:"total_ventas"`
}

func TopProductos(rw http.ResponseWriter, r *http.Request) {
	rows, err := db.BaseDeDatos.Raw(`
  SELECT p.id, p.nombre, p.descripcion, p.precio, p.stock, p.imagen_url, p.animal_id, p.disponible, p.fecha_ingreso, SUM(fd.cantidad) AS total_vendido
  FROM factura_detalle fd
  JOIN producto p ON fd.producto_id = p.id
  GROUP BY p.id, p.nombre, p.descripcion, p.precio, p.stock, p.imagen_url, p.animal_id, p.disponible, p.fecha_ingreso
  ORDER BY total_vendido DESC
  LIMIT 6;
`).Rows()

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var result []TopProducto
	for rows.Next() {
		var p TopProducto
		// OJO: El orden aquí debe ser IGUAL al SELECT
		if err := rows.Scan(
			&p.ID,
			&p.Nombre,
			&p.Descripcion,
			&p.Precio,
			&p.Stock,
			&p.ImagenURL,
			&p.AnimalID,
			&p.Disponible,
			&p.FechaIngreso, // Si da error por formato de fecha, usa time.Time y adapta
			&p.TotalVendido,
		); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		result = append(result, p)
	}
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(result)
}

func VentasPorDia2026(rw http.ResponseWriter, r *http.Request) {
	rows, err := db.BaseDeDatos.Raw(`
        SELECT 
            fecha::date AS dia,
            SUM(total) AS total_ventas
        FROM 
            factura
        WHERE 
            EXTRACT(YEAR FROM fecha) = 2026
        GROUP BY 
            dia
        ORDER BY 
            dia;
    `).Rows()

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var result []VentasPorDia
	for rows.Next() {
		var v VentasPorDia
		if err := rows.Scan(&v.Dia, &v.TotalVentas); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		result = append(result, v)
	}

	// Si no hay ventas, enviamos un array vacío en lugar de null
	if result == nil {
		result = []VentasPorDia{}
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(result)
}

// Estructura para el reporte de mejores clientes
type TopCliente struct {
	Nombre        string  `json:"prsnnmbr"`
	TotalComprado float64 `json:"total_comprado"`
}

// Handler para el reporte de los 5 mejores clientes por compras
func TopClientes(rw http.ResponseWriter, r *http.Request) {
	rows, err := db.BaseDeDatos.Raw(`
		SELECT p.prsnnmbr, SUM(f.total) AS total_comprado
		FROM factura f
		JOIN prsn p ON f.persona_id = p.prsn__id
		GROUP BY p.prsnnmbr
		ORDER BY total_comprado DESC
		LIMIT 5;
	`).Rows()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var result []TopCliente
	for rows.Next() {
		var c TopCliente
		if err := rows.Scan(&c.Nombre, &c.TotalComprado); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		result = append(result, c)
	}
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(result)
}

// Estructura para reporte de stock disponible de todos los productos
type ProductoStockDisponible struct {
	Nombre string `json:"nombre"`
	Stock  int    `json:"stock"`
}

// Handler para stock disponible de todos los productos
func StockDisponible(rw http.ResponseWriter, r *http.Request) {
	rows, err := db.BaseDeDatos.Raw(`SELECT nombre, stock 
FROM producto 
ORDER BY stock DESC 
LIMIT 10;
`).Rows()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var result []ProductoStockDisponible
	for rows.Next() {
		var p ProductoStockDisponible
		if err := rows.Scan(&p.Nombre, &p.Stock); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		result = append(result, p)
	}
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(result)
}

// Estructura para el valor total del stock
type ValorTotalStock struct {
	Valor float64 `json:"valor_total_stock"`
}

// Handler para consultar el valor total del stock disponible
func ValorStock(rw http.ResponseWriter, r *http.Request) {
	rows, err := db.BaseDeDatos.Raw(`
		SELECT COALESCE(SUM(precio * stock), 0) AS valor_total_stock
		FROM producto
		WHERE stock > 0;
	`).Rows()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var v ValorTotalStock
	if rows.Next() {
		if err := rows.Scan(&v.Valor); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(v)
}

// Estructura para el total recaudado
type TotalRecaudado struct {
	Total float64 `json:"total_recaudado"`
}

// Handler para consultar el total recaudado de todas las facturas
func TotalRecaudadoHandler(rw http.ResponseWriter, r *http.Request) {
	rows, err := db.BaseDeDatos.Raw(`
		SELECT COALESCE(SUM(total), 0) AS total_recaudado
		FROM factura;
	`).Rows()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var t TotalRecaudado
	if rows.Next() {
		if err := rows.Scan(&t.Total); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(t)
}

// Estructura para el total de productos vendidos
type TotalProductosVendidos struct {
	Total int `json:"total_productos_vendidos"`
}

// Handler para consultar el total de productos vendidos (suma de cantidad)
func TotalProductosVendidosHandler(rw http.ResponseWriter, r *http.Request) {
	rows, err := db.BaseDeDatos.Raw(`
		SELECT COALESCE(SUM(cantidad), 0) AS total_productos_vendidos
		FROM factura_detalle;
	`).Rows()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var t TotalProductosVendidos
	if rows.Next() {
		if err := rows.Scan(&t.Total); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(t)
}

// Handler para productos con stock menor a 10
func ProductosBajoStock(rw http.ResponseWriter, r *http.Request) {
	rows, err := db.BaseDeDatos.Raw(`
		SELECT nombre, stock 
		FROM producto 
		WHERE stock <= 10
		ORDER BY stock ASC;
	`).Rows()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var result []ProductoStockDisponible
	for rows.Next() {
		var p ProductoStockDisponible
		if err := rows.Scan(&p.Nombre, &p.Stock); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		result = append(result, p)
	}
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(result)
}
