package main

import (
	"fmt"
	"net/http"
	"time"
)

//GESTION DE PROVINCIAS (provincias.html)
func provincias(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	accion := r.FormValue("accion")
	//DAR DE ALTAR UNA NUEVA PROVINCIA
	if accion == "provincia" {
		username := r.FormValue("username")
		region := r.FormValue("region")
		provincia := r.FormValue("provincia")
		var output string
		if provincia == "" {
			output = "<div class='form-group text-warning'>El campo provincia no puede estar vacio</div>"
		} else if region == "" {
			output = "<div class='form-group text-warning'>Debe haber almenos una región</div>"
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
				//Si es un usuario super-admin o un usuario que tiene creador super-admin, le permitimos crear provincias
				if padre_id == 0 || padre_id == id_admin {
					timestamp := time.Now().Unix()
					db_mu.Lock()
					_, err1 := db.Exec("INSERT INTO provincia (`provincia`, `creador_id`, `timestamp`, `region_id`) VALUES (?, ?, ?, ?)", provincia, id, timestamp, region)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						output = "<div class='form-group text-danger'>Fallo al añadir provincia</div>"
					} else {
						output = "OK"
					}
				} else {
					output = "<div class='form-group text-danger'>Solo un usuario ROOT puede añadir una provincia</div>"
				}
			}
		}
		fmt.Fprint(w, output) 
	}
	//MODIFICAR / EDITAR UNA PROVINCIA
	if accion == "edit_provincia" {
		edit_id := r.FormValue("edit_id")
		username := r.FormValue("username")
		region := r.FormValue("region")
		provincia := r.FormValue("provincia")
	
		if provincia == "" {
			empty = "El campo provincia no puede estar vacío"
			fmt.Fprintf(w, "<div class='form-group text-warning'>%s</div>", empty)
		} else if region == "" {
			empty = "El campo región no puede estar vacío"
			fmt.Fprintf(w, "<div class='form-group text-warning'>%s</div>", empty)
		} else {
			query, err := db.Query("SELECT id, entidad_id FROM usuarios WHERE user = ?", username)
			if err != nil {
				Error.Println(err)
			}
			for query.Next() {
				var id, entidad_id int
				err = query.Scan(&id, &entidad_id)
				if err != nil {
					Error.Println(err)
				}
				if entidad_id == 0 {
					db_mu.Lock()
					_, err1 := db.Exec("UPDATE provincia SET provincia=?, region_id=? WHERE id = ?", provincia, region, edit_id)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						bad = "Fallo al modificar provincia"
						fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
					} else {
						fmt.Fprint(w, "OK")
					}
				} else {
					bad = "Solo un usuario ROOT puede editar una provincia"
					fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
				}
			}
		}
	}
	//MOSTRAR PROVINCIAS EN UNA TABLA
	if accion == "tabla_provincia" {
		username := r.FormValue("username")
		query, err := db.Query("SELECT id FROM usuarios WHERE user = ?", username)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			var id, creador_id int
			var tiempo int64
			var provincia, region string
			err = query.Scan(&creador_id)
			if err != nil {
				Error.Println(err)
			}
			query, err := db.Query("SELECT provincia.id, provincia.provincia, provincia.timestamp, region.region FROM provincia INNER JOIN region ON provincia.region_id = region.id WHERE provincia.creador_id = ?", creador_id)
			if err != nil {
				Warning.Println(err)
			}
			for query.Next() {
				err = query.Scan(&id, &provincia, &tiempo, &region)
				if err != nil {
					Error.Println(err)
				}
				creacion := time.Unix(tiempo, 0)
				fmt.Fprintf(w, "<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Pulsa para editar provincia'>%s</a></td><td>%s</td><td>%s</td></tr>", 
							id, provincia, creacion, region)
			}
		}
	}
	//CARGA LOS DATOS DE UNA PROVINCIA EN UN FORMULARIO
	if accion == "load_provincia" { 
		edit_id := r.FormValue("edit_id")
		var id, region_id int
		var provincia string
		query, err := db.Query("SELECT id, provincia, region_id FROM provincia WHERE id = ?", edit_id)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			err = query.Scan(&id, &provincia, &region_id)
			if err != nil {
				Error.Println(err)
			}
			fmt.Fprintf(w, "id=%d&provincia=%s&region=%d", id, provincia, region_id)
		}
	}
	//MOSTRAR UN SELECT DE REGIONES SEGUN SU CREADOR (provincias.html)
	if accion == "provincia_region" {
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
			//Muestra un select de regiones por usuario
			var list string
			query, err := db.Query("SELECT id, region FROM region WHERE creador_id = ?", id)
			if err != nil {
				Error.Println(err)
			}
			list = "<div class='panel-heading'>Región</div><div class='panel-body'><select id='region' name='region'>"
			if query.Next() {
				var id_region int
				var name string
				err = query.Scan(&id_region, &name)
				if err != nil {
					Error.Println(err)
				}
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_region, name)
				for query.Next() {
					err = query.Scan(&id_region, &name)
					if err != nil {
						Error.Println(err)
					}
					list += fmt.Sprintf("<option value='%d'>%s</option>", id_region, name)
				}
				list += "</select></div>"
				fmt.Fprint(w, list)
			} else {
				list += "<option value=''>No hay regiones</option></select></div>"
				fmt.Fprint(w, list)
			}
		}
	}
}

