package main

import (
	"fmt"
	"net/http"
	"time"
)

//GESTION DE PAISES (paises.html)
func paises(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	accion := r.FormValue("accion")
	//DAR DE ALTAR UN NUEVO PAIS
	if accion == "pais" {
		username := r.FormValue("username")
		almacen := r.FormValue("almacen")
		pais := r.FormValue("pais")
		var output string
		if pais == "" {
			output = "<div class='form-group text-warning'>El campo pais no puede estar vacio</div>"
		} else if almacen == "" {
			output = "<div class='form-group text-warning'>Debe haber almenos un almacen</div>"
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
				//Si es un usuario super-admin o un usuario que tiene creador super-admin, le permitimos crear paises
				if padre_id == 0 || padre_id == id_admin {
					timestamp := time.Now().Unix()
					db_mu.Lock()
					_, err1 := db.Exec("INSERT INTO pais (`pais`, `creador_id`, `timestamp`, `almacen_id`) VALUES (?, ?, ?, ?)", pais, id, timestamp, almacen)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						output = "<div class='form-group text-danger'>Fallo al añadir pais</div>"
					} else {
						output = "OK"
					}
				} else {
					output = "<div class='form-group text-danger'>Solo un usuario ROOT puede añadir un pais</div>"
				}
			}
		}
		fmt.Fprint(w, output)
	}
	//MODIFICAR / EDITAR UN PAIS
	if accion == "edit_pais" {
		edit_id := r.FormValue("edit_id")
		username := r.FormValue("username")
		almacen := r.FormValue("almacen")
		pais := r.FormValue("pais")

		if almacen == "" {
			empty = "El campo almacen no puede estar vacío"
			fmt.Fprintf(w, "<div class='form-group text-warning'>%s</div>", empty)
		} else if pais == "" {
			empty = "El campo pais no puede estar vacío"
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
					_, err1 := db.Exec("UPDATE pais SET pais=?, almacen_id=? WHERE id = ?", pais, almacen, edit_id)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						bad = "Fallo al modificar país"
						fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
					} else {
						fmt.Fprint(w, "OK")
					}
				} else {
					bad = "Solo un usuario ROOT puede editar un país"
					fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
				}
			}
		}
	}
	//MOSTRAR PAISES EN UNA TABLA
	if accion == "tabla_pais" {
		username := r.FormValue("username")
		query, err := db.Query("SELECT id FROM usuarios WHERE user = ?", username)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			var id, creador_id int
			var tiempo int64
			var pais, almacen string
			err = query.Scan(&creador_id)
			if err != nil {
				Error.Println(err)
			}
			query, err := db.Query("SELECT pais.id, pais.pais, pais.timestamp, almacenes.almacen FROM pais INNER JOIN almacenes ON pais.almacen_id = almacenes.id WHERE pais.creador_id = ?", creador_id)
			if err != nil {
				Warning.Println(err)
			}
			for query.Next() {
				err = query.Scan(&id, &pais, &tiempo, &almacen)
				if err != nil {
					Error.Println(err)
				}
				creacion := time.Unix(tiempo, 0)
				fmt.Fprintf(w, "<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Pulsa para editar país'>%s</a></td><td>%s</td><td>%s</td></tr>",
					id, pais, creacion, almacen)
			}
		}
	}
	//CARGA LOS DATOS DE UN PAIS EN EL FORMULARIO
	if accion == "load_pais" {
		edit_id := r.FormValue("edit_id")
		var id, almacen_id int
		var pais string
		query, err := db.Query("SELECT id, pais, almacen_id FROM pais WHERE id = ?", edit_id)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			err = query.Scan(&id, &pais, &almacen_id)
			if err != nil {
				Error.Println(err)
			}
			fmt.Fprintf(w, "id=%d&pais=%s&almacen=%d", id, pais, almacen_id)
		}
	}
	//MOSTRAR UN SELECT DE ALMACENES SEGUN SU CREADOR (almacenes.html)
	if accion == "pais_almacen" {
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
			//Muestra un select de almacenes por usuario
			var list string
			query, err := db.Query("SELECT id, almacen FROM almacenes WHERE creador_id = ?", id)
			if err != nil {
				Error.Println(err)
			}
			list = "<div class='panel-heading'>Almacen</div><div class='panel-body'><select id='almacen' name='almacen'>"
			if query.Next() {
				var id_alm int
				var name string
				err = query.Scan(&id_alm, &name)
				if err != nil {
					Error.Println(err)
				}
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_alm, name)
				for query.Next() {
					err = query.Scan(&id_alm, &name)
					if err != nil {
						Error.Println(err)
					}
					list += fmt.Sprintf("<option value='%d'>%s</option>", id_alm, name)
				}
				list += "</select></div>"
				fmt.Fprint(w, list)
			} else {
				list += "<option value=''>No hay almacenes</option></select></div>"
				fmt.Fprint(w, list)
			}
		}
	}
}
