package main

import (
	"fmt"
	"net/http"
)

func organizaciones(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	fmt.Println(username)
	var creador_id, ent_id int
	var ent_name string
	query1, err := db.Query("SELECT id FROM usuarios WHERE user = ?", username)
	if err != nil {
		Error.Println(err)
	}
	for query1.Next() {
		err = query1.Scan(&creador_id)
		if err != nil {
			Error.Println(err)
		}
		query2, err := db.Query("SELECT id, nombre FROM entidades WHERE creador_id = ?", creador_id)
		if err != nil {
			Warning.Println(err)
		}
		for query2.Next() {
			err = query2.Scan(&ent_id, &ent_name)
			if err != nil {
				Error.Println(err)
			}
			fmt.Println(ent_id, ent_name)
			fmt.Fprintf(w, "<input type='text' value='%s' readonly><input type='text' name='almacenes%d'><br>", ent_name, ent_id)
		}
	}
}
