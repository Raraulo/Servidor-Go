package seguridad

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"srvr/db"
	"srvr/dominios"
	"srvr/handlers/utiles"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// Obtener todos los productos
func GetProductos(w http.ResponseWriter, r *http.Request) {
	var productos []dominios.Producto
	err := db.BaseDeDatos.Find(&productos).Error
	if err != nil {
		utiles.EnviaError(w, http.StatusInternalServerError)
		return
	}
	utiles.EnviaData(w, productos, http.StatusOK)
}

// Obtener un solo producto
func GetProducto(w http.ResponseWriter, r *http.Request) {
	if prod, err := GetProductoById(r); err != nil {
		utiles.EnviaError(w, http.StatusNotFound)
	} else {
		utiles.EnviaData(w, prod, http.StatusOK)
	}
}

// Función auxiliar para buscar producto por ID
func GetProductoById(r *http.Request) (dominios.Producto, *gorm.DB) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	var result = struct{ Id int }{}
	producto := dominios.Producto{}

	db.BaseDeDatos.Raw(`SELECT id AS Id FROM producto WHERE id = ?`, id).Scan(&result)
	if err := db.BaseDeDatos.First(&producto, id); err.Error != nil {
		fmt.Println("Error al obtener producto")
		return producto, err
	}
	return producto, nil
}
func GetCategorias(w http.ResponseWriter, r *http.Request) {
	var categorias []dominios.Categoria
	db.BaseDeDatos.Find(&categorias)
	json.NewEncoder(w).Encode(categorias)
}

// Crear un nuevo producto
func CreateProducto(w http.ResponseWriter, r *http.Request) {
	p := dominios.Producto{}
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		utiles.EnviaError(w, http.StatusBadRequest)
		return
	}
	if err := db.BaseDeDatos.Create(&p).Error; err != nil {
		utiles.EnviaError(w, http.StatusInternalServerError)
		return
	}
	utiles.EnviaData(w, p, http.StatusCreated)
}

// Actualizar un producto existente
func UpdateProducto(w http.ResponseWriter, r *http.Request) {
	var prodID int64
	if oldProd, err := GetProductoById(r); err != nil {
		utiles.EnviaError(w, http.StatusNotFound)
	} else {
		prodID = oldProd.ID
		var p dominios.Producto
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			fmt.Println("Error al decodificar JSON:", err)
			utiles.EnviaError(w, http.StatusUnprocessableEntity)
			return
		}

		// Log de datos recibidos
		fmt.Printf("JSON recibido para actualizar producto: %+v\n", p)

		p.ID = prodID
		result := db.BaseDeDatos.Save(&p)

		if result.Error != nil {
			fmt.Println("Error al actualizar:", result.Error)
			utiles.EnviaError(w, http.StatusInternalServerError)
			return
		}

		fmt.Println("Producto actualizado correctamente")
		utiles.EnviaData(w, p, http.StatusOK)
	}
}

// Eliminar un producto por ID
func DeleteProducto(w http.ResponseWriter, r *http.Request) {
	if prod, err := GetProductoById(r); err != nil {
		utiles.EnviaError(w, http.StatusNotFound)
	} else {
		if err := db.BaseDeDatos.Delete(&prod).Error; err != nil {
			utiles.EnviaError(w, http.StatusInternalServerError)
			return
		}
		utiles.EnviaData(w, prod, http.StatusOK)
	}
}
