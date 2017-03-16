package main

import (
	"fmt"
	"net/http"
)

func config_shop(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	user := r.FormValue("username")
	query, err := db.Query("SELECT entidad_id FROM usuarios WHERE user = ?", user)
	if err != nil {
		Error.Println(err)
	}
	for query.Next() {
		var entidad int
		err = query.Scan(&entidad)
		if err != nil {
			Error.Println(err)
		}
		var list string
		query, err := db.Query("SELECT id, nombre FROM entidades WHERE id=?", entidad)
		if err != nil {
			Error.Println(err)
		}
		list = "<div class='panel-heading'>Entidad</div><div class='panel-body'><select id='entidad' name='entidad'>"
		for query.Next() {
			var id int
			var name string
			err = query.Scan(&id, &name)
			if err != nil {
				Error.Println(err)
			}
			list += fmt.Sprintf("<option value='%d'>%s</option>", id, name)
		}
		list += "</select></div>"
		fmt.Fprint(w, list)
	}
}
