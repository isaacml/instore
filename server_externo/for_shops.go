package main

import (
	"fmt"
	"net/http"
)

func config_shop(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	accion := r.FormValue("action")
	
	if accion == "entidad" {
		user := r.FormValue("username")
		query, err := db.Query("SELECT entidad_id FROM usuarios WHERE user = ?", user)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			var entidad, id_ent int
			var list, name string
			err = query.Scan(&entidad)
			if err != nil {
				Error.Println(err)
			}
			query, err := db.Query("SELECT id, nombre FROM entidades WHERE id=?", entidad)
			if err != nil {
				Error.Println(err)
			}
			if query.Next() {
				  list = "<div class='panel-heading'>Entidad</div><div class='panel-body'><select id='entidad' name='entidad'><option value='' selected>Selecciona una entidad</option>"
			      query.Scan(&id_ent, &name)
			      list += fmt.Sprintf("<option value='%d'>%s</option>", id_ent, name)
			      for query.Next() {
			          query.Scan(&id_ent, &name)
				      if err != nil {
						Error.Println(err)
					  }
					  list += fmt.Sprintf("<option value='%d'>%s</option>", id_ent, name)
			      }
			      list += "</select></div>"
			} else {
			      list = "<div class='panel-heading'>Entidad</div><div class='panel-body'><select id='entidad' name='entidad'><option value='' selected>No hay entidades</option></select></div>"
			}
			fmt.Fprint(w, list)
		}
	}
	if accion == "almacen" {
		ent := r.FormValue("entidad")
		fmt.Println(ent)
	}
}
