package seguridad

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"srvr/db"
	"srvr/dominios"

	"github.com/gorilla/mux"
)

// ListarFacturas devuelve todas las facturas registradas (solo datos básicos)
func ListarFacturas(w http.ResponseWriter, r *http.Request) {
	var facturas []dominios.Factura
	if err := db.BaseDeDatos.Find(&facturas).Error; err != nil {
		http.Error(w, `{"ok":false,"msg":"Error al obtener facturas"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":       true,
		"facturas": facturas,
	})
}

// ObtenerFactura devuelve una sola factura por su ID incluyendo detalles
func ObtenerFactura(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"ok":false,"msg":"ID inválido"}`, http.StatusBadRequest)
		return
	}

	var factura dominios.Factura
	if err := db.BaseDeDatos.
		Preload("Detalle").
		Preload("Persona"). // ✅ Preload para obtener nombre, apellido, etc.
		First(&factura, id).
		Error; err != nil {
		http.Error(w, `{"ok":false,"msg":"Factura no encontrada"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":      true,
		"factura": factura,
	})
}

// CrearFactura crea una nueva factura básica (sin detalles)
func CrearFactura(w http.ResponseWriter, r *http.Request) {
	var nueva dominios.Factura
	if err := json.NewDecoder(r.Body).Decode(&nueva); err != nil {
		http.Error(w, `{"ok":false,"msg":"Error al leer los datos"}`, http.StatusBadRequest)
		return
	}
	if nueva.PersonaID == 0 || nueva.Total <= 0 || nueva.MetodoPago == "" {
		http.Error(w, `{"ok":false,"msg":"Datos incompletos"}`, http.StatusBadRequest)
		return
	}
	if err := db.BaseDeDatos.Create(&nueva).Error; err != nil {
		http.Error(w, `{"ok":false,"msg":"Error al guardar"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":      true,
		"factura": nueva,
	})
}

// CrearFacturaConDetalle crea una factura con sus líneas, actualiza stock y disponibilidad

func CrearFacturaConDetalle(w http.ResponseWriter, r *http.Request) {
	type EntradaDetalle struct {
		ProductoID int64 `json:"producto_id"`
		Cantidad   int   `json:"cantidad"`
	}

	var entrada struct {
		PersonaID  int64            `json:"persona_id"`
		MetodoPago string           `json:"metodo_pago"`
		Detalles   []EntradaDetalle `json:"detalles"`
	}

	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		http.Error(w, `{"ok":false,"msg":"JSON inválido"}`, http.StatusBadRequest)
		return
	}
	if entrada.PersonaID == 0 || len(entrada.Detalles) == 0 {
		http.Error(w, `{"ok":false,"msg":"Datos incompletos"}`, http.StatusBadRequest)
		return
	}

	var total float64
	productos := make(map[int64]dominios.Producto)

	for _, item := range entrada.Detalles {
		var prod dominios.Producto
		if err := db.BaseDeDatos.First(&prod, item.ProductoID).Error; err != nil {
			http.Error(w, `{"ok":false,"msg":"Producto no encontrado"}`, http.StatusNotFound)
			return
		}
		if prod.Stock < item.Cantidad {
			http.Error(w, `{"ok":false,"msg":"Stock insuficiente para `+prod.Nombre+`"}`, http.StatusConflict)
			return
		}
		productos[item.ProductoID] = prod
		total += float64(item.Cantidad) * prod.Precio
	}

	factura := dominios.Factura{
		PersonaID:  entrada.PersonaID,
		MetodoPago: entrada.MetodoPago,
		Total:      total,
		Fecha:      time.Now(),
	}
	if err := db.BaseDeDatos.Create(&factura).Error; err != nil {
		http.Error(w, `{"ok":false,"msg":"Error al crear factura"}`, http.StatusInternalServerError)
		return
	}

	for _, item := range entrada.Detalles {
		prod := productos[item.ProductoID]

		det := dominios.FacturaDetalle{
			FacturaID:      factura.ID,
			ProductoID:     item.ProductoID,
			Cantidad:       item.Cantidad,
			PrecioUnitario: prod.Precio,
			ProductoNombre: prod.Nombre,    // ✅ Asegurado
			ImagenURL:      prod.ImagenURL, // ✅ Asegurado
			Activo:         true,           // ✅ Activo por defecto
		}

		if err := db.BaseDeDatos.Create(&det).Error; err != nil {
			http.Error(w, `{"ok":false,"msg":"Error al guardar detalle"}`, http.StatusInternalServerError)
			return
		}

		prod.Stock -= item.Cantidad
		if prod.Stock <= 0 {
			prod.Stock = 0
			prod.Disponible = false
		}
		if err := db.BaseDeDatos.Save(&prod).Error; err != nil {
			http.Error(w, `{"ok":false,"msg":"Error al actualizar producto"}`, http.StatusInternalServerError)
			return
		}
	}

	db.BaseDeDatos.Preload("Detalle").First(&factura, factura.ID)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":      true,
		"factura": factura,
	})
}

