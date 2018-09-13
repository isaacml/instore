package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
	"time"
)

//GESTION DE REGIONES (regiones.html)
func regiones(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	accion := r.FormValue("accion")
	//DAR DE ALTAR UNA NUEVA REGION
	if accion == "region" {
		var output, reg_name string
		var id, padre_id, id_admin, cont int
		username := r.FormValue("username")
		almacen := r.FormValue("almacen")
		region := r.FormValue("region")
		pais := r.FormValue("pais")
		if region == "" || almacen == "" || pais == "" {
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
			//Si es un usuario super-admin o un usuario que tiene creador super-admin, le permitimos crear regiones
			if padre_id == 0 || padre_id == id_admin {
				//Buscamos las regiones asociados a un determinado pais
				db_mu.Lock()
				regs, err := db.Query("SELECT region FROM region WHERE pais_id = ?", pais)
				db_mu.Unlock()
				if err != nil {
					Error.Println(err)
					return
				}
				for regs.Next() {
					err = regs.Scan(&reg_name)
					if err != nil {
						Error.Println(err)
						continue
					}
					//Se comprueba que no hay dos regiones con el mismo nombre
					if region == reg_name {
						cont++ //Si hay alguna region, el contador incrementa
					}
				}
				//Cont = 0, no hay ninguna region
				if cont == 0 {
					timestamp := time.Now().Unix()
					db_mu.Lock()
					_, err1 := db.Exec("INSERT INTO region (`region`, `creador_id`, `timestamp`, `pais_id`) VALUES (?, ?, ?, ?)", region, id, timestamp, pais)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						output = "<div class='form-group text-danger'>Fallo al añadir region</div>"
						return
					} else {
						output = "<div class='form-group text-success'>Región añadida correctamente</div>"
					}
				} else {
					output = "<div class='form-group text-danger'>El país ya tiene esa región asociada</div>"
				}
			} else {
				output = "<div class='form-group text-danger'>Solo un usuario ROOT puede añadir una region</div>"
			}
		}
		fmt.Fprintf(w, output)
	}
	//MODIFICAR / EDITAR UNA REGION
	if accion == "edit_region" {
		var output, reg_name string
		var id, padre_id, cont int
		edit_id := r.FormValue("edit_id")
		username := r.FormValue("username")
		region := r.FormValue("region")
		pais := r.FormValue("pais")
		if region == "" {
			output = "<div class='form-group text-warning'>La región no puede estar vacía</div>"
		} else {
			db_mu.Lock()
			err := db.QueryRow("SELECT id, padre_id FROM usuarios WHERE user = ?", username).Scan(&id, &padre_id)
			db_mu.Unlock()
			if err != nil {
				Error.Println(err)
				return
			}
			if padre_id == 0 || padre_id == 1 {
				//Buscamos las regiones asociadas a un determinado pais
				db_mu.Lock()
				query, err := db.Query("SELECT region FROM region WHERE pais_id = ? AND id != ?", pais, edit_id)
				db_mu.Unlock()
				if err != nil {
					Warning.Println(err)
					return
				}
				for query.Next() {
					err = query.Scan(&reg_name)
					if err != nil {
						Error.Println(err)
						continue
					}
					//Si hay alguno, el contador incrementa
					if reg_name == region {
						cont++
					}
				}
				//Cont = 0, no hay region asociada a pais
				if cont == 0 {
					db_mu.Lock()
					_, err1 := db.Exec("UPDATE region SET region=? WHERE id = ?", region, edit_id)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						output = "<div class='form-group text-danger'>Fallo al modificar región</div>"
						return
					} else {
						output = "<div class='form-group text-success'>Región modificada correctamente</div>"
					}
				} else {
					output = "<div class='form-group text-danger'>El país ya tiene esa región asociada</div>"
				}
			} else {
				output = "<div class='form-group text-danger'>Solo un usuario ROOT puede editar una región</div>"
			}
		}
		fmt.Fprintf(w, output)
	}
	//MOSTRAR REGIONES EN UNA TABLA
	if accion == "tabla_region" {
		var id, creador_id int
		var tiempo int64
		var pais, region, almacen string
		username := r.FormValue("username")
		db_mu.Lock()
		err := db.QueryRow("SELECT id FROM usuarios WHERE user = ?", username).Scan(&creador_id)
		db_mu.Unlock()
		if err != nil {
			Error.Println(err)
			return
		}
		db_mu.Lock()
		query, err := db.Query("SELECT region.id, region.region, region.timestamp, almacenes.almacen, pais.pais FROM region INNER JOIN pais ON region.pais_id = pais.id INNER JOIN almacenes ON almacenes.id = pais.almacen_id WHERE region.creador_id = ?", creador_id)
		db_mu.Unlock()
		if err != nil {
			Warning.Println(err)
			return
		}
		for query.Next() {
			err = query.Scan(&id, &region, &tiempo, &almacen, &pais)
			if err != nil {
				Error.Println(err)
				continue
			}
			//Se obtiene la fecha de creacion de un almacen
			f_creacion := libs.FechaCreacion(tiempo)
			cadena := "<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Pulsa para editar región'>%s</a>"
			cadena += "<a href='#' onclick='borrar(%d)' title='Borrar región' style='float:right'><span class='fa fa-trash-o'></a></td>"
			cadena += "<td>%s</td><td>%s</td><td>%s</td></tr>"
			fmt.Fprintf(w, cadena, id, region, id, f_creacion, almacen, pais)
		}
	}
	//CARGA LOS DATOS DE UNA REGION EN UN FORMULARIO
	if accion == "load_region" {
		edit_id := r.FormValue("edit_id")
		var id_reg, id_pais int
		var region string
		db_mu.Lock()
		err := db.QueryRow("SELECT id, region, pais_id FROM region WHERE id = ?", edit_id).Scan(&id_reg, &region, &id_pais)
		db_mu.Unlock()
		if err != nil {
			Error.Println(err)
			return
		}
		fmt.Fprintf(w, "id=%d&region=%s&id_pais=%d", id_reg, region, id_pais)
	}
	//BORRAR UNA REGION
	if accion == "del_reg" {
		var cont int
		db_mu.Lock()
		query, err := db.Query("SELECT * FROM provincia WHERE region_id = ?", r.FormValue("borrar"))
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
			_, err := db.Exec("DELETE FROM region WHERE id = ?", r.FormValue("borrar"))
			db_mu.Unlock()
			if err != nil {
				Error.Println(err)
				return
			}
			fmt.Fprint(w, "OK")
		} else {
			fmt.Fprint(w, "<div class='form-group text-danger'>Necesario borrar provincias de las que depende</div>")
		}
	}
	//MOSTRAR UN SELECT DE PAISES SEGUN SU ALMACEN
	if accion == "show_paises" {
		var list, name string
		var id_pais int
		//Muestra un select de paises asociado a un almacen
		db_mu.Lock()
		query, err := db.Query("SELECT id, pais FROM pais WHERE almacen_id = ?", r.FormValue("alm"))
		db_mu.Unlock()
		if err != nil {
			Error.Println(err)
			return
		}
		list = "<option value=''>[Seleccionar País]</option>"
		if query.Next() {
			err = query.Scan(&id_pais, &name)
			if err != nil {
				Error.Println(err)
				return
			}
			list += fmt.Sprintf("<option value='%d'>%s</option>", id_pais, name)
			for query.Next() {
				err = query.Scan(&id_pais, &name)
				if err != nil {
					Error.Println(err)
					continue
				}
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_pais, name)
			}
		} else {
			list += "<option value=''>No hay paises</option>"
		}
		fmt.Fprint(w, list)
	}
}
