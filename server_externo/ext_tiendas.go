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
		almacen := r.FormValue("almacen")
		pais := r.FormValue("pais")
		region := r.FormValue("region")
		provincia := r.FormValue("provincia")
		address := r.FormValue("address")
		phone := r.FormValue("phone")
		extra := r.FormValue("extra")
		if tienda == "" || address == "" || phone == "" || almacen == "" || pais == "" || region == "" || provincia == "" {
			output = "<div class='form-group text-warning'>Los campos no pueden estar vacíos</div>"
		} else {
			db_mu.Lock()
			err := db.QueryRow("SELECT id, padre_id FROM usuarios WHERE user = ?", username).Scan(&id, &padre_id)
			db_mu.Unlock()
			if err != nil {
				Error.Println(err)
				return
			}
			//Hacemos un select para obtener el id del usuario super-admin
			db_mu.Lock()
			err = db.QueryRow("SELECT id FROM usuarios WHERE padre_id = 0 AND entidad_id = 0").Scan(&id_admin)
			db_mu.Unlock()
			if err != nil {
				Error.Println(err)
				return
			}
			//Si es un usuario super-admin o un usuario que tiene creador super-admin, le permitimos crear tiendas
			if padre_id == 0 || padre_id == id_admin {
				//Buscamos las tiendas asociadas a una provincia
				db_mu.Lock()
				shop, err := db.Query("SELECT tienda FROM tiendas WHERE provincia_id = ?", provincia)
				db_mu.Unlock()
				if err != nil {
					Error.Println(err)
					return
				}
				for shop.Next() {
					err = shop.Scan(&shop_name)
					if err != nil {
						Error.Println(err)
						continue
					}
					//Se comprueba que no hay dos tiendas con el mismo nombre
					if tienda == shop_name {
						cont++ //Si hay alguna tienda, el contador incrementa
					}
				}
				//Cont = 0, no hay ninguna tienda
				if cont == 0 {
					timestamp := time.Now().Unix()
					db_mu.Lock()
					_, err1 := db.Exec("INSERT INTO tiendas (`tienda`, `creador_id`, `timestamp`, `provincia_id`, `address`, `phone`, `extra`, `last_connect`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", tienda, id, timestamp, provincia, address, phone, extra, timestamp)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						output = "<div class='form-group text-danger'>Fallo al añadir tienda</div>"
						return
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
		var id, padre_id, cont int
		var output, shop_name string
		edit_id := r.FormValue("edit_id")
		username := r.FormValue("username")
		tienda := r.FormValue("tienda")
		address := r.FormValue("address")
		phone := r.FormValue("phone")
		extra := r.FormValue("extra")
		provincia := r.FormValue("provincia")
		if tienda == "" || address == "" || phone == "" {
			output = "<div class='form-group text-warning'>No puede haber campos vacíos</div>"
		} else {
			db_mu.Lock()
			err := db.QueryRow("SELECT id, padre_id FROM usuarios WHERE user = ?", username).Scan(&id, &padre_id)
			db_mu.Unlock()
			if err != nil {
				Error.Println(err)
				return
			}
			if padre_id == 0 || padre_id == 1 {
				//Buscamos las tiendas asociadas a una provincia
				db_mu.Lock()
				query, err := db.Query("SELECT tienda FROM tiendas WHERE provincia_id = ? AND id != ?", provincia, edit_id)
				db_mu.Unlock()
				if err != nil {
					Warning.Println(err)
					return
				}
				for query.Next() {
					err = query.Scan(&shop_name)
					if err != nil {
						Error.Println(err)
						continue
					}
					//Si hay alguno, el contador incrementa
					if shop_name == tienda {
						cont++
					}
				}
				//Cont = 0, no hay tienda asociada a provincia
				if cont == 0 {
					db_mu.Lock()
					_, err1 := db.Exec("UPDATE tiendas SET tienda=?, address=?, phone=?, extra=? WHERE id = ?", tienda, address, phone, extra, edit_id)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						output = "<div class='form-group text-danger'>Fallo al añadir la tienda</div>"
						return
					} else {
						output = "<div class='form-group text-success'>Tienda añadida correctamente</div>"
					}
				} else {
					output = "<div class='form-group text-danger'>Esa tienda ya existe</div>"
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
		var provincia, region, pais, almacen, tienda, address, phone, extra string
		username := r.FormValue("username")
		db_mu.Lock()
		err := db.QueryRow("SELECT id FROM usuarios WHERE user = ?", username).Scan(&creador_id)
		db_mu.Unlock()
		if err != nil {
			Error.Println(err)
			return
		}
		db_mu.Lock()
		query, err := db.Query("SELECT tiendas.id, tiendas.tienda, tiendas.timestamp, provincia.provincia, region.region, pais.pais, almacenes.almacen, tiendas.address, tiendas.phone, tiendas.extra FROM tiendas INNER JOIN provincia ON tiendas.provincia_id = provincia.id INNER JOIN region ON region.id = provincia.region_id INNER JOIN pais ON region.pais_id = pais.id INNER JOIN almacenes ON almacenes.id = pais.almacen_id WHERE tiendas.creador_id = ?", creador_id)
		db_mu.Unlock()
		if err != nil {
			Warning.Println(err)
			return
		}
		for query.Next() {
			err = query.Scan(&id, &tienda, &tiempo, &provincia, &region, &pais, &almacen, &address, &phone, &extra)
			if err != nil {
				Error.Println(err)
				continue
			}
			//Se obtiene la fecha de creacion de un almacen
			f_creacion := libs.FechaCreacion(tiempo)
			cadena := "<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Pulsa para editar tienda'>%s</a><a href='#' onclick='borrar(%d)' title='Borrar tienda' style='float:right'>"
			cadena += "<span class='fa fa-trash-o'></a></td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>"
			fmt.Fprintf(w, cadena, id, tienda, id, f_creacion, provincia, region, pais, almacen, address, phone, extra)
		}
	}
	//CARGA LOS DATOS DE UNA TIENDA EN UN FORMULARIO
	if accion == "load_tienda" {
		var id, id_prov int
		var tienda, address, phone, extra string
		edit_id := r.FormValue("edit_id")
		db_mu.Lock()
		query, err := db.Query("SELECT id, tienda, address, phone, extra, provincia_id FROM tiendas WHERE id = ?", edit_id)
		db_mu.Unlock()
		if err != nil {
			Error.Println(err)
			return
		}
		for query.Next() {
			err = query.Scan(&id, &tienda, &address, &phone, &extra, &id_prov)
			if err != nil {
				Error.Println(err)
				continue
			}
			fmt.Fprintf(w, "id=%d&tienda=%s&address=%s&phone=%s&extra=%s&id_prov=%d", id, tienda, address, phone, extra, id_prov)
		}
	}
	//BORRAR UNA TIENDA
	if accion == "del_tienda" {
		db_mu.Lock()
		_, err := db.Exec("DELETE FROM tiendas WHERE id = ?", r.FormValue("borrar"))
		db_mu.Unlock()
		if err != nil {
			Error.Println(err)
			return
		}
	}
	//MOSTRAR UN SELECT DE PROVINCIAS SEGUN SU REGION
	if accion == "show_prov" {
		var list string
		//Muestra un select de provincias por usuario
		db_mu.Lock()
		query, err := db.Query("SELECT id, provincia FROM provincia WHERE region_id = ?", r.FormValue("reg"))
		db_mu.Unlock()
		if err != nil {
			Error.Println(err)
			return
		}
		list = "<option value=''>[Seleccionar Provincia]</option>"
		if query.Next() {
			var id_prov int
			var name string
			err = query.Scan(&id_prov, &name)
			if err != nil {
				Error.Println(err)
				return
			}
			list += fmt.Sprintf("<option value='%d'>%s</option>", id_prov, name)
			for query.Next() {
				err = query.Scan(&id_prov, &name)
				if err != nil {
					Error.Println(err)
					continue
				}
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_prov, name)
			}
		} else {
			list += "<option value=''>No hay provincias</option></select></div>"
		}
		fmt.Fprint(w, list)
	}
}
