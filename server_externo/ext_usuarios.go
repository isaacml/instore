package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
)

//Gestion de usuarios: tanto el propio (edit_own_user.html) como el resto de usuarios (alta_users.html)
func usuarios(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	accion := r.FormValue("accion")
	//Modificar nombre y contraseña del usuario propio
	if accion == "own_user" {
		var output string
		user := r.FormValue("username")
		old_user := r.FormValue("old_user")
		password := r.FormValue("password")
		repeat_password := r.FormValue("repeat-password")
		if user == "" || password == "" || repeat_password == "" {
			output = "<div class='form-group text-warning'>Los campos no pueden estar vacíos</div>"
		} else {
			//Solo si las contraseñas son iguales modificamos
			if password == repeat_password {
				db_mu.Lock()
				query, err := db.Query("SELECT id, user FROM usuarios WHERE old_user = ?", old_user)
				db_mu.Unlock()
				if err != nil {
					Error.Println(err)
					return
				}
				for query.Next() {
					var id int
					var user_bd string
					err = query.Scan(&id, &user_bd)
					if err != nil {
						Error.Println(err)
						continue
					}
					//Comprobamos que no hay dos Usuarios con el mismo nombre
					if user_bd == user {
						output = "<div class='form-group text-danger'>El usuario ya existe</div>"
					} else {
						db_mu.Lock()
						_, err1 := db.Exec("UPDATE usuarios SET user=?, old_user=?, pass=? WHERE id = ?", user, user, password, id)
						db_mu.Unlock()
						if err1 != nil {
							Error.Println(err1)
							continue
						}
						output = "OK"
					}
				}
			} else {
				output = "<div class='form-group text-danger'>Las contraseñas no coinciden</div>"
			}
		}
		fmt.Fprint(w, output)
	}
	//DAR DE ALTA NUEVOS USUARIOS
	if accion == "alta_users" {
		var output string //FORMATO DE SALIDA: usuario;nom_completo;contraseña;status
		//VARIABLE DE FORMULARIO
		user := r.FormValue("user")
		name_user := r.FormValue("name_user")
		pass := r.FormValue("pass")
		father := r.FormValue("padre")
		input_entidad := r.FormValue("input_entidad")
		//Comprobamos que ninguno de los campos esté vacio
		if user == "" || name_user == "" || pass == "" || input_entidad == "" {
			output = fmt.Sprintf("%s;%s;;<div class='form-group text-warning'>Hay campos vacios</div>", user, name_user)
		} else {
			var existe_usuario int
			//Comprobamos si existe o NO el usuario en base de datos
			db_mu.Lock()
			err1 := db.QueryRow("SELECT count(*) FROM usuarios WHERE user = ?", user).Scan(&existe_usuario)
			db_mu.Unlock()
			if err1 != nil {
				Error.Println(err1)
				return
			}
			//Usuario no existe, continuamos...
			if existe_usuario == 0 {
				var id_admin, bitmap int
				//Tomamos el identificador del padre
				db_mu.Lock()
				err2 := db.QueryRow("SELECT id FROM usuarios WHERE user = ?", father).Scan(&id_admin)
				db_mu.Unlock()
				if err2 != nil {
					Error.Println(err2)
					return
				}
				//Generar el bitmap de acciones en hexadecimal
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
				if r.FormValue("msg_normal") != "" {
					bitmap = bitmap + MSG_NORMAL
				}
				bitmap_hex := fmt.Sprintf("%x", bitmap) //Se guarda el valor del bitmap en hexadecimal
				//Insertamos los datos en BD
				db_mu.Lock()
				_, err1 := db.Exec("INSERT INTO usuarios (`user`, `old_user`, `pass`, `nombre_completo`, `entidad_id`, `padre_id`, `bitmap_acciones`) VALUES (?,?,?,?,?,?,?)",
					user, user, pass, name_user, input_entidad, id_admin, bitmap_hex)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
					output = fmt.Sprintf(";;;<div class='form-group text-danger'>Fallo al añadir usuario</div>")
					return
				} else {
					output = fmt.Sprintf(";;;<div class='form-group text-success'>Usuario añadido correctamente</div>")
				}
			} else {
				//ERROR: el usuario ya existe
				output = fmt.Sprintf(";%s;;<div class='form-group text-danger'>El usuario ya existe, prueba con otro</div></div>", name_user)
			}
		}
		fmt.Fprint(w, output)
	}
	//MOSTRAR LOS USUARIOS EN UNA TABLA
	if accion == "tabla_users" {
		username := r.FormValue("username")
		var tabla string
		var id_user, dad_id int
		db_mu.Lock()
		err0 := db.QueryRow("SELECT id, padre_id FROM usuarios WHERE user = ?", username).Scan(&id_user, &dad_id)
		db_mu.Unlock()
		if err0 != nil {
			Error.Println(err0)
			return
		}
		tabla = "<table class='table table-striped table-bordered table-hover' id='dataTables-example'>"
		//padre = 0 : es un usuario SUPER-ADMIN, muestra todos los usuarios
		if dad_id == 0 {
			var id, padre_id int
			var user, all_name, pass, creador string
			db_mu.Lock()
			query, err := db.Query("SELECT id, user, nombre_completo, pass, padre_id FROM usuarios")
			db_mu.Unlock()
			if err != nil {
				Warning.Println(err)
				return
			}
			tabla += "<thead><tr><th>Usuario</th><th class='hidden-xs'>Nombre Completo</th><th>Contraseña</th><th>Creador</th></tr></thead><tbody>"
			for query.Next() {
				err = query.Scan(&id, &user, &all_name, &pass, &padre_id)
				if err != nil {
					Error.Println(err)
					continue
				}
				if padre_id != 0 {
					db_mu.Lock()
					err = db.QueryRow("SELECT user FROM usuarios WHERE id = ?", padre_id).Scan(&creador)
					db_mu.Unlock()
					if err != nil {
						Warning.Println(err)
						continue
					}
				}
				tabla += fmt.Sprintf("<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Pulsa para editar el usuario'>%s</a></td class='hidden-xs'><td>%s</td><td>%s</td><td>%s</td></tr>",
					id, user, all_name, pass, creador)
			}
		} else if dad_id == 1 { //padre = 1, su creador es el super-admin, muestra todos los usuarios que ha creado el y sus hijos
			var id, padre_id int
			var user, all_name, pass, creador string
			db_mu.Lock()
			query, err := db.Query("SELECT id, user, nombre_completo, pass, padre_id FROM usuarios WHERE entidad_id IN (SELECT id FROM entidades WHERE creador_id = ?)", id_user)
			db_mu.Unlock()
			if err != nil {
				Warning.Println(err)
				return
			}
			tabla += "<thead><tr><th>Usuario</th><th class='hidden-xs'>Nombre Completo</th><th>Contraseña</th><th>Creador</th></tr></thead><tbody>"
			for query.Next() {
				err = query.Scan(&id, &user, &all_name, &pass, &padre_id)
				if err != nil {
					Error.Println(err)
					continue
				}
				db_mu.Lock()
				err = db.QueryRow("SELECT user FROM usuarios WHERE id = ?", padre_id).Scan(&creador)
				db_mu.Unlock()
				if err != nil {
					Warning.Println(err)
					continue
				}
				tabla += fmt.Sprintf("<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Pulsa para editar el usuario'>%s</a></td><td class='hidden-xs'>%s</td><td>%s</td><td>%s</td></tr>",
					id, user, all_name, pass, creador)
			}
		} else { //Usuario Normal: Solo puede ver los usuarios que él ha creado
			var id int
			var user, all_name, pass string
			db_mu.Lock()
			query, err := db.Query("SELECT id, user, nombre_completo, pass FROM usuarios WHERE padre_id = ?", id_user)
			db_mu.Unlock()
			if err != nil {
				Warning.Println(err)
				return
			}
			tabla += "<thead><tr><th>Usuario</th><th class='hidden-xs'>Nombre Completo</th><th>Contraseña</th></tr></thead><tbody>"
			for query.Next() {
				err = query.Scan(&id, &user, &all_name, &pass)
				if err != nil {
					Error.Println(err)
					continue
				}
				tabla += fmt.Sprintf("<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Pulsa para editar el usuario'>%s</a></td><td class='hidden-xs'>%s</td><td>%s</td></tr>",
					id, user, all_name, pass)
			}
		}
		tabla += "</tbody></table>"
		fmt.Fprint(w, tabla)
	}
	//CARGA LOS DATOS DE USUARIO EN UN FORMULARIO
	if accion == "load_user" {
		edit_id := r.FormValue("edit_id")
		var id, ent_id int
		var user, all_name, pass string
		db_mu.Lock()
		query, err := db.Query("SELECT id, user, nombre_completo, pass, entidad_id FROM usuarios WHERE id = ?", edit_id)
		db_mu.Unlock()
		if err != nil {
			Error.Println(err)
			return
		}
		for query.Next() {
			err = query.Scan(&id, &user, &all_name, &pass, &ent_id)
			if err != nil {
				Error.Println(err)
				continue
			}
			fmt.Fprintf(w, "id=%d&user=%s&name_user=%s&pass=%s&entidad=%d", id, user, all_name, pass, ent_id)
		}
	}
	//CHEQUEO DEL BITMAP: marca o desmarca las acciones del bitmap
	if accion == "bitmap_checked" {
		var output, bitmap string
		edit_id := r.FormValue("edit_id")
		db_mu.Lock()
		err := db.QueryRow("SELECT bitmap_acciones FROM usuarios WHERE id = ?", edit_id).Scan(&bitmap)
		db_mu.Unlock()
		if err != nil {
			Error.Println(err)
			return
		}
		//Checkeado o No, segun el resultado al pasarle la mascara
		prog_pub := libs.BitmapParsing(bitmap, PROG_PUB) //res[0]
		if prog_pub != 0 {
			output += "<input type='checkbox' name='prog_pub_edit' value='1' checked/> Programar Publicidad<br>"
		} else {
			output += "<input type='checkbox' name='prog_pub_edit' value='1'/> Programar Publicidad<br>"
		}
		prog_mus := libs.BitmapParsing(bitmap, PROG_MUS) //res[1]
		if prog_mus != 0 {
			output += "<input type='checkbox' name='prog_mus_edit' value='2' checked/> Programar Música<br>"
		} else {
			output += "<input type='checkbox' name='prog_mus_edit' value='2'/> Programar Música<br>"
		}
		prog_msg := libs.BitmapParsing(bitmap, PROG_MSG) //res[2]
		if prog_msg != 0 {
			output += "<input type='checkbox' name='prog_msg_edit' value='4' checked/> Programar Mensajes Nuevos<br>"
		} else {
			output += "<input type='checkbox' name='prog_msg_edit' value='4'/> Programar Mensajes Nuevos<br>"
		}
		add_mus := libs.BitmapParsing(bitmap, ADD_MUS) //res[3]
		if add_mus != 0 {
			output += "<input type='checkbox' name='add_mus_edit' value='8' checked/> Añadir Música No Cifrada<br>"
		} else {
			output += "<input type='checkbox' name='add_mus_edit' value='8'/> Añadir Música No Cifrada<br>"
		}
		msg_normal := libs.BitmapParsing(bitmap, MSG_NORMAL) //res[4]
		if msg_normal != 0 {
			output += "<input type='checkbox' name='msg_normal_edit' value='32' checked/> Reproducir Mensajes Normales<br>"
		} else {
			output += "<input type='checkbox' name='msg_normal_edit' value='32'/> Reproducir Mensajes Normales<br>"
		}
		fmt.Fprint(w, output)
	}
	//MODIFICAR / EDITAR USUARIO (por su identificador)
	if accion == "edit_user" {
		var output string
		edit_id := r.FormValue("edit_id")
		user := r.FormValue("user")
		name_user := r.FormValue("name_user")
		pass := r.FormValue("pass")
		entidad := r.FormValue("entidad")
		admin_user := r.FormValue("admin_user")
		if user == "" || name_user == "" || pass == "" {
			output = "<div class='form-group text-danger'>Error al editar: hay campos vacíos</div>"
		} else {
			//Generar el bitmap de acciones en hexadecimal
			var bitmap int
			if r.FormValue("prog_pub_edit") != "" {
				bitmap = PROG_PUB
			}
			if r.FormValue("prog_mus_edit") != "" {
				bitmap = bitmap + PROG_MUS
			}
			if r.FormValue("prog_msg_edit") != "" {
				bitmap = bitmap + PROG_MSG
			}
			if r.FormValue("add_mus_edit") != "" {
				bitmap = bitmap + ADD_MUS
			}
			if r.FormValue("msg_normal_edit") != "" {
				bitmap = bitmap + MSG_NORMAL
			}
			//Aquí se guarda el valor del bitmap en hexadecimal
			bitmap_hex := fmt.Sprintf("%x", bitmap)
			db_mu.Lock()
			query, err := db.Query("SELECT padre_id FROM usuarios WHERE user = ?", admin_user)
			db_mu.Unlock()
			if err != nil {
				Error.Println(err)
				return
			}
			for query.Next() {
				var padre_id int
				err = query.Scan(&padre_id)
				if err != nil {
					Error.Println(err)
					continue
				}
				if padre_id == 0 || padre_id == 1 {
					db_mu.Lock()
					_, err1 := db.Exec("UPDATE usuarios SET user=?, old_user=?, nombre_completo=?, pass=?, entidad_id=?, bitmap_acciones=? WHERE id = ?",
						user, user, name_user, pass, entidad, bitmap_hex, edit_id)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						output = "<div class='form-group text-danger'>Fallo al modificar usuario</div>"
						continue
					} else {
						output = "<div class='form-group text-success'>Usuario modificado correctamente</div>"
					}
				} else {
					db_mu.Lock()
					_, err1 := db.Exec("UPDATE usuarios SET user=?, old_user=?, nombre_completo=?, pass=?, bitmap_acciones=? WHERE id = ?", user, user, name_user, pass, bitmap_hex, edit_id)
					db_mu.Unlock()
					if err1 != nil {
						Error.Println(err1)
						output = "<div class='form-group text-danger'>Fallo al modificar usuario</div>"
					} else {
						output = "<div class='form-group text-success'>Usuario modificado correctamente</div>"
					}
				}
			}
		}
		fmt.Fprint(w, output)
	}
	//MOSTRAR UN SELECT DE ENTIDADES PARA UN USUARIO (alta_users.html)
	if accion == "user_entidad" {
		user := r.FormValue("username") //Recogemos usuario autentificado en el panel de administrador
		var list string
		var id, entidad, padre int
		//tomamos el id_usuario, el id_entidad a la que pertenece y el id de su padre
		db_mu.Lock()
		err0 := db.QueryRow("SELECT id, entidad_id, padre_id FROM usuarios WHERE user = ?", user).Scan(&id, &entidad, &padre)
		db_mu.Unlock()
		if err0 != nil {
			Error.Println(err0)
			return
		}
		//padre = 0 : es un usuario SUPER-ADMIN, puede ver todas las entidades
		if padre == 0 {
			var name string
			var id_ent int
			db_mu.Lock()
			query, err := db.Query("SELECT id, nombre FROM entidades")
			db_mu.Unlock()
			if err != nil {
				Error.Println(err)
				return
			}
			if query.Next() {
				list = ";<div class='panel-heading'>Entidad</div><div class='panel-body'><select id='entidad' name='entidad'>"
				query.Scan(&id_ent, &name)
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_ent, name)
				for query.Next() {
					query.Scan(&id_ent, &name)
					if err != nil {
						Error.Println(err)
						continue
					}
					list += fmt.Sprintf("<option value='%d'>%s</option>", id_ent, name)
				}
				list += "</select></div>"
			} else {
				list = ";<div class='panel-heading'>Entidad</div><div class='panel-body'><select id='entidad' name='entidad'><option value='' selected>No hay entidades</option></select></div>"
			}
		} else if padre == 1 { //padre = 1, su creador es el super-admin, puede ver todas las entidades que ha creado.
			var name string
			var id_ent int
			db_mu.Lock()
			query, err := db.Query("SELECT id, nombre FROM entidades WHERE creador_id=?", id)
			db_mu.Unlock()
			if err != nil {
				Error.Println(err)
				return
			}
			if query.Next() {
				list = ";<div class='panel-heading'>Entidad</div><div class='panel-body'><select id='entidad' name='entidad'>"
				query.Scan(&id_ent, &name)
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_ent, name)
				for query.Next() {
					query.Scan(&id_ent, &name)
					if err != nil {
						Error.Println(err)
						continue
					}
					list += fmt.Sprintf("<option value='%d'>%s</option>", id_ent, name)
				}
				list += "</select></div>"
			} else {
				list = ";<div class='panel-heading'>Entidad</div><div class='panel-body'><select id='entidad' name='entidad'><option value='' selected>No hay entidades</option></select></div>"
			}
		} else { //Es un usuario normal: No puede ver ninguna entidad, los usuarios que añade, se añaden a su propia entidad
			input := fmt.Sprintf("<input id='entidad' name='entidad' value='%d'>", entidad)
			list = "DisableEnt;" + input
		}
		fmt.Fprint(w, list)
	}
}
