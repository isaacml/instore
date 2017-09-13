package main

import (
	"fmt"
	"net/http"
	"time"
)

//GESTION DE ALMACENES (almacenes.html)
func almacenes(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	accion := r.FormValue("accion")
	//DAR DE ALTAR UN NUEVO ALMACEN
	if accion == "almacen" {
		username := r.FormValue("username")
		almacen := r.FormValue("almacen")
		entidad := r.FormValue("entidad")
		var output string
		if almacen == "" {
			output = "<div class='form-group text-warning'>El campo no puede estar vacio</div>"
		} else if entidad == "" {
			output = "<div class='form-group text-warning'>Debe haber almenos una entidad</div>"
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
				//Si es un usuario super-admin o un usuario que tiene creador super-admin, le permitimos crear almacenes
				if padre_id == 0 || padre_id == id_admin {
					timestamp := time.Now().Unix()
					db_mu.Lock()
					_, err1 := db.Exec("INSERT INTO almacenes (`almacen`, `creador_id`, `timestamp`, `entidad_id`) VALUES (?, ?, ?, ?)", almacen, id, timestamp, entidad)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						output = "<div class='form-group text-danger'>Fallo al añadir almacen</div>"
					} else {
						output = "OK"
					}
				} else {
					output = "<div class='form-group text-danger'>Solo un usuario ROOT puede añadir un almacen</div>"
				}
			}
		}
		fmt.Fprint(w, output)
	}
	//MODIFICAR / EDITAR UN ALMACEN
	if accion == "edit_almacen" {
		edit_id := r.FormValue("edit_id")
		username := r.FormValue("username")
		almacen := r.FormValue("almacen")
		entidad := r.FormValue("entidad")
		if entidad == "" {
			empty = "El campo entidad no puede estar vacío"
			fmt.Fprintf(w, "<div class='form-group text-warning'>%s</div>", empty)
		} else if almacen == "" {
			empty = "El campo almacen no puede estar vacío"
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
					_, err1 := db.Exec("UPDATE almacenes SET almacen=?, entidad_id=? WHERE id = ?", almacen, entidad, edit_id)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						bad = "Fallo al modificar almacen"
						fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
					} else {
						fmt.Fprint(w, "OK")
					}
				} else {
					bad = "Solo un usuario ROOT puede editar un almacen"
					fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
				}
			}
		}
	}
	//MOSTRAR ALMACENES EN UNA TABLA
	if accion == "tabla_almacen" {
		username := r.FormValue("username")
		query, err := db.Query("SELECT id FROM usuarios WHERE user = ?", username)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			var id, creador_id int
			var tiempo int64
			var almacen, entidad string
			err = query.Scan(&creador_id)
			if err != nil {
				Error.Println(err)
			}
			query, err := db.Query("SELECT almacenes.id, almacenes.almacen, almacenes.timestamp, entidades.nombre FROM entidades INNER JOIN almacenes ON almacenes.entidad_id = entidades.id WHERE almacenes.creador_id = ?", creador_id)
			if err != nil {
				Warning.Println(err)
			}
			for query.Next() {
				err = query.Scan(&id, &almacen, &tiempo, &entidad)
				if err != nil {
					Error.Println(err)
				}
				creacion := time.Unix(tiempo, 0)
				fmt.Fprintf(w, "<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Pulsa para editar almacen'>%s</a></td><td>%s</td><td>%s</td></tr>",
					id, almacen, creacion, entidad)
			}
		}
	}
	//CARGA LOS DATOS DE UN ALMACEN EN UN FORMULARIO
	if accion == "load_almacen" {
		edit_id := r.FormValue("edit_id")
		var id, ent_id int
		var almacen string
		query, err := db.Query("SELECT id, almacen, entidad_id FROM almacenes WHERE id = ?", edit_id)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			err = query.Scan(&id, &almacen, &ent_id)
			if err != nil {
				Error.Println(err)
			}
			fmt.Fprintf(w, "id=%d&almacen=%s&entidad=%d", id, almacen, ent_id)
		}
	}
	//MOSTRAR UN SELECT DE ENTIDADES SEGUN SU CREADOR (almacenes.html)
	if accion == "almacen_entidad" {
		user := r.FormValue("username")
		query, err := db.Query("SELECT id FROM usuarios WHERE user = ?", user)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			var id int
			err = query.Scan(&id)
			if err != nil {
				Error.Println(err)
			}
			//Muestra un select de entidades por usuario
			var list string
			query, err := db.Query("SELECT id, nombre FROM entidades WHERE creador_id = ?", id)
			if err != nil {
				Error.Println(err)
			}
			list = "<div class='panel-heading'>Entidad</div><div class='panel-body'><select id='entidad' name='entidad'>"
			if query.Next() {
				var id_ent int
				var name string
				err = query.Scan(&id_ent, &name)
				if err != nil {
					Error.Println(err)
				}
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_ent, name)
				for query.Next() {
					err = query.Scan(&id_ent, &name)
					if err != nil {
						Error.Println(err)
					}
					list += fmt.Sprintf("<option value='%d'>%s</option>", id_ent, name)
				}
				list += "</select></div>"
				fmt.Fprint(w, list)
			} else {
				list += "<option value=''>No hay entidades</option></select></div>"
				fmt.Fprint(w, list)
			}
		}
	}
}
