package dominios

import "srvr/db"

// Categoria representa la tabla de categorías generales
type Categoria struct {
	ID     int64  `json:"id" gorm:"primaryKey"`
	Nombre string `json:"nombre"`
}

type Categorias []Categoria

// TableName especifica el nombre físico de la tabla
func (Categoria) TableName() string {
	return "categoria"
}

// MigrarCategoria crea la tabla si no existe
func MigrarCategoria() {
	db.BaseDeDatos.AutoMigrate(Categoria{})
}
