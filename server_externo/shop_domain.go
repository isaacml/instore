package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
	"time"
)

func recoger_dominio(w http.ResponseWriter, r *http.Request) {
	arr_domain := libs.DomainGenerator(r.FormValue("dominio"))
	fecha_actual := time.Now()
	string_fecha := fmt.Sprintf("%4d%02d%02d", fecha_actual.Year(), int(fecha_actual.Month()), fecha_actual.Day())
	fmt.Println(string_fecha)
	for _, val := range arr_domain {
		query, err := db.Query("SELECT fichero FROM publi WHERE destino = ? AND fecha_inicio = ?", val, string_fecha)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			var file string
			err = query.Scan(&file)
			if err != nil {
				Error.Println(err)
			}
			fmt.Fprint(w, file)
		}
	}
}
