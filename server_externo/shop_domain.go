package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
	"time"
)

func recoger_dominio(w http.ResponseWriter, r *http.Request) {
	var output string
	arr_domain := libs.DomainGenerator(r.FormValue("dominio"))
	fecha_actual := time.Now()
	//Formato de la fecha actual --> 20070405
	string_fecha := fmt.Sprintf("%4d%02d%02d", fecha_actual.Year(), int(fecha_actual.Month()), fecha_actual.Day())

	output += "[publi]"
	for _, val := range arr_domain {
		publicidad, err := db.Query("SELECT fichero FROM publi WHERE destino = ? AND fecha_inicio = ?", val, string_fecha)
		if err != nil {
			Error.Println(err)
		}
		for publicidad.Next() {
			var f_publi string
			err = publicidad.Scan(&f_publi)
			if err != nil {
				Error.Println(err)
			}
			output += ";" + f_publi
		}
	}
	output += "[mensaje]"
	for _, val := range arr_domain {
		mensajes, err := db.Query("SELECT fichero, playtime FROM mensaje WHERE destino = ? AND fecha_inicio = ?", val, string_fecha)
		if err != nil {
			Error.Println(err)
		}
		for mensajes.Next() {
			var f_msg, playtime string
			err = mensajes.Scan(&f_msg, &playtime)
			if err != nil {
				Error.Println(err)
			}
			output += ";" + f_msg + "<=>" + playtime
		}
	}
	fmt.Fprint(w, output)
}
