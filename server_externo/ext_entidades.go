package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
	"time"
)

//GESTION DE ENTIDADES (entidades.html)
func entidades(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	accion := r.FormValue("accion")
	//DAR DE ALTAR UNA NUEVA ENTIDAD
	if accion == "entidad" {
		var id_user, padre_id, id_admin, cont int
		var output, ent_name string
		username := r.FormValue("username")
		entidad := r.FormValue("entidad")
		if entidad == "" {
			output = "<div class='form-group text-warning'>El campo no puede estar vacío</div>"
		} else {
			//Obtenemos el id y el padre(creador) de un usuario concreto
			err := db.QueryRow("SELECT id, padre_id FROM usuarios WHERE user = ?", username).Scan(&id_user, &padre_id)
			if err != nil {
				Error.Println(err)
			}
			//Obtener el id del usuario super-admin
			err = db.QueryRow("SELECT id FROM usuarios WHERE padre_id = 0 AND entidad_id = 0").Scan(&id_admin)
			if err != nil {
				Error.Println(err)
			}
			//Si es un usuario super-admin o un usuario que tiene creador super-admin, le permitimos crear entidades
			if padre_id == 0 || padre_id == id_admin {
				//Select que muestra todas las entidades de un usuario concreto
				ents, err := db.Query("SELECT nombre FROM entidades WHERE creador_id = ?", id_user)
				if err != nil {
					Error.Println(err)
				}
				for ents.Next() {
					err = ents.Scan(&ent_name)
					if err != nil {
						Error.Println(err)
					}
					//Se comprueba que no hay dos entidades con el mismo nombre
					if ent_name == entidad {
						cont++ //Si hay alguna entidad, el contador incrementa
					}
				}
				//Cont = 0, no hay ninguna entidad
				if cont == 0 {
					timestamp := time.Now().Unix()
					db_mu.Lock()
					_, err1 := db.Exec("INSERT INTO entidades (`nombre`, `creador_id`, `timestamp`, `last_access`, `status`) VALUES (?, ?, ?, ?, ?)", entidad, id_user, timestamp, timestamp, 1)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						output = "<div class='form-group text-danger'>Fallo al añadir entidad</div>"
					} else {
						output = "<div class='form-group text-success'>Entidad añadida correctamente</div>"
					}
				} else {
					output = "<div class='form-group text-danger'>Ya existe una entidad con ese nombre</div>"
				}
			} else {
				output = "<div class='form-group text-danger'>Solo un usuario ROOT puede añadir una entidad</div>"
			}
		}
		fmt.Fprint(w, output)
	}
	//MODIFICAR / EDITAR UNA ENTIDAD
	if accion == "edit_entidad" {
		var id_user, padre_id, cont int
		var output, ent_name string
		edit_id := r.FormValue("edit_id")
		username := r.FormValue("username")
		entidad := r.FormValue("entidad")
		if entidad == "" {
			output = "<div class='form-group text-warning'>El campo no puede estar vacío</div>"
		} else {
			//Obtenemos el id y el padre(creador) del usuario conectado
			err := db.QueryRow("SELECT id, padre_id FROM usuarios WHERE user = ?", username).Scan(&id_user, &padre_id)
			if err != nil {
				Error.Println(err)
			}
			//Si es un usuario super-admin o  tiene creador super-admin, le permitimos modificar
			if padre_id == 0 || padre_id == 1 {
				//Select que muestra todas las entidades de un usuario concreto
				ents, err := db.Query("SELECT nombre FROM entidades WHERE creador_id = ?", id_user)
				if err != nil {
					Error.Println(err)
				}
				for ents.Next() {
					err = ents.Scan(&ent_name)
					if err != nil {
						Error.Println(err)
					}
					//Se comprueba que no hay dos entidades con el mismo nombre
					if ent_name == entidad {
						cont++ //Si hay alguna entidad, el contador incrementa
					}
				}
				//Cont = 0, no hay ninguna entidad
				if cont == 0 {
					db_mu.Lock()
					_, err1 := db.Exec("UPDATE entidades SET nombre=? WHERE id = ?", entidad, edit_id)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						output = "<div class='form-group text-danger'>Fallo al modificar entidad</div>"
					} else {
						output = "<div class='form-group text-success'>Entidad modificada correctamente</div>"
					}
				} else {
					output = "<div class='form-group text-danger'>Entidad existente, no se puede modificar</div>"
				}
			} else {
				output = "<div class='form-group text-danger'>Solo un usuario ROOT puede editar una entidad</div>"
			}
		}
		fmt.Fprint(w, output)
	}
	//MOSTRAR ENTIDADES EN UNA TABLA
	if accion == "tabla_entidad" {
		var id, creador_id, status int
		var tiempo int64
		var nombre, f_creacion, st string
		username := r.FormValue("username")
		//Obtenemos el id de usuario conectado
		err := db.QueryRow("SELECT id FROM usuarios WHERE user = ?", username).Scan(&creador_id)
		if err != nil {
			Error.Println(err)
		}
		query, err := db.Query("SELECT id, nombre, timestamp, status FROM entidades WHERE creador_id = ?", creador_id)
		if err != nil {
			Warning.Println(err)
		}
		for query.Next() {
			err = query.Scan(&id, &nombre, &tiempo, &status)
			if err != nil {
				Error.Println(err)
			}
			//Se obtiene la fecha de creacion de una entidad
			f_creacion = libs.FechaCreacion(tiempo)
			//Se evalua el estado
			if status == 1 {
				st = "<option value='1' selected>ON</option><option value='0'>OFF</option>"
			} else {
				st = "<option value='1'>ON</option><option value='0' selected>OFF</option>"
			}
			cadena := "<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Pulsa para editar entidad'>%s</a></td>"
			cadena += "<td>%s</td><td><select onchange='status(this.value, %d)'>%s</select></td></tr>"
			fmt.Fprintf(w, cadena, id, nombre, f_creacion, id, st)
		}
	}
	//CARGA LOS DATOS DE ENTIDAD EN UN FORMULARIO
	if accion == "load_entidad" {
		edit_id := r.FormValue("edit_id")
		var id, st int
		var nombre string
		query, err := db.Query("SELECT id, nombre, status FROM entidades WHERE id = ?", edit_id)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			err = query.Scan(&id, &nombre, &st)
			if err != nil {
				Error.Println(err)
			}
			fmt.Fprintf(w, "id=%d&entidad=%s&st_ent=%d", id, nombre, st)
		}
	}
	//MODIFICA EL ESTADO DE LA ENTIDAD
	if accion == "edit_status" {
		edit_id := r.FormValue("edit_id")
		st_ent := r.FormValue("st_ent")
		db_mu.Lock()
		_, err1 := db.Exec("UPDATE entidades SET status=? WHERE id = ?", st_ent, edit_id)
		db_mu.Unlock()
		if err1 != nil {
			Error.Println(err1)
		}
		fmt.Fprint(w, "<div class='form-group text-success'>Se ha cambiado el estado</div>")
	}
}
