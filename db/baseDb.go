package db

import (
	"database/sql"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var BaseDeDatos *gorm.DB
var conexionDb *sql.DB

func crearBaseDeDatosLocalSiNoExiste(dsnProduccion string) {
	if dsnProduccion != "" {
		return // En producción (Render) la DB ya viene creada y configurada
	}

	// Conexión a la base de datos 'postgres' por defecto para administrar
	dsnDefault := "host=127.0.0.1 user=postgres password=12345 dbname=postgres port=5432 sslmode=disable"
	dbDefault, err := gorm.Open(postgres.Open(dsnDefault), &gorm.Config{})
	if err != nil {
		fmt.Println("No se pudo conectar a PostgreSQL para verificar/crear la BD:", err)
		return
	}

	sqlDB, _ := dbDefault.DB()
	defer sqlDB.Close()

	var exists bool
	dbDefault.Raw("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = 'srvr_go')").Scan(&exists)
	if !exists {
		fmt.Println("La base de datos 'srvr_go' no existe. Creándola...")
		if err := dbDefault.Exec("CREATE DATABASE srvr_go").Error; err != nil {
			fmt.Println("Error creando la BD:", err)
		} else {
			fmt.Println("Base de datos 'srvr_go' creada exitosamente.")
		}
	}
}

func InitDB() {
	dsn := os.Getenv("DATABASE_URL")

	crearBaseDeDatosLocalSiNoExiste(dsn)

	if dsn == "" {
		dsn = "host=127.0.0.1 user=postgres password=12345 dbname=srvr_go port=5432 sslmode=disable"
		fmt.Println("Conectando a desarrollo local...")
	} else {
		fmt.Println("Conectando a producción (Render)...")
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		PrepareStmt: false,
	})

	if err != nil {
		panic("Error al conectar a la base de datos: " + err.Error())
	}

	BaseDeDatos = db
	fmt.Println("Conexión exitosa con GORM")
}
