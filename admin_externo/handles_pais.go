package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
)

//Envia el formulario para dar de alta un nuevo país
func pais(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		//Eliminamos puntos, dos puntos y puntos comas
		correct_res := libs.DeleteSplitsChars(r.FormValue("pais"))
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/pais.cgi", "pais;"+correct_res, "username;"+username, "almacen;"+r.FormValue("almacen")))
		if respuesta == "OK" {
			good = "País añadido correctamente"
			fmt.Fprintf(w, "<div class='form-group text-success'>%s</div>", good)
		} else {
			fmt.Fprint(w, respuesta)
		}
	}
}

//Envia el formulario para mostrar una tabla con los paises
func get_pais(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/get_pais.cgi", "username;"+username))
		fmt.Fprint(w, respuesta)
	}
}

//Envia el formulario para cargar un país concreto
func load_pais(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/load_pais.cgi", "edit_id;"+r.FormValue("load")))
		fmt.Fprint(w, respuesta)
	}
}

//Envia el formulario para modificar un país concreto
func edit_pais(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		//Eliminamos puntos, dos puntos y puntos comas
		correct_res := libs.DeleteSplitsChars(r.FormValue("pais"))
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/edit_pais.cgi", "edit_id;"+r.FormValue("id"), "almacen;"+r.FormValue("almacen"), "username;"+username, "pais;"+correct_res))
		if respuesta == "OK" {
			good = "País modificado correctamente"
			fmt.Fprintf(w, "<div class='form-group text-success'>%s</div>", good)
		} else {
			fmt.Fprint(w, respuesta)
		}
	}
}
