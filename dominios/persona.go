package dominios

import (
	"srvr/db"
	"time"
)

type Persona struct {
	Id          int64     `json:"id" gorm:"primaryKey;column:prsn__id"`
	Cedula      string    `json:"cedula" gorm:"column:prsncdla"`
	Nombre      string    `json:"nombre" gorm:"column:prsnnmbr"`
	Apellido    string    `json:"apellido" gorm:"column:prsnapll"`
	FechaInicio time.Time `json:"fechaInicio" gorm:"column:prsnfcin"`
	FechaFin    time.Time `json:"fechaFin" gorm:"column:prsnfcfn"`
	Mail        string    `json:"mail" gorm:"column:prsnmail"`
	Login       string    `json:"login" gorm:"column:prsnlogn"`
	Password    string    `json:"password" gorm:"column:prsnpass"`
	Activo      uint8     `json:"activo" gorm:"column:prsnactv"`
	Telefono    string    `json:"telefono" gorm:"column:prsntelf"`
	Sexo        string    `json:"sexo" gorm:"column:prsnsexo"`
	Direccion   string    `json:"direccion" gorm:"column:prsndire"`
	
}

type Personas []Persona

func (Persona) TableName() string {
	return "prsn"
}

func MigrarPersona() {
	//	db.BaseDeDatos.AutoMigrate(User{})
	db.BaseDeDatos.AutoMigrate(Persona{})
}