// ActualizarFactura modifica los datos de una factura existente
func ActualizarFactura(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"ok":false,"msg":"ID inválido"}`, http.StatusBadRequest)
		return
	}
	var actualizada dominios.Factura
	if err := json.NewDecoder(r.Body).Decode(&actualizada); err != nil {
		http.Error(w, `{"ok":false,"msg":"Error al leer los datos"}`, http.StatusBadRequest)
		return
	}
	actualizada.ID = id
	if err := db.BaseDeDatos.Model(&dominios.Factura{}).Where("id = ?", id).Updates(actualizada).Error; err != nil {
		http.Error(w, `{"ok":false,"msg":"Error al actualizar"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":      true,
		"factura": actualizada,
	})
}

// EliminarFactura elimina una factura existente
func EliminarFactura(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"ok":false,"msg":"ID inválido"}`, http.StatusBadRequest)
		return
	}
	if err := db.BaseDeDatos.Delete(&dominios.Factura{}, id).Error; err != nil {
		http.Error(w, `{"ok":false,"msg":"Error al eliminar"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":  true,
		"msg": "Factura eliminada correctamente",
	})
}

// VerFacturasConDetallePorUsuario muestra facturas con detalle (JOIN explícito)
func VerFacturasConDetallePorUsuario(w http.ResponseWriter, r *http.Request) {
	personaIDStr := r.Header.Get("PersonaID")
	personaID, err := strconv.ParseInt(personaIDStr, 10, 64)
	if err != nil || personaID == 0 {
		http.Error(w, `{"ok":false,"msg":"ID de usuario no válido"}`, http.StatusBadRequest)
		return
	}

	// 1️⃣ Obtén tus facturas básicas
	var facturas []dominios.Factura
	if err := db.BaseDeDatos.
		Where("persona_id = ?", personaID).
		Order("fecha DESC").
		Find(&facturas).Error; err != nil {
		http.Error(w, `{"ok":false,"msg":"Error al obtener facturas"}`, http.StatusInternalServerError)
		return
	}

	// 2️⃣ Estructuras de salida
	type detalleOut struct {
		ProductoNombre string  `json:"producto_nombre"`
		Cantidad       int     `json:"cantidad"`
		PrecioUnitario float64 `json:"precio_unitario"`
		ImagenURL      string  `json:"imagen"`
	}
	type facturaOut struct {
		ID         int64        `json:"id"`
		PersonaID  int64        `json:"persona_id"`
		Fecha      time.Time    `json:"fecha"`
		Total      float64      `json:"total"`
		MetodoPago string       `json:"metodo_pago"`
		Detalle    []detalleOut `json:"detalle"`
	}

	// 3️⃣ JOIN con la columna correcta
	var salida []facturaOut
	for _, f := range facturas {
		var dets []detalleOut
		db.BaseDeDatos.
			Table("factura_detalle fd").
			Select(`
                p.nombre            AS producto_nombre,
                fd.cantidad         AS cantidad,
                fd.precio_unitario  AS precio_unitario,
                p.imagen_url        AS imagen_url`).
			Joins("JOIN producto p ON p.id = fd.producto_id").
			Where("fd.factura_id = ?", f.ID).
			Scan(&dets)

		salida = append(salida, facturaOut{
			ID:         f.ID,
			PersonaID:  f.PersonaID,
			Fecha:      f.Fecha,
			Total:      f.Total,
			MetodoPago: f.MetodoPago,
			Detalle:    dets,
		})
	}

	// 4️⃣ Envía la respuesta
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":       true,
		"facturas": salida,
		"persona":  personaID,
		"usuario":  r.Header.Get("Usuario"),
	})
}

// VerFacturasPorUsuario permite ver facturas sin detalle
func VerFacturasPorUsuario(w http.ResponseWriter, r *http.Request) {
	personaIDHeader := r.Header.Get("PersonaID")
	username := r.Header.Get("Usuario")
	paramID := mux.Vars(r)["persona_id"]

	personaIDToken, err := strconv.ParseInt(personaIDHeader, 10, 64)
	if err != nil || personaIDToken == 0 {
		http.Error(w, `{"ok":false,"msg":"ID en token inválido"}`, http.StatusUnauthorized)
		return
	}

	var personaIDFinal int64
	if username == "admin1" && paramID != "" {
		personaIDFinal, err = strconv.ParseInt(paramID, 10, 64)
		if err != nil || personaIDFinal == 0 {
			http.Error(w, `{"ok":false,"msg":"persona_id inválido"}`, http.StatusBadRequest)
			return
		}
	} else {
		personaIDFinal = personaIDToken
	}

	var facturas []dominios.Factura
	if err := db.BaseDeDatos.Where("persona_id = ?", personaIDFinal).Find(&facturas).Error; err != nil {
		http.Error(w, `{"ok":false,"msg":"Error al obtener facturas"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":       true,
		"usuario":  username,
		"persona":  personaIDFinal,
		"facturas": facturas,
	})

}
