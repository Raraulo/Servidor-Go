package db

import (
	"database/sql"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dsn = os.Getenv("DATABASE_URL")

var conexionDb *sql.DB

// función anónima que se autoejecuta para conectarse a la BBDD
var BaseDeDatos = func() (db *gorm.DB) {
	if dsn == "" {
		dsn = "host=127.0.0.1 user=postgres password=12345 dbname=srvr_go port=5432 sslmode=disable"
	}
	if db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		fmt.Println("Error al conectarse a la BBDD")
		panic(err)
	} else {
		fmt.Println("Conexión existosa con GORM")
		return db
	}
}()
