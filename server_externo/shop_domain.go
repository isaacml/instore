package main

import (
	"fmt"
	//"github.com/isaacml/instore/libs"
	"net/http"
)

func recoger_dominio(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.FormValue("dominio"))

	Q1, err := db.Query("SELECT fichero FROM publi WHERE destino = ?", r.FormValue("dominio"))
	if err != nil {
		Error.Println(err)
	}
	for Q1.Next() {
		var file string
		err = Q1.Scan(&file)
		if err != nil {
			Error.Println(err)
		}
		fmt.Fprint(w, file)
	}
}
