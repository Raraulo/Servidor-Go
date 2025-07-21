package utiles

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type retorna struct {
	Registro interface{}
	Ok       bool `json:"ok"`
}

func EnviaData(rw http.ResponseWriter, data interface{}, estado int) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(estado)
	fmt.Println(">>>data:", data)
	objetoRetorna := retorna{data, true}
	// fmt.Println("retorna:", objetoRetorna)
	// salida, _ := json.Marshal(&data)
	salida, _ := json.Marshal(&objetoRetorna)
	// fmt.Println(salida)
	fmt.Fprintln(rw, string(salida))
}

func EnviaError(rw http.ResponseWriter, estado int) {
	rw.WriteHeader(estado)
	objetoRetorna := retorna{nil, false}
	salida, _ := json.Marshal(&objetoRetorna)
	// fmt.Println(salida)
	fmt.Fprintln(rw, string(salida))
}
