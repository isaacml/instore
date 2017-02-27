package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
)

//Envia el formulario para dar de alta un nuevo almacen
func almacen(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		//Eliminamos puntos, dos puntos y puntos comas
		correct_res := libs.DeleteSplitsChars(r.FormValue("almacen"))
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/almacen.cgi", "almacen;"+correct_res, "username;"+username, "entidad;"+r.FormValue("entidad")))
		if respuesta == "OK" {
			good = "Almacen añadido correctamente"
			fmt.Fprintf(w, "<div class='form-group text-success'>%s</div>", good)
		} else {
			fmt.Fprint(w, respuesta)
		}
	}
}

//Envia el formulario para escribir una tabla que muestra los almacenes
func get_almacen(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/get_almacen.cgi", "username;"+username))
		fmt.Fprint(w, respuesta)
	}
}

//Envia el formulario para cargar los datos de un almacen concreto
func load_almacen(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/load_almacen.cgi", "edit_id;"+r.FormValue("load")))
		fmt.Fprint(w, respuesta)
	}
}

//Envia el formulario para modificar los datos de un almacen concreto
func edit_almacen(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		//Eliminamos puntos, dos puntos y puntos comas
		correct_res := libs.DeleteSplitsChars(r.FormValue("almacen"))
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/edit_almacen.cgi", "edit_id;"+r.FormValue("id"), "almacen;"+correct_res, "username;"+username, "entidad;"+r.FormValue("entidad")))
		if respuesta == "OK" {
			good = "Almacen modificado correctamente"
			fmt.Fprintf(w, "<div class='form-group text-success'>%s</div>", good)
		} else {
			fmt.Fprint(w, respuesta)
		}
	}
}
