package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
	"time"
)

//GESTION DE PROVINCIAS (provincias.html)
func provincias(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	accion := r.FormValue("accion")
	//DAR DE ALTAR UNA NUEVA PROVINCIA
	if accion == "provincia" {
		var id, padre_id, id_admin, cont int
		var output, prov_name string
		username := r.FormValue("username")
		almacen := r.FormValue("almacen")
		pais := r.FormValue("pais")
		region := r.FormValue("region")
		provincia := r.FormValue("provincia")
		if almacen == "" || pais == "" || region == "" || provincia == "" {
			output = "<div class='form-group text-warning'>Los campos no pueden estar vacíos</div>"
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
			//Si es un usuario super-admin o un usuario que tiene creador super-admin, le permitimos crear provincias
			if padre_id == 0 || padre_id == id_admin {
				//Buscamos las provincias asociadas a una region
				prov, err := db.Query("SELECT provincia FROM provincia WHERE region_id = ?", region)
				if err != nil {
					Error.Println(err)
				}
				for prov.Next() {
					err = prov.Scan(&prov_name)
					if err != nil {
						Error.Println(err)
					}
					//Se comprueba que no hay dos provincias con el mismo nombre
					if provincia == prov_name {
						cont++ //Si hay alguna provincia, el contador incrementa
					}
				}
				//Cont = 0, no hay ninguna region
				if cont == 0 {
					timestamp := time.Now().Unix()
					db_mu.Lock()
					_, err1 := db.Exec("INSERT INTO provincia (`provincia`, `creador_id`, `timestamp`, `region_id`) VALUES (?, ?, ?, ?)", provincia, id, timestamp, region)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						output = "<div class='form-group text-danger'>Fallo al añadir provincia</div>"
					} else {
						output = "<div class='form-group text-success'>Provincia añadida correctamente</div>"
					}
				} else {
					output = "<div class='form-group text-danger'>La región ya tiene esa provincia asociada</div>"
				}

			} else {
				output = "<div class='form-group text-danger'>Solo un usuario ROOT puede añadir una provincia</div>"
			}
		}
		fmt.Fprint(w, output)
	}
	//MODIFICAR / EDITAR UNA PROVINCIA
	if accion == "edit_provincia" {
		var output, prov_name string
		var id, padre_id, cont int
		edit_id := r.FormValue("edit_id")
		username := r.FormValue("username")
		provincia := r.FormValue("provincia")
		region := r.FormValue("region")
		if provincia == "" {
			output = "<div class='form-group text-warning'>El campo provincia no puede estar vacío</div>"
		} else {
			err := db.QueryRow("SELECT id, padre_id FROM usuarios WHERE user = ?", username).Scan(&id, &padre_id)
			if err != nil {
				Error.Println(err)
			}
			if padre_id == 0 || padre_id == 1 {
				//Buscamos las provincias asociadas a una region
				query, err := db.Query("SELECT provincia FROM provincia WHERE region_id = ? AND id != ?", region, edit_id)
				if err != nil {
					Warning.Println(err)
				}
				for query.Next() {
					err = query.Scan(&prov_name)
					if err != nil {
						Error.Println(err)
					}
					//Si hay alguna, el contador incrementa
					if prov_name == provincia {
						cont++
					}
				}
				//Cont = 0, no hay provincia asociada a region
				if cont == 0 {
					db_mu.Lock()
					_, err1 := db.Exec("UPDATE provincia SET provincia=? WHERE id = ?", provincia, edit_id)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						output = "<div class='form-group text-danger'>Fallo al modificar provincia</div>"
					} else {
						output = "<div class='form-group text-success'>Provincia modificada correctamente</div>"
					}
				} else {
					output = "<div class='form-group text-danger'>La región ya tiene esa provincia asociada</div>"
				}
			} else {
				output = "<div class='form-group text-danger'>Solo un usuario ROOT puede editar una provincia</div>"
			}
		}
		fmt.Fprint(w, output)
	}
	//MOSTRAR PROVINCIAS EN UNA TABLA
	if accion == "tabla_provincia" {
		var id, creador_id int
		var tiempo int64
		var provincia, region, pais, almacen string
		username := r.FormValue("username")
		err := db.QueryRow("SELECT id FROM usuarios WHERE user = ?", username).Scan(&creador_id)
		if err != nil {
			Error.Println(err)
		}
		query, err := db.Query("SELECT provincia.id, provincia.provincia, provincia.timestamp, region.region, pais.pais, almacenes.almacen FROM provincia INNER JOIN region ON provincia.region_id = region.id INNER JOIN pais ON region.pais_id = pais.id INNER JOIN almacenes ON almacenes.id = pais.almacen_id WHERE provincia.creador_id = ?", creador_id)
		if err != nil {
			Warning.Println(err)
		}
		for query.Next() {
			err = query.Scan(&id, &provincia, &tiempo, &region, &pais, &almacen)
			if err != nil {
				Error.Println(err)
			}
			//Se obtiene la fecha de creacion de una provincia
			f_creacion := libs.FechaCreacion(tiempo)
			cadena := "<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Pulsa para editar provincia'>%s</a>"
			cadena += "<a href='#' onclick='borrar(%d)' title='Borrar provincia' style='float:right'><span class='fa fa-trash-o'></a></td>"
			cadena += "<td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>"
			fmt.Fprintf(w, cadena, id, provincia, id, f_creacion, region, pais, almacen)
		}
	}
	//CARGA LOS DATOS DE UNA PROVINCIA EN UN FORMULARIO
	if accion == "load_provincia" {
		var id, id_reg int
		var provincia string
		edit_id := r.FormValue("edit_id")
		query, err := db.Query("SELECT id, provincia, region_id FROM provincia WHERE id = ?", edit_id)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			err = query.Scan(&id, &provincia, &id_reg)
			if err != nil {
				Error.Println(err)
			}
			fmt.Fprintf(w, "id=%d&provincia=%s&id_reg=%d", id, provincia, id_reg)
		}
	}
	//BORRAR UNA PROVINCIA
	if accion == "del_prov" {
		var cont int
		query, err := db.Query("SELECT * FROM tiendas WHERE provincia_id = ?", r.FormValue("borrar"))
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			cont++
		}
		if cont == 0 {
			db_mu.Lock()
			_, err := db.Exec("DELETE FROM provincia WHERE id = ?", r.FormValue("borrar"))
			db_mu.Unlock()
			if err != nil {
				Error.Println(err)
			}
			fmt.Fprint(w, "OK")
		} else {
			fmt.Fprint(w, "<div class='form-group text-danger'>Necesario borrar tiendas de las que depende</div>")
		}
	}
	//MOSTRAR UN SELECT DE REGIONES SEGUN PAIS
	if accion == "show_region" {
		var list string
		//Muestra un select de regiones por provincia
		query, err := db.Query("SELECT id, region FROM region WHERE pais_id = ?", r.FormValue("pais"))
		if err != nil {
			Error.Println(err)
		}
		list = "<option value=''>[Seleccionar Región]</option>"
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
		} else {
			list += "<option value=''>No hay regiones</option></select></div>"
		}
		fmt.Fprint(w, list)
	}
}
