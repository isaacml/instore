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
		region := r.FormValue("region")
		pais := r.FormValue("pais")
		if region == "" {
			output = "<div class='form-group text-warning'>El campo region no puede estar vacio</div>"
		} else if pais == "" {
			output = "<div class='form-group text-warning'>Debe haber almenos un almacen</div>"
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
			//Si es un usuario super-admin o un usuario que tiene creador super-admin, le permitimos crear regiones
			if padre_id == 0 || padre_id == id_admin {
				//Buscamos las regiones asociados a un determinado pais
				regs, err := db.Query("SELECT region FROM region WHERE pais_id = ?", pais)
				if err != nil {
					Error.Println(err)
				}
				for regs.Next() {
					err = regs.Scan(&reg_name)
					if err != nil {
						Error.Println(err)
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
			output = "<div class='form-group text-warning'>El campo región no puede estar vacío</div>"
		} else if pais == "" {
			output = "<div class='form-group text-warning'>El campo país no puede estar vacío</div>"
		} else {
			err := db.QueryRow("SELECT id, padre_id FROM usuarios WHERE user = ?", username).Scan(&id, &padre_id)
			if err != nil {
				Error.Println(err)
			}
			if padre_id == 0 || padre_id == 1 {
				//Buscamos las regiones asociados a un determinado pais
				regs, err := db.Query("SELECT region FROM region WHERE pais_id = ?", pais)
				if err != nil {
					Error.Println(err)
				}
				for regs.Next() {
					err = regs.Scan(&reg_name)
					if err != nil {
						Error.Println(err)
					}
					//Se comprueba que no hay dos regiones con el mismo nombre
					if region == reg_name {
						cont++ //Si hay alguna region, el contador incrementa
					}
				}
				//Cont = 0, no hay ninguna region
				if cont == 0 {
					db_mu.Lock()
					_, err1 := db.Exec("UPDATE region SET region=?, pais_id=? WHERE id = ?", region, pais, edit_id)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						output = "<div class='form-group text-danger'>Fallo al modificar región</div>"
					} else {
						output = "<div class='form-group text-success'>Región modificada correctamente</div>"
					}
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
		var pais, region string
		username := r.FormValue("username")
		err := db.QueryRow("SELECT id FROM usuarios WHERE user = ?", username).Scan(&creador_id)
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
			//Se obtiene la fecha de creacion de un almacen
			f_creacion := libs.FechaCreacion(tiempo)
			fmt.Fprintf(w, "<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Pulsa para editar región'>%s</a></td><td>%s</td><td>%s</td></tr>",
				id, region, f_creacion, pais)
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
		var id int
		var list string
		user := r.FormValue("username")
		err := db.QueryRow("SELECT id FROM usuarios WHERE user = ?", user).Scan(&id)
		if err != nil {
			Error.Println(err)
		}
		//Muestra un select de paises por usuario
		query, err := db.Query("SELECT id, pais FROM pais WHERE creador_id = ?", id)
		if err != nil {
			Error.Println(err)
		}
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
		} else {
			list += "<option value=''>No hay paises</option>"
		}
		fmt.Fprint(w, list)
	}
	//Busca las organizaciones padre de un determinado país
	if accion == "orgs_before" {
		var ent_name, alm_name, pais_name string
		err := db.QueryRow("SELECT entidades.nombre, almacenes.almacen, pais.pais FROM pais INNER JOIN almacenes ON pais.almacen_id = almacenes.id INNER JOIN entidades ON almacenes.entidad_id = entidades.id WHERE pais.id = ?", r.FormValue("pais_id")).Scan(&ent_name, &alm_name, &pais_name)
		if err != nil {
			Error.Println(err)
		}
		gen_orgs_before := fmt.Sprintf("<div class='form-group text-warning'>%s.%s.%s</div>", ent_name, alm_name, pais_name)
		fmt.Fprint(w, gen_orgs_before)
	}
}
