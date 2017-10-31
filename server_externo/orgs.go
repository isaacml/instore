package main

import (
	"fmt"
	"net/http"
)

func organizaciones(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	accion := r.FormValue("accion")
	username := r.FormValue("username")
	if accion == "ent_org" {
		var creador_id, ent_id int
		var ent_name, list string
		err := db.QueryRow("SELECT id FROM usuarios WHERE user = ?", username).Scan(&creador_id)
		if err != nil {
			Error.Println(err)
		}
		query2, err := db.Query("SELECT id, nombre FROM entidades WHERE creador_id = ?", creador_id)
		if err != nil {
			Warning.Println(err)
		}
		list = "<div class='panel-heading'>Entidad</div><select id='entidad' name='entidad'>"
		for query2.Next() {
			err = query2.Scan(&ent_id, &ent_name)
			if err != nil {
				Error.Println(err)
			}
			list += fmt.Sprintf("<option value='%d'>%s</option>", ent_id, ent_name)
		}
		list += "</select>"
		fmt.Fprint(w, list)
	}
	if accion == "alm_org" {
		var creador_id, ent_id int
		var ent_name, list string
		err := db.QueryRow("SELECT id FROM usuarios WHERE user = ?", username).Scan(&creador_id)
		if err != nil {
			Error.Println(err)
		}
		query2, err := db.Query("SELECT id, almacen FROM almacenes WHERE creador_id = ?", creador_id)
		if err != nil {
			Warning.Println(err)
		}
		list = "<div class='panel-heading'>Almacen</div><select id='almacen' name='almacen'>"
		for query2.Next() {
			err = query2.Scan(&ent_id, &ent_name)
			if err != nil {
				Error.Println(err)
			}
			list += fmt.Sprintf("<option value='%d'>%s</option>", ent_id, ent_name)
		}
		list += "</select>"
		fmt.Fprint(w, list)
	}
	if accion == "pais_org" {

	}
	if accion == "reg_org" {

	}
}
