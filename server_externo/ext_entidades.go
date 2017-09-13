package main

import (
	"fmt"
	"net/http"
	"time"
)

//GESTION DE ENTIDADES (entidades.html)
func entidades(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	accion := r.FormValue("accion")
	//DAR DE ALTAR UNA NUEVA ENTIDAD
	if accion == "entidad" {
		username := r.FormValue("username")
		entidad := r.FormValue("entidad")
		var output string
		if entidad == "" {
			output = "<div class='form-group text-warning'>El campo no puede estar vacio</div>"
		} else {
			query, err := db.Query("SELECT id, padre_id FROM usuarios WHERE user = ?", username)
			if err != nil {
				Error.Println(err)
			}
			for query.Next() {
				var id, padre_id, id_admin int
				err = query.Scan(&id, &padre_id)
				if err != nil {
					Error.Println(err)
				}
				//Hacemos un select para obtener el id del usuario super-admin
				err = db.QueryRow("SELECT id FROM usuarios WHERE padre_id = 0 AND entidad_id = 0").Scan(&id_admin)
				if err != nil {
					Error.Println(err)
				}
				//Si es un usuario super-admin o un usuario que tiene creador super-admin, le permitimos crear entidades
				if padre_id == 0 || padre_id == id_admin {
					timestamp := time.Now().Unix()
					db_mu.Lock()
					_, err1 := db.Exec("INSERT INTO entidades (`nombre`, `creador_id`, `timestamp`, `last_access`) VALUES (?, ?, ?, ?)", entidad, id, timestamp, timestamp)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						output = "<div class='form-group text-danger'>Fallo al añadir entidad</div>"
					} else {
						output = "OK"
					}
				} else {
					output = "<div class='form-group text-danger'>Solo un usuario ROOT puede añadir una entidad</div>"
				}
			}
		}
		fmt.Fprint(w, output)
	}
	//MODIFICAR / EDITAR UNA ENTIDAD
	if accion == "edit_entidad" {
		edit_id := r.FormValue("edit_id")
		username := r.FormValue("username")
		entidad := r.FormValue("entidad")
		if entidad == "" {
			empty = "El campo no puede estar vacío"
			fmt.Fprintf(w, "<div class='form-group text-warning'>%s</div>", empty)
		} else {
			query, err := db.Query("SELECT id, padre_id FROM usuarios WHERE user = ?", username)
			if err != nil {
				Error.Println(err)
			}
			for query.Next() {
				var id, padre_id int
				err = query.Scan(&id, &padre_id)
				if err != nil {
					Error.Println(err)
				}
				if padre_id == 0 || padre_id == 1 {
					db_mu.Lock()
					_, err1 := db.Exec("UPDATE entidades SET nombre=? WHERE id = ?", entidad, edit_id)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						bad = "Fallo al modificar entidad"
						fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
					} else {
						fmt.Fprint(w, "OK")
					}
				} else {
					bad = "Solo un usuario ROOT puede editar una entidad"
					fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
				}
			}
		}
	}
	//MOSTRAR ENTIDADES EN UNA TABLA
	if accion == "tabla_entidad" {
		username := r.FormValue("username")
		query, err := db.Query("SELECT id FROM usuarios WHERE user = ?", username)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			var id, creador_id int
			var tiempo int64
			var nombre string
			err = query.Scan(&creador_id)
			if err != nil {
				Error.Println(err)
			}
			query, err := db.Query("SELECT id, nombre, timestamp FROM entidades WHERE creador_id = ?", creador_id)
			if err != nil {
				Warning.Println(err)
			}
			for query.Next() {
				err = query.Scan(&id, &nombre, &tiempo)
				if err != nil {
					Error.Println(err)
				}
				creacion := time.Unix(tiempo, 0)
				fmt.Fprintf(w, "<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Pulsa para editar entidad'>%s</a></td><td>%s</td></tr>",
					id, nombre, creacion)
			}
		}
	}
	//CARGA LOS DATOS DE ENTIDAD EN UN FORMULARIO
	if accion == "load_entidad" {
		edit_id := r.FormValue("edit_id")
		var id int
		var nombre string
		query, err := db.Query("SELECT id, nombre FROM entidades WHERE id = ?", edit_id)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			err = query.Scan(&id, &nombre)
			if err != nil {
				Error.Println(err)
			}
			fmt.Fprintf(w, "id=%d&entidad=%s", id, nombre)
		}
	}
}
