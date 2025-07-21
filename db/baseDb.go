package db

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dsn = `host=127.0.0.1 user=postgres password=12345 dbname=srvr_go port=5432 sslmode=disable`

var conexionDb *sql.DB

// función anónima que se autoejecuta para conectarse a la BBDD
var BaseDeDatos = func() (db *gorm.DB) {
	// if db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
	if db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		fmt.Println("Error al conectarse a la BBDD")
		panic(err)
	} else {
		fmt.Println("Conexión existosa con GORM")
		return db
	}
}()
