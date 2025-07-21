package dominios

import "srvr/db"

// Animal representa los tipos de mascota
type Animal struct {
	ID          int64  `json:"id" gorm:"primaryKey"`
	Nombre      string `json:"nombre"`
	CategoriaID int64  `json:"categoria_id"`
}

type Animales []Animal

// TableName especifica el nombre físico de la tabla
func (Animal) TableName() string {
	return "animal"
}

// MigrarAnimal crea la tabla si no existe
func MigrarAnimal() {
	db.BaseDeDatos.AutoMigrate(Animal{})
}
