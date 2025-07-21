package dominios

import (
	"srvr/db"
	"time"
)

// Producto representa los artículos disponibles para venta
type Producto struct {
	ID           int64     `json:"id"`
	Nombre       string    `json:"nombre"`
	Descripcion  string    `json:"descripcion"`
	Precio       float64   `json:"precio"`
	Stock        int       `json:"stock"`
	ImagenURL    string    `json:"imagen_url"`
	AnimalID     int64     `json:"animal_id"` // ⚠️ IMPORTANTE: debe ser int
	Disponible   bool      `json:"disponible"`
	FechaIngreso time.Time `json:"fecha_ingreso"`
}

type Productos []Producto

// TableName especifica el nombre físico de la tabla
func (Producto) TableName() string {
	return "producto"
}

// MigrarProducto crea la tabla si no existe
func MigrarProducto() {
	db.BaseDeDatos.AutoMigrate(Producto{})
}
