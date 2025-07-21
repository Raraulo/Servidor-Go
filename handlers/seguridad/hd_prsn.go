package seguridad

import (
	"encoding/json"
	"fmt"
	"net/http"
	"srvr/db"
	"srvr/dominios"
	"srvr/handlers/utiles"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func limpiarFechas(p *dominios.Persona) {
	if p.FechaInicio.Year() == 0 {
		p.FechaInicio = time.Time{}
	}
	if p.FechaFin.Year() == 0 {
		p.FechaFin = time.Time{}
	}
}

func GetPersonas(rw http.ResponseWriter, r *http.Request) {
	personas := []dominios.Persona{}
	db.BaseDeDatos.Find(&personas)

	// Limpiar fechas inválidas
	for i := range personas {
		limpiarFechas(&personas[i])
	}

	utiles.EnviaData(rw, personas, http.StatusOK)
}

func GetPersona(rw http.ResponseWriter, r *http.Request) {
	if persona, err := GetPersonaById(r); err != nil {
		utiles.EnviaError(rw, http.StatusNotFound)
	} else {
		limpiarFechas(&persona)
		utiles.EnviaData(rw, persona, http.StatusOK)
	}
}

func GetPersonaById(r *http.Request) (dominios.Persona, *gorm.DB) {
	vars := mux.Vars(r)
	usroId, _ := strconv.Atoi(vars["id"])
	var result = struct {
		Id int
	}{}

	persona := dominios.Persona{}

	db.BaseDeDatos.Raw(`select prsn__id Id from prsn where prsn__id = ?`, usroId).Scan(&result)
	fmt.Println("salida:", result.Id)
	if err := db.BaseDeDatos.First(&persona, usroId); err.Error != nil {
		fmt.Println("Error al obtener persona")
		return persona, err
	} else {
		return persona, nil
	}
}

func GetUsuario(id int64) dominios.Persona {
	prsn := dominios.Persona{}
	db.BaseDeDatos.First(&prsn, id)
	limpiarFechas(&prsn)
	return prsn
}

func ExisteUsuario(login, password string) dominios.Persona {
	prsn := dominios.Persona{}

	// Usa nombres de columnas reales: prsnlogn y prsnpass
	err := db.BaseDeDatos.Where("prsnlogn = ? AND prsnpass = ?", login, password).First(&prsn).Error
	if err != nil {
		fmt.Println("Login fallido para:", login, password)
	} else {
		fmt.Println("Login exitoso:", prsn.Id, prsn.Login)
	}

	return prsn
}

func CreaPersona(rw http.ResponseWriter, r *http.Request) {
	persona := dominios.Persona{}
	decoder := json.NewDecoder(r.Body)
	fmt.Println(r.Body)

	if err := decoder.Decode(&persona); err != nil {
		utiles.EnviaError(rw, http.StatusUnprocessableEntity)
	} else {
		if persona.Id == 0 && persona.Nombre != "" {
			db.BaseDeDatos.Save(&persona)
			limpiarFechas(&persona)
			utiles.EnviaData(rw, persona, http.StatusCreated)
		} else {
			utiles.EnviaError(rw, http.StatusUnprocessableEntity)
		}
	}
}

func UpdatePersona(rw http.ResponseWriter, r *http.Request) {
	var usro_id int64

	if persona_antiguo, err := GetPersonaById(r); err != nil {
		utiles.EnviaError(rw, http.StatusNotFound)
	} else {
		usro_id = persona_antiguo.Id
		persona := dominios.Persona{}
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&persona); err != nil {
			utiles.EnviaError(rw, http.StatusUnprocessableEntity)
		} else {
			persona.Id = usro_id
			db.BaseDeDatos.Save(&persona)
			limpiarFechas(&persona)
			utiles.EnviaData(rw, persona, http.StatusOK)
		}
	}
}

func DeletePersona(rw http.ResponseWriter, r *http.Request) {
	if usro, err := GetPersonaById(r); err != nil {
		utiles.EnviaError(rw, http.StatusNotFound)
	} else {
		db.BaseDeDatos.Delete(&usro)
		utiles.EnviaData(rw, usro, http.StatusOK)
	}
}

// ----------- LOGIN (con debug extra) -------------
func Login(rw http.ResponseWriter, r *http.Request) {
	type Credenciales struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	var creds Credenciales
	err := json.NewDecoder(r.Body).Decode(&creds)

	// ----------- DEBUG: Muestra lo recibido -----------
	fmt.Println("DEBUG Login recibido -> Login:", creds.Login)
	fmt.Println("DEBUG Login recibido -> Password:", creds.Password)
	// ---------------------------------------------------

	if err != nil {
		http.Error(rw, "Error al procesar entrada", http.StatusUnprocessableEntity)
		return
	}

	if creds.Login == "" || creds.Password == "" {
		fmt.Println("❌ ERROR: Login o Password vacío (verifica tu JSON y el tag de la struct)")
	}

	persona := ExisteUsuario(creds.Login, creds.Password)

	if persona.Id != 0 {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"ok":       true,
			"Registro": persona,
		})
	} else {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"ok":  false,
			"msg": "Login o password incorrectos",
		})
	}
}
