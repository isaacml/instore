package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
	"time"
)

//GESTION DE PAISES (paises.html)
func paises(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	accion := r.FormValue("accion")
	//DAR DE ALTAR UN NUEVO PAIS
	if accion == "pais" {
		var output, pais_name string
		var id, padre_id, id_admin, cont int
		username := r.FormValue("username")
		almacen := r.FormValue("almacen")
		pais := r.FormValue("pais")
		if pais == "" {
			output = "<div class='form-group text-warning'>El campo pais no puede estar vacio</div>"
		} else if almacen == "" {
			output = "<div class='form-group text-warning'>Debe haber almenos un almacen</div>"
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
			//Si es un usuario super-admin o un usuario que tiene creador super-admin, le permitimos crear paises
			if padre_id == 0 || padre_id == id_admin {
				//Buscamos los paises asociados a un determinado almacen
				db_mu.Lock()
				query, err := db.Query("SELECT pais FROM pais WHERE almacen_id = ?", almacen)
				db_mu.Unlock()
				if err != nil {
					Warning.Println(err)
					return
				}
				for query.Next() {
					err = query.Scan(&pais_name)
					if err != nil {
						Error.Println(err)
						continue
					}
					//Si hay alguno, el contador incrementa
					if pais_name == pais {
						cont++
					}
				}
				//Cont = 0, no hay ningun pais asociado
				if cont == 0 {
					timestamp := time.Now().Unix()
					db_mu.Lock()
					_, err1 := db.Exec("INSERT INTO pais (`pais`, `creador_id`, `timestamp`, `almacen_id`) VALUES (?, ?, ?, ?)", pais, id, timestamp, almacen)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						output = "<div class='form-group text-danger'>Fallo al añadir pais</div>"
						return
					} else {
						output = "<div class='form-group text-success'>País añadido correctamente</div>"
					}
				} else {
					output = "<div class='form-group text-danger'>El almacen ya tiene ese país asociado</div>"
				}
			} else {
				output = "<div class='form-group text-danger'>Solo un usuario ROOT puede añadir un pais</div>"
			}
		}
		fmt.Fprint(w, output)
	}
	//MODIFICAR / EDITAR UN PAIS
	if accion == "edit_pais" {
		var output, pais_name string
		var id, padre_id, cont int
		edit_id := r.FormValue("edit_id")
		username := r.FormValue("username")
		almacen := r.FormValue("almacen")
		pais := r.FormValue("pais")
		if almacen == "" {
			output = "<div class='form-group text-warning'>El campo almacen no puede estar vacío</div>"
		} else if pais == "" {
			output = "<div class='form-group text-warning'>El campo pais no puede estar vacío</div>"
		} else {
			db_mu.Lock()
			err := db.QueryRow("SELECT id, padre_id FROM usuarios WHERE user = ?", username).Scan(&id, &padre_id)
			db_mu.Unlock()
			if err != nil {
				Error.Println(err)
				return
			}
			if padre_id == 0 || padre_id == 1 {
				//Buscamos los paises asociados a un determinado almacen
				db_mu.Lock()
				query, err := db.Query("SELECT pais FROM pais WHERE almacen_id = ? AND id != ?", almacen, edit_id)
				db_mu.Unlock()
				if err != nil {
					Warning.Println(err)
					return
				}
				for query.Next() {
					err = query.Scan(&pais_name)
					if err != nil {
						Error.Println(err)
						continue
					}
					//Si hay alguno, el contador incrementa
					if pais_name == pais {
						cont++
					}
				}
				//Cont = 0, no hay ningun pais asociado
				if cont == 0 {
					db_mu.Lock()
					_, err1 := db.Exec("UPDATE pais SET pais=?, almacen_id=? WHERE id = ?", pais, almacen, edit_id)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						output = "<div class='form-group text-danger'>Fallo al modificar país</div>"
						return
					} else {
						output = "<div class='form-group text-success'>País modificado correctamente</div>"
					}
				} else {
					output = "<div class='form-group text-danger'>El almacen ya tiene ese país asociado</div>"
				}
			} else {
				output = "<div class='form-group text-danger'>Solo un usuario ROOT puede editar un país</div>"
			}
		}
		fmt.Fprint(w, output)
	}
	//MOSTRAR PAISES EN UNA TABLA
	if accion == "tabla_pais" {
		var id, creador_id int
		var tiempo int64
		var pais, almacen string
		username := r.FormValue("username")
		db_mu.Lock()
		err := db.QueryRow("SELECT id FROM usuarios WHERE user = ?", username).Scan(&creador_id)
		db_mu.Unlock()
		if err != nil {
			Error.Println(err)
			return
		}
		db_mu.Lock()
		query, err := db.Query("SELECT pais.id, pais.pais, pais.timestamp, almacenes.almacen FROM pais INNER JOIN almacenes ON pais.almacen_id = almacenes.id WHERE pais.creador_id = ?", creador_id)
		db_mu.Unlock()
		if err != nil {
			Warning.Println(err)
			return
		}
		for query.Next() {
			err = query.Scan(&id, &pais, &tiempo, &almacen)
			if err != nil {
				Error.Println(err)
				continue
			}
			//Se obtiene la fecha de creacion de un pais
			f_creacion := libs.FechaCreacion(tiempo)
			cadena := "<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Pulsa para editar país'>%s</a>"
			cadena += "<a href='#' onclick='borrar(%d)' title='Borrar país' style='float:right'><span class='fa fa-trash-o'></a></td>"
			cadena += "<td>%s</td><td>%s</td></tr>"
			fmt.Fprintf(w, cadena, id, pais, id, f_creacion, almacen)
		}
	}
	//CARGA LOS DATOS DE UN PAIS EN EL FORMULARIO
	if accion == "load_pais" {
		edit_id := r.FormValue("edit_id")
		var id, almacen_id int
		var pais string
		db_mu.Lock()
		query, err := db.Query("SELECT id, pais, almacen_id FROM pais WHERE id = ?", edit_id)
		db_mu.Unlock()
		if err != nil {
			Error.Println(err)
			return
		}
		for query.Next() {
			err = query.Scan(&id, &pais, &almacen_id)
			if err != nil {
				Error.Println(err)
				continue
			}
			fmt.Fprintf(w, "id=%d&pais=%s&almacen=%d", id, pais, almacen_id)
		}
	}
	//BORRAR UN PAIS
	if accion == "del_pais" {
		var cont int
		db_mu.Lock()
		query, err := db.Query("SELECT * FROM region WHERE pais_id = ?", r.FormValue("borrar"))
		db_mu.Unlock()
		if err != nil {
			Error.Println(err)
			return
		}
		for query.Next() {
			cont++
		}
		if cont == 0 {
			db_mu.Lock()
			_, err := db.Exec("DELETE FROM pais WHERE id = ?", r.FormValue("borrar"))
			db_mu.Unlock()
			if err != nil {
				Error.Println(err)
				return
			}
			fmt.Fprint(w, "OK")
		} else {
			fmt.Fprint(w, "<div class='form-group text-danger'>Necesario borrar regiones de las que depende</div>")
		}
	}
	//MOSTRAR UN SELECT DE ALMACENES SEGUN SU CREADOR (almacenes.html)
	if accion == "show_almacen" {
		var id int
		var list string
		user := r.FormValue("username")
		db_mu.Lock()
		err := db.QueryRow("SELECT id FROM usuarios WHERE user = ?", user).Scan(&id)
		db_mu.Unlock()
		if err != nil {
			Error.Println(err)
			return
		}
		//Muestra un select de almacenes por usuario
		db_mu.Lock()
		query, err := db.Query("SELECT id, almacen FROM almacenes WHERE creador_id = ?", id)
		db_mu.Unlock()
		if err != nil {
			Error.Println(err)
			return
		}
		list = "<option value=''>[Seleccionar Almacen]</option>"
		if query.Next() {
			var id_alm int
			var name string
			err = query.Scan(&id_alm, &name)
			if err != nil {
				Error.Println(err)
				return
			}
			list += fmt.Sprintf("<option value='%d'>%s</option>", id_alm, name)
			for query.Next() {
				err = query.Scan(&id_alm, &name)
				if err != nil {
					Error.Println(err)
					continue
				}
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_alm, name)
			}
		} else {
			list += "<option value=''>No hay almacenes</option></select></div>"
		}
		fmt.Fprint(w, list)
	}
}
