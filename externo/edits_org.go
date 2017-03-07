package main

import (
	"fmt"
	"net/http"
)

//Función para modificar el nombre y contraseña del usuario propio
func edit_own_user(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	user := r.FormValue("username")
	old_user := r.FormValue("old_user")
	password := r.FormValue("password")
	repeat_password := r.FormValue("repeat-password")

	if user == "" || password == "" || repeat_password == "" {
		empty = "Los campos no pueden estar vacíos"
		fmt.Fprintf(w, "<div class='form-group text-warning'>%s</div>", empty)
	} else {
		//Solo si las contraseñas son iguales modificamos
		if password == repeat_password {
			query, err := db.Query("SELECT id, user FROM usuarios WHERE old_user = ?", old_user)
			if err != nil {
				Error.Println(err)
			}
			for query.Next() {
				var id int
				var user_bd string
				err = query.Scan(&id, &user_bd)
				if err != nil {
					Error.Println(err)
				}
				//Comprobamos que no hay dos Usuarios con el mismo nombre
				if user_bd == user {
					bad = "El usuario ya existe"
					fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
				} else {
					db_mu.Lock()
					_, err1 := db.Exec("UPDATE usuarios SET user=?, old_user=?, pass=? WHERE id = ?", user, user, password, id)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
					}
				}
			}
			fmt.Fprintf(w, "OK")
		} else {
			bad = "Las contraseñas no coinciden"
			fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
		}
	}
}

//Función que va a modificar a un usuario concreto por su identificador
func edit_user(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	edit_id := r.FormValue("edit_id")
	user := r.FormValue("user")
	name_user := r.FormValue("name_user")
	pass := r.FormValue("pass")
	padre := r.FormValue("padre")
	entidad := r.FormValue("entidad")
	admin_user := r.FormValue("admin_user")

	if user == "" || name_user == "" || pass == "" {
		empty = "Los campos no pueden estar vacíos"
		fmt.Fprintf(w, "<div class='form-group text-warning'>%s</div>", empty)
	} else {
		//Generar el bitmap de acciones en hexadecimal
		var bitmap int
		if r.FormValue("prog_pub") != "" {
			bitmap = PROG_PUB
		}
		if r.FormValue("prog_mus") != "" {
			bitmap = bitmap + PROG_MUS
		}
		if r.FormValue("prog_msg") != "" {
			bitmap = bitmap + PROG_MSG
		}
		if r.FormValue("add_mus") != "" {
			bitmap = bitmap + ADD_MUS
		}
		if r.FormValue("msg_auto") != "" {
			bitmap = bitmap + MSG_AUTO
		}
		if r.FormValue("msg_normal") != "" {
			bitmap = bitmap + MSG_NORMAL
		}

		//Aquí se guarda el valor del bitmap en hexadecimal
		bitmap_hex := fmt.Sprintf("%x", bitmap)

		query, err := db.Query("SELECT padre_id FROM usuarios WHERE user = ?", admin_user)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			var padre_id int
			err = query.Scan(&padre_id)
			if err != nil {
				Error.Println(err)
			}
			if padre_id == 0 {
				db_mu.Lock()
				_, err1 := db.Exec("UPDATE usuarios SET user=?, old_user=?, nombre_completo=?, pass=?, entidad_id=?, padre_id=?, bitmap_acciones=? WHERE id = ?",
					user, user, name_user, pass, entidad, padre, bitmap_hex, edit_id)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
					bad = "Fallo al modificar usuario"
					fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
				} else {
					good := "Usuario modificado correctamente"
					fmt.Fprintf(w, "<div class='form-group text-success'>%s</div>", good)
				}
			} else {
				db_mu.Lock()
				_, err1 := db.Exec("UPDATE usuarios SET user=?, old_user=?, nombre_completo=?, pass=?, bitmap_acciones=? WHERE id = ?", user, user, name_user, pass, bitmap_hex, edit_id)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
					bad = "Fallo al modificar usuario"
					fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
				} else {
					good := "Usuario modificado correctamente"
					fmt.Fprintf(w, "<div class='form-group text-success'>%s</div>", good)
				}
			}
		}
	}
}

