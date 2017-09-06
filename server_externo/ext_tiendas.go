package main

import (
	"fmt"
	"net/http"
	"time"
)

//GESTION DE TIENDAS (tiendas.html)
func tiendas(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	accion := r.FormValue("accion")
	//DAR DE ALTAR UNA NUEVA TIENDA
	if accion == "tienda" {
		username := r.FormValue("username")
		tienda := r.FormValue("tienda")
		provincia := r.FormValue("provincia")
		address := r.FormValue("address")
		phone := r.FormValue("phone")
		extra := r.FormValue("extra")
		var output string
		if tienda == "" || address == "" || phone == "" {
			output = "<div class='form-group text-warning'>No pueden haber campos vacíos</div>"
		} else if provincia == "" {
			output = "<div class='form-group text-warning'>Debe haber almenos una provincia</div>"
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
				//Si es un usuario super-admin o un usuario que tiene creador super-admin, le permitimos crear tiendas
				if padre_id == 0 || padre_id == id_admin {
					timestamp := time.Now().Unix()
					db_mu.Lock()
					_, err1 := db.Exec("INSERT INTO tiendas (`tienda`, `creador_id`, `timestamp`, `provincia_id`, `address`, `phone`, `extra`) VALUES (?, ?, ?, ?, ?, ?, ?)", tienda, id, timestamp, provincia, address, phone, extra)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						output = "<div class='form-group text-danger'>Fallo al añadir tienda</div>"
					} else {
						output = "OK"
					}
				} else {
					output = "<div class='form-group text-danger'>Solo un usuario ROOT puede añadir una tienda</div>"
				}
			}
		}
		fmt.Fprint(w, output)
	}
	//MODIFICAR / EDITAR UNA TIENDA
	if accion == "edit_tienda" {
		edit_id := r.FormValue("edit_id")
		username := r.FormValue("username")
		tienda := r.FormValue("tienda")
		provincia := r.FormValue("provincia")
		address := r.FormValue("address")
		phone := r.FormValue("phone")
		extra := r.FormValue("extra")
	
		if provincia == "" || tienda == "" || address == "" || phone == "" {
			empty = "No puede haber campos vacíos"
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
					_, err1 := db.Exec("UPDATE tiendas SET tienda=?, provincia_id=?, address=?, phone=?, extra=? WHERE id = ?", tienda, provincia, address, phone, extra, edit_id)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						bad = "Fallo al modificar tienda"
						fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
					} else {
						fmt.Fprint(w, "OK")
					}
				} else {
					bad = "Solo un usuario ROOT puede editar una tienda"
					fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
				}
			}
		}
	}
	//MOSTRAR TIENDAS EN UNA TABLA
	if accion == "tabla_tienda" {
		username := r.FormValue("username")
		query, err := db.Query("SELECT id FROM usuarios WHERE user = ?", username)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			var id, creador_id int
			var tiempo int64
			var provincia, tienda, address, phone, extra string
			err = query.Scan(&creador_id)
			if err != nil {
				Error.Println(err)
			}
			query, err := db.Query("SELECT tiendas.id, tiendas.tienda, tiendas.timestamp, provincia.provincia, tiendas.address, tiendas.phone, tiendas.extra FROM tiendas INNER JOIN provincia ON tiendas.provincia_id = provincia.id WHERE tiendas.creador_id = ?", creador_id)
			if err != nil {
				Warning.Println(err)
			}
			for query.Next() {
				err = query.Scan(&id, &tienda, &tiempo, &provincia, &address, &phone, &extra)
				if err != nil {
					Error.Println(err)
				}
				creacion := time.Unix(tiempo, 0)
				fmt.Fprintf(w, "<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Pulsa para editar tienda'>%s</a></td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", 
							id, tienda, creacion, provincia, address, phone, extra)
			}
		}
	}
	//CARGA LOS DATOS DE UNA TIENDA EN UN FORMULARIO
	if accion == "load_tienda" { 
		edit_id := r.FormValue("edit_id")
		var id, prov_id int
		var tienda, address, phone, extra string
		query, err := db.Query("SELECT id, tienda, provincia_id, address, phone, extra FROM tiendas WHERE id = ?", edit_id)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			err = query.Scan(&id, &tienda, &prov_id, &address, &phone, &extra)
			if err != nil {
				Error.Println(err)
			}
			fmt.Fprintf(w, "id=%d&tienda=%s&provincia=%d&address=%s&phone=%s&extra=%s", id, tienda, prov_id, address, phone, extra)
		}
	}
	//MOSTRAR UN SELECT DE PROVINCIAS SEGUN SU CREADOR (tiendas.html)
	if accion == "tienda_provincia" {
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
			//Muestra un select de provincias por usuario
			var list string
			query, err := db.Query("SELECT id, provincia FROM provincia WHERE creador_id = ?", id)
			if err != nil {
				Error.Println(err)
			}
			list = "<div class='panel-heading'>Provincia</div><div class='panel-body'><select id='provincia' name='provincia'>"
			if query.Next() {
				var id_prov int
				var name string
				err = query.Scan(&id_prov, &name)
				if err != nil {
					Error.Println(err)
				}
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_prov, name)
				for query.Next() {
					err = query.Scan(&id_prov, &name)
					if err != nil {
						Error.Println(err)
					}
					list += fmt.Sprintf("<option value='%d'>%s</option>", id_prov, name)
				}
				list += "</select></div>"
				fmt.Fprint(w, list)
			} else {
				list += "<option value=''>No hay provincias</option></select></div>"
				fmt.Fprint(w, list)
			}
		}
	}
}