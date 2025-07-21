package dominios

import (
	"srvr/db"
	"time"
)

// FacturaDetalle representa cada línea de la factura con el producto asociado
// Contiene información de producto y precio para cada elemento de la factura.
type FacturaDetalle struct {
	ID             int64   `json:"id" gorm:"primaryKey"`
	FacturaID      int64   `json:"-"`
	ProductoID     int64   `json:"producto_id"`
	Cantidad       int     `json:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario"`
	ProductoNombre string  `json:"producto_nombre"`
	ImagenURL      string  `json:"imagen"`
	Activo         bool    `json:"activo"`
}

// TableName especifica la tabla física para FacturaDetalle
func (FacturaDetalle) TableName() string {
	return "factura_detalle"
}

// Factura representa una venta realizada a un cliente
type Factura struct {
	ID         int64            `json:"id" gorm:"primaryKey"`
	PersonaID  int64            `json:"persona_id"`
	Fecha      time.Time        `json:"fecha"`
	Total      float64          `json:"total"`
	MetodoPago string           `json:"metodo_pago"`
	Detalle    []FacturaDetalle `json:"detalle" gorm:"foreignKey:FacturaID"`
	Persona    Persona          `json:"persona" gorm:"foreignKey:PersonaID;references:Id"`
}

// TableName especifica la tabla física para Factura
func (Factura) TableName() string {
	return "factura"
}

// MigrarFactura crea o actualiza las tablas factura y factura_detalle en la base de datos
func MigrarFactura() {
	db.BaseDeDatos.AutoMigrate(&Factura{}, &FacturaDetalle{})
}
