package dominios

import (
	"srvr/db"
	"time"
)

// Venta representa una venta realizada a un cliente (equivalente a factura)
type Venta struct {
	ID         int64     `json:"id" gorm:"primaryKey;column:vnta__id"`
	PersonaID  int64     `json:"persona_id" gorm:"column:vntprsnid"`
	Fecha      time.Time `json:"fecha" gorm:"column:vntfc"`
	Total      float64   `json:"total" gorm:"column:vnttotal"`
	MetodoPago string    `json:"metodo_pago" gorm:"column:vntmtdpago"`
}

type Ventas []Venta

// TableName especifica el nombre físico de la tabla
func (Venta) TableName() string {
	return "venta"
}

// MigrarVenta crea la tabla si no existe
func MigrarVenta() {
	db.BaseDeDatos.AutoMigrate(Venta{})
}
