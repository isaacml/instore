package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
	"time"
)

//GESTION DE ALMACENES (almacenes.html)
func almacenes(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	accion := r.FormValue("accion")
	//DAR DE ALTAR UN NUEVO ALMACEN
	if accion == "almacen" {
		var output, alm_name string
		var id_user, padre_id, id_admin, cont int
		username := r.FormValue("username")
		almacen := r.FormValue("almacen")
		entidad := r.FormValue("entidad")
		if almacen == "" {
			output = "<div class='form-group text-warning'>El campo no puede estar vacio</div>"
		} else if entidad == "" {
			output = "<div class='form-group text-warning'>Debe haber almenos una entidad</div>"
		} else {
			err := db.QueryRow("SELECT id, padre_id FROM usuarios WHERE user = ?", username).Scan(&id_user, &padre_id)
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
				//Select que muestra todas los almacenes de un usuario concreto
				alms, err := db.Query("SELECT almacen FROM almacenes WHERE creador_id = ?", id_user)
				if err != nil {
					Error.Println(err)
				}
				for alms.Next() {
					err = alms.Scan(&alm_name)
					if err != nil {
						Error.Println(err)
					}
					//Se comprueba que no hay dos almacenes con el mismo nombre
					if alm_name == almacen {
						cont++ //Si hay algun almacen, el contador incrementa
					}
				}
				//Cont = 0, no hay ningun almacen
				if cont == 0 {
					timestamp := time.Now().Unix()
					db_mu.Lock()
					_, err1 := db.Exec("INSERT INTO almacenes (`almacen`, `creador_id`, `timestamp`, `entidad_id`) VALUES (?, ?, ?, ?)", almacen, id_user, timestamp, entidad)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						output = "<div class='form-group text-danger'>Fallo al añadir almacen</div>"
					} else {
						output = "<div class='form-group text-success'>Almacen añadido correctamente</div>"
					}
				} else {
					output = "<div class='form-group text-danger'>Ya existe un almacen con ese nombre</div>"
				}
			} else {
				output = "<div class='form-group text-danger'>Solo un usuario ROOT puede añadir un almacen</div>"
			}
		}
		fmt.Fprint(w, output)
	}
	//MODIFICAR / EDITAR UN ALMACEN
	if accion == "edit_almacen" {
		var output, alm_name string
		var id_user, padre_id, cont int
		edit_id := r.FormValue("edit_id")
		username := r.FormValue("username")
		almacen := r.FormValue("almacen")
		entidad := r.FormValue("entidad")
		if entidad == "" {
			output = "<div class='form-group text-warning'>El campo entidad no puede estar vacío</div>"
		} else if almacen == "" {
			output = "<div class='form-group text-warning'>El campo almacen no puede estar vacío</div>"
		} else {
			//Obtenemos el id y el padre(creador) del usuario conectado
			err := db.QueryRow("SELECT id, padre_id FROM usuarios WHERE user = ?", username).Scan(&id_user, &padre_id)
			if err != nil {
				Error.Println(err)
			}
			if padre_id == 0 || padre_id == 1 {
				//Select que muestra todos los almacenes de un usuario concreto
				alms, err := db.Query("SELECT almacen FROM almacenes WHERE creador_id = ? AND id != ?", id_user, edit_id)
				if err != nil {
					Error.Println(err)
				}
				for alms.Next() {
					err = alms.Scan(&alm_name)
					if err != nil {
						Error.Println(err)
					}
					//Se comprueba que no hay dos almacenes con el mismo nombre
					if alm_name == almacen {
						cont++ //Si hay alguna entidad, el contador incrementa
					}
				} //Cont = 0, no hay ningun almacen
				if cont == 0 {
					db_mu.Lock()
					_, err1 := db.Exec("UPDATE almacenes SET almacen=?, entidad_id=? WHERE id = ?", almacen, entidad, edit_id)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						output = "<div class='form-group text-danger'>Fallo al modificar almacen</div>"
					} else {
						output = "<div class='form-group text-success'>Almacen modificado correctamente</div>"
					}
				} else {
					output = "<div class='form-group text-danger'>El almacen no se puede modificar o ya existe</div>"
				}
			} else {
				output = "<div class='form-group text-danger'>Solo un usuario ROOT puede editar un almacen</div>"
			}
		}
		fmt.Fprint(w, output)
	}
	//MOSTRAR ALMACENES EN UNA TABLA
	if accion == "tabla_almacen" {
		var id, creador_id int
		var tiempo int64
		var almacen, entidad string
		username := r.FormValue("username")
		err := db.QueryRow("SELECT id FROM usuarios WHERE user = ?", username).Scan(&creador_id)
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
			//Se obtiene la fecha de creacion de un almacen
			f_creacion := libs.FechaCreacion(tiempo)
			cadena := "<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Pulsa para editar almacen'>%s</a>"
			cadena += "<a href='#' onclick='borrar(%d)' title='Borrar almacen' style='float:right'><span class='fa fa-trash-o'></a></td>"
			cadena += "<td>%s</td><td>%s</td></tr>"
			fmt.Fprintf(w, cadena, id, almacen, id, f_creacion, entidad)
		}
	}
	//CARGA LOS DATOS DE UN ALMACEN EN UN FORMULARIO
	if accion == "load_almacen" {
		var id, ent_id int
		var almacen string
		edit_id := r.FormValue("edit_id")
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
	//BORRAR UN ALMACEN
	if accion == "del_alm" {
		var cont int
		query, err := db.Query("SELECT * FROM pais WHERE almacen_id = ?", r.FormValue("borrar"))
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			cont++
		}
		if cont == 0 {
			db_mu.Lock()
			_, err := db.Exec("DELETE FROM almacenes WHERE id = ?", r.FormValue("borrar"))
			db_mu.Unlock()
			if err != nil {
				Error.Println(err)
			}
			fmt.Fprint(w, "OK")
		} else {
			fmt.Fprint(w, "<div class='form-group text-danger'>Necesario borrar paises de los que depende</div>")
		}
	}
	//MOSTRAR UN SELECT DE ENTIDADES SEGUN SU CREADOR (almacenes.html)
	if accion == "show_ent" {
		var id int
		var list string
		user := r.FormValue("username")
		err := db.QueryRow("SELECT id FROM usuarios WHERE user = ?", user).Scan(&id)
		if err != nil {
			Error.Println(err)
		}
		//Muestra un select de entidades por usuario
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
		} else {
			list += "<option value=''>No hay entidades</option></select></div>"
		}
		fmt.Fprint(w, list)
	}
}
