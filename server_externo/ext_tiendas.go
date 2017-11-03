package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
	"time"
)

//GESTION DE TIENDAS (tiendas.html)
func tiendas(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	accion := r.FormValue("accion")
	//DAR DE ALTAR UNA NUEVA TIENDA
	if accion == "tienda" {
		var output, shop_name string
		var id, padre_id, id_admin, cont int
		username := r.FormValue("username")
		tienda := r.FormValue("tienda")
		provincia := r.FormValue("provincia")
		address := r.FormValue("address")
		phone := r.FormValue("phone")
		extra := r.FormValue("extra")
		if tienda == "" || address == "" || phone == "" {
			output = "<div class='form-group text-warning'>No pueden haber campos vacíos</div>"
		} else if provincia == "" {
			output = "<div class='form-group text-warning'>Debe haber almenos una provincia</div>"
		} else {
			err := db.QueryRow("SELECT id, padre_id FROM usuarios WHERE user = ?", username).Scan(&id, &padre_id)
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
				//Buscamos las tiendas asociadas a una provincia
				shop, err := db.Query("SELECT tienda FROM tiendas WHERE provincia_id = ?", provincia)
				if err != nil {
					Error.Println(err)
				}
				for shop.Next() {
					err = shop.Scan(&shop_name)
					if err != nil {
						Error.Println(err)
					}
					//Se comprueba que no hay dos tiendas con el mismo nombre
					if tienda == shop_name {
						cont++ //Si hay alguna tienda, el contador incrementa
					}
				}
				//Cont = 0, no hay ninguna region
				if cont == 0 {
					timestamp := time.Now().Unix()
					db_mu.Lock()
					_, err1 := db.Exec("INSERT INTO tiendas (`tienda`, `creador_id`, `timestamp`, `provincia_id`, `address`, `phone`, `extra`) VALUES (?, ?, ?, ?, ?, ?, ?)", tienda, id, timestamp, provincia, address, phone, extra)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						output = "<div class='form-group text-danger'>Fallo al añadir tienda</div>"
					} else {
						output = "<div class='form-group text-success'>Tienda añadida correctamente</div>"
					}
				} else {
					output = "<div class='form-group text-danger'>Esa tienda ya existe</div>"
				}
			} else {
				output = "<div class='form-group text-danger'>Solo un usuario ROOT puede añadir una tienda</div>"
			}
		}
		fmt.Fprint(w, output)
	}
	//MODIFICAR / EDITAR UNA TIENDA
	if accion == "edit_tienda" {
		var id, padre_id int
		var output string
		edit_id := r.FormValue("edit_id")
		username := r.FormValue("username")
		tienda := r.FormValue("tienda")
		address := r.FormValue("address")
		phone := r.FormValue("phone")
		extra := r.FormValue("extra")
		if tienda == "" || address == "" || phone == "" {
			output = "<div class='form-group text-warning'>No puede haber campos vacíos</div>"
		} else {
			err := db.QueryRow("SELECT id, padre_id FROM usuarios WHERE user = ?", username).Scan(&id, &padre_id)
			if err != nil {
				Error.Println(err)
			}
			if padre_id == 0 || padre_id == 1 {
				db_mu.Lock()
				_, err1 := db.Exec("UPDATE tiendas SET tienda=?, address=?, phone=?, extra=? WHERE id = ?", tienda, address, phone, extra, edit_id)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
					output = "<div class='form-group text-danger'>Fallo al modificar tienda</div>"
				} else {
					output = "<div class='form-group text-success'>Tienda modificada correctamente</div>"
				}
			} else {
				output = "<div class='form-group text-danger'>Solo un usuario ROOT puede editar una tienda</div>"
			}
			fmt.Fprint(w, output)
		}
	}
	//MOSTRAR TIENDAS EN UNA TABLA
	if accion == "tabla_tienda" {
		var id, creador_id int
		var tiempo int64
		var provincia, tienda, address, phone, extra string
		username := r.FormValue("username")
		err := db.QueryRow("SELECT id FROM usuarios WHERE user = ?", username).Scan(&creador_id)
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
			//Se obtiene la fecha de creacion de un almacen
			f_creacion := libs.FechaCreacion(tiempo)
			fmt.Fprintf(w, "<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Pulsa para editar tienda'>%s</a></td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>",
				id, tienda, f_creacion, provincia, address, phone, extra)
		}
	}
	//CARGA LOS DATOS DE UNA TIENDA EN UN FORMULARIO
	if accion == "load_tienda" {
		var id int
		var tienda, address, phone, extra string
		edit_id := r.FormValue("edit_id")
		query, err := db.Query("SELECT id, tienda, address, phone, extra FROM tiendas WHERE id = ?", edit_id)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			err = query.Scan(&id, &tienda, &address, &phone, &extra)
			if err != nil {
				Error.Println(err)
			}
			fmt.Fprintf(w, "id=%d&tienda=%s&address=%s&phone=%s&extra=%s", id, tienda, address, phone, extra)
		}
	}
	//MOSTRAR UN SELECT DE PROVINCIAS SEGUN SU REGION
	if accion == "show_prov" {
		var list string
		//Muestra un select de provincias por usuario
		query, err := db.Query("SELECT id, provincia FROM provincia WHERE region_id = ?", r.FormValue("reg"))
		if err != nil {
			Error.Println(err)
		}
		list = "<option value=''>[Seleccionar Provincia]</option>"
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
		} else {
			list += "<option value=''>No hay provincias</option></select></div>"
		}
		fmt.Fprint(w, list)
	}
}