//Función que va a modificar a una entidad concreta
func edit_entidad(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	edit_id := r.FormValue("edit_id")
	username := r.FormValue("username")
	entidad := r.FormValue("entidad")

	if entidad == "" {
		empty = "El campo no puede estar vacío"
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
			if padre_id == 0 {
				db_mu.Lock()
				_, err1 := db.Exec("UPDATE entidades SET nombre=? WHERE id = ?", entidad, edit_id)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
					bad = "Fallo al modificar entidad"
					fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
				} else {
					fmt.Fprint(w, "OK")
				}
			} else {
				bad = "Solo un usuario ROOT puede editar una entidad"
				fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
			}
		}
	}
}

//Función que va a modificar un almacen concreto
func edit_almacen(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	edit_id := r.FormValue("edit_id")
	username := r.FormValue("username")
	almacen := r.FormValue("almacen")
	entidad := r.FormValue("entidad")

	if entidad == "" {
		empty = "El campo entidad no puede estar vacío"
		fmt.Fprintf(w, "<div class='form-group text-warning'>%s</div>", empty)
	} else if almacen == "" {
		empty = "El campo almacen no puede estar vacío"
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
			if padre_id == 0 {
				db_mu.Lock()
				_, err1 := db.Exec("UPDATE almacenes SET almacen=?, entidad_id=? WHERE id = ?", almacen, entidad, edit_id)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
					bad = "Fallo al modificar almacen"
					fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
				} else {
					fmt.Fprint(w, "OK")
				}
			} else {
				bad = "Solo un usuario ROOT puede editar un almacen"
				fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
			}
		}
	}
}

//Función que va a modificar un pais concreto
func edit_pais(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	edit_id := r.FormValue("edit_id")
	username := r.FormValue("username")
	almacen := r.FormValue("almacen")
	pais := r.FormValue("pais")

	if almacen == "" {
		empty = "El campo almacen no puede estar vacío"
		fmt.Fprintf(w, "<div class='form-group text-warning'>%s</div>", empty)
	} else if pais == "" {
		empty = "El campo pais no puede estar vacío"
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
			if padre_id == 0 {
				db_mu.Lock()
				_, err1 := db.Exec("UPDATE pais SET pais=?, almacen_id=? WHERE id = ?", pais, almacen, edit_id)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
					bad = "Fallo al modificar país"
					fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
				} else {
					fmt.Fprint(w, "OK")
				}
			} else {
				bad = "Solo un usuario ROOT puede editar un país"
				fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
			}
		}
	}
}

//Función que va a modificar una región concreta
func edit_region(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
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
			if padre_id == 0 {
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

//Función que va a modificar una provincia concreta
func edit_provincia(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	edit_id := r.FormValue("edit_id")
	username := r.FormValue("username")
	region := r.FormValue("region")
	provincia := r.FormValue("provincia")

	if provincia == "" {
		empty = "El campo provincia no puede estar vacío"
		fmt.Fprintf(w, "<div class='form-group text-warning'>%s</div>", empty)
	} else if region == "" {
		empty = "El campo región no puede estar vacío"
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
			if padre_id == 0 {
				db_mu.Lock()
				_, err1 := db.Exec("UPDATE provincia SET provincia=?, region_id=? WHERE id = ?", provincia, region, edit_id)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
					bad = "Fallo al modificar provincia"
					fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
				} else {
					fmt.Fprint(w, "OK")
				}
			} else {
				bad = "Solo un usuario ROOT puede editar una provincia"
				fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
			}
		}
	}
}

//Función que va a modificar una tienda concreta
func edit_tienda(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
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
			if padre_id == 0 {
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
