package main

import (
	"fmt"
	"net/http"
	"time"
)

//GESTION DE REGIONES (regiones.html)
func regiones(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	accion := r.FormValue("accion")
	//DAR DE ALTAR UNA NUEVA REGION
	if accion == "region" {
		username := r.FormValue("username")
		region := r.FormValue("region")
		pais := r.FormValue("pais")
		var output string
		if region == "" {
			output = "<div class='form-group text-warning'>El campo region no puede estar vacio</div>"
		} else if pais == "" {
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
				//Si es un usuario super-admin o un usuario que tiene creador super-admin, le permitimos crear regiones
				if padre_id == 0 || padre_id == id_admin {
					timestamp := time.Now().Unix()
					db_mu.Lock()
					_, err1 := db.Exec("INSERT INTO region (`region`, `creador_id`, `timestamp`, `pais_id`) VALUES (?, ?, ?, ?)", region, id, timestamp, pais)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						output = "<div class='form-group text-danger'>Fallo al añadir region</div>"
					} else {
						output = "OK"
					}
				} else {
					output = "<div class='form-group text-danger'>Solo un usuario ROOT puede añadir una region</div>"
				}
			}
		}
		fmt.Fprintf(w, output)
	}
	//MODIFICAR / EDITAR UNA REGION
	if accion == "edit_region" {
		edit_id := r.FormValue("edit_id")
		username := r.FormValue("username")
		region := r.FormValue("region")
		pais := r.FormValue("pais")

		if region == "" {
			empty = "El campo región no puede estar vacío"
			fmt.Fprintf(w, "<div class='form-group text-warning'>%s</div>", empty)
		} else if pais == "" {
			empty = "El campo país no puede estar vacío"
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
					_, err1 := db.Exec("UPDATE region SET region=?, pais_id=? WHERE id = ?", region, pais, edit_id)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						bad = "Fallo al modificar región"
						fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
					} else {
						fmt.Fprint(w, "OK")
					}
				} else {
					bad = "Solo un usuario ROOT puede editar una región"
					fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
				}
			}
		}
	}
	//MOSTRAR REGIONES EN UNA TABLA
	if accion == "tabla_region" {
		username := r.FormValue("username")
		query, err := db.Query("SELECT id FROM usuarios WHERE user = ?", username)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			var id, creador_id int
			var tiempo int64
			var pais, region string
			err = query.Scan(&creador_id)
			if err != nil {
				Error.Println(err)
			}
			query, err := db.Query("SELECT region.id, region.region, region.timestamp, pais.pais FROM region INNER JOIN pais ON region.pais_id = pais.id WHERE region.creador_id = ?", creador_id)
			if err != nil {
				Warning.Println(err)
			}
			for query.Next() {
				err = query.Scan(&id, &region, &tiempo, &pais)
				if err != nil {
					Error.Println(err)
				}
				creacion := time.Unix(tiempo, 0)
				fmt.Fprintf(w, "<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Pulsa para editar región'>%s</a></td><td>%s</td><td>%s</td></tr>",
					id, region, creacion, pais)
			}
		}
	}
	//CARGA LOS DATOS DE UNA REGION EN UN FORMULARIO
	if accion == "load_region" {
		edit_id := r.FormValue("edit_id")
		var id, pais_id int
		var region string
		query, err := db.Query("SELECT id, region, pais_id FROM region WHERE id = ?", edit_id)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			err = query.Scan(&id, &region, &pais_id)
			if err != nil {
				Error.Println(err)
			}
			fmt.Fprintf(w, "id=%d&region=%s&pais=%d", id, region, pais_id)
		}
	}
	//MOSTRAR UN SELECT DE PAISES SEGUN SU CREADOR (regiones.html)
	if accion == "region_pais" {
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
			//Muestra un select de paises por usuario
			var list string
			query, err := db.Query("SELECT id, pais FROM pais WHERE creador_id = ?", id)
			if err != nil {
				Error.Println(err)
			}
			list = "<div class='panel-heading'>País</div><div class='panel-body'><select id='pais' name='pais'>"
			if query.Next() {
				var id_pais int
				var name string
				err = query.Scan(&id_pais, &name)
				if err != nil {
					Error.Println(err)
				}
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_pais, name)
				for query.Next() {
					err = query.Scan(&id_pais, &name)
					if err != nil {
						Error.Println(err)
					}
					list += fmt.Sprintf("<option value='%d'>%s</option>", id_pais, name)
				}
				list += "</select></div>"
				fmt.Fprint(w, list)
			} else {
				list += "<option value=''>No hay paises</option></select></div>"
				fmt.Fprint(w, list)
			}
		}
	}
}
