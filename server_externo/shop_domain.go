package main

import (
	"fmt"
	//"github.com/isaacml/instore/libs"
	"net/http"
)

func recoger_dominio(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.FormValue("dominio"))

	query, err := db.Query("SELECT id FROM usuarios WHERE user = ?", r.FormValue("userAdmin"))
	if err != nil {
		Error.Println(err)
	}
}
