package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
)

//Enviar los datos de formulario para cambiar el usuario y contraseña del user activo
func edit_own_user(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		//Eliminamos puntos, dos puntos y puntos comas
		clearUser := libs.DeleteSplitsChars(r.FormValue("username"))
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/edit_own_user.cgi", "username;"+clearUser, "old_user;"+username, "password;"+r.FormValue("pass"), "repeat-password;"+r.FormValue("repeat-password")))
		if respuesta == "OK" {
			good = "Datos modificados correctamente"
			fmt.Fprintf(w, "<div class='form-group text-success'>%s</div>", good)
			username = clearUser //Al cambiar el nombre de usuario, es necesario actualizar la variable global username
		} else {
			fmt.Fprint(w, respuesta)
		}
	}
}

//Enviar los datos de formulario para dar de alta un usuario nuevo
func alta_users(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		//Eliminamos puntos, dos puntos y puntos comas
		clearUser := libs.DeleteSplitsChars(r.FormValue("user"))
		name_user := libs.DeleteSplitsChars(r.FormValue("name_user"))
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/alta_users.cgi", "user;"+clearUser, "name_user;"+name_user, "pass;"+r.FormValue("pass"), "padre;"+username, "input_padre;"+r.FormValue("padre"), "input_entidad;"+r.FormValue("entidad"),
			"prog_pub;"+r.FormValue("prog_pub"), "prog_mus;"+r.FormValue("prog_mus"), "prog_msg;"+r.FormValue("prog_msg"), "add_mus;"+r.FormValue("add_mus"), "msg_auto;"+r.FormValue("msg_auto"), "msg_normal;"+r.FormValue("msg_normal")))
		fmt.Fprint(w, respuesta)
	}
}

//Envia el formulario para mostrar los usuarios en una tabla
func get_users(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/get_users.cgi", "username;"+username))
		fmt.Fprint(w, respuesta)
	}
}

//Envia el formulario para cargar un usuario concreto
func load_user(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/load_user.cgi", "edit_id;"+r.FormValue("load")))
		respuesta2 := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/bitmap_checked.cgi", "edit_id;"+r.FormValue("load")))
		fmt.Fprint(w, respuesta+":.:"+respuesta2)
	}
}

//Envia el formulario para modificar un usuario concreto
func edit_user(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		//Eliminamos puntos, dos puntos y puntos comas
		clearUser := libs.DeleteSplitsChars(r.FormValue("user"))
		name_user := libs.DeleteSplitsChars(r.FormValue("name_user"))
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/edit_user.cgi", "edit_id;"+r.FormValue("id"), "user;"+clearUser, "name_user;"+name_user, "pass;"+r.FormValue("pass"), "padre;"+r.FormValue("padre"), "entidad;"+r.FormValue("entidad"), "admin_user;"+username,
			"prog_pub_edit;"+r.FormValue("prog_pub_edit"), "prog_mus_edit;"+r.FormValue("prog_mus_edit"), "prog_msg_edit;"+r.FormValue("prog_msg_edit"), "add_mus_edit;"+r.FormValue("add_mus_edit"), "msg_auto_edit;"+r.FormValue("msg_auto_edit"), "msg_normal_edit;"+r.FormValue("msg_normal_edit")))
		fmt.Fprint(w, respuesta)
	}
}

//Pasamos el nombre de usuario al servidor para obtener los bitmaps de acciones
func bitmaps(w http.ResponseWriter, r *http.Request) {
	respuesta := libs.GenerateFORM(serverext["serverroot"]+"/bitmap_actions.cgi", "user;"+username)
	fmt.Fprint(w, respuesta)
}
