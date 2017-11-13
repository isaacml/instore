package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
)

//Gestion de usuarios: tanto el propio (edit_own_user.html) como el resto de usuarios (alta_users.html)
func usuarios(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		accion := r.FormValue("accion")
		loadSettings(serverRoot)
		updateExpires(sid)
		//Muestra el usuario en activo
		if accion == "user_admin" {
			fmt.Fprint(w, username)
		}
		//Envio de datos al server_ext: Cambiar el usuario y contraseña del usuario activo/propio
		if accion == "own_user" {
			var output string
			//Eliminamos puntos, dos puntos y puntos comas
			clearUser := libs.DeleteSplitsChars(r.FormValue("username"))
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/usuarios.cgi", "accion;own_user", "username;"+clearUser, "old_user;"+username, "password;"+r.FormValue("pass"), "repeat-password;"+r.FormValue("repeat-password")))
			if respuesta == "OK" {
				username = clearUser //Al cambiar el nombre de usuario, es necesario actualizar la variable global username
				output = "<div class='form-group text-success'>Datos modificados correctamente</div>"
			} else {
				output = respuesta
			}
			fmt.Fprint(w, output)
		}
		//Envio de datos al server_ext: Alta de nuevo usuario
		if accion == "new_user" {
			//Eliminamos puntos, dos puntos y puntos comas
			clearUser := libs.DeleteSplitsChars(r.FormValue("user"))
			name_user := libs.DeleteSplitsChars(r.FormValue("name_user"))
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/usuarios.cgi", "accion;alta_users", "user;"+clearUser, "name_user;"+name_user, "pass;"+r.FormValue("pass"), "padre;"+username, "input_entidad;"+r.FormValue("entidad"),
				"prog_pub;"+r.FormValue("prog_pub"), "prog_mus;"+r.FormValue("prog_mus"), "prog_msg;"+r.FormValue("prog_msg"), "add_mus;"+r.FormValue("add_mus"), "msg_auto;"+r.FormValue("msg_auto"), "msg_normal;"+r.FormValue("msg_normal")))
			fmt.Fprint(w, respuesta)
		}
		//Envio de datos al server_ext: Tabla de usuarios
		if accion == "get_users" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/usuarios.cgi", "accion;tabla_users", "username;"+username))
			fmt.Fprint(w, respuesta)
		}
		//Envio de datos al server_ext: Cargar datos de un usuario concreto
		if accion == "load_user" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/usuarios.cgi", "accion;load_user", "edit_id;"+r.FormValue("load")))
			respuesta2 := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/usuarios.cgi", "accion;bitmap_checked", "edit_id;"+r.FormValue("load")))
			fmt.Fprint(w, respuesta+":.:"+respuesta2)
		}
		//Envio de datos al server_ext: Modificar/Editar un usuario concreto
		if accion == "edit_user" {
			//Eliminamos puntos, dos puntos y puntos comas
			clearUser := libs.DeleteSplitsChars(r.FormValue("user"))
			name_user := libs.DeleteSplitsChars(r.FormValue("name_user"))
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/usuarios.cgi", "accion;edit_user", "edit_id;"+r.FormValue("id"), "user;"+clearUser, "name_user;"+name_user, "pass;"+r.FormValue("pass"), "entidad;"+r.FormValue("entidad"), "admin_user;"+username,
				"prog_pub_edit;"+r.FormValue("prog_pub_edit"), "prog_mus_edit;"+r.FormValue("prog_mus_edit"), "prog_msg_edit;"+r.FormValue("prog_msg_edit"), "add_mus_edit;"+r.FormValue("add_mus_edit"), "msg_auto_edit;"+r.FormValue("msg_auto_edit"), "msg_normal_edit;"+r.FormValue("msg_normal_edit")))
			fmt.Fprint(w, respuesta)
		}
		//Envio de datos al server_ext: Mostrar un select de las entidades que tiene un usuario
		if accion == "listEnt" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/usuarios.cgi", "accion;user_entidad", "username;"+username))
			fmt.Fprint(w, respuesta)
		}
	}
}
