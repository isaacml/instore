package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
)

////Envia el formulario para dar de alta una nueva entidad
func entidad(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		//Eliminamos puntos, dos puntos y puntos comas
		correct_res := libs.DeleteSplitsChars(r.FormValue("entidad"))
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/entidad.cgi", "entidad;"+correct_res, "username;"+username))
		if respuesta == "OK" {
			good = "Entidad añadida correctamente"
			fmt.Fprintf(w, "<div class='form-group text-success'>%s</div>", good)
		} else {
			fmt.Fprint(w, respuesta)
		}
	}
}

//Envia el formulario para mostrar en una tabla los datos de una entidad
func get_entidad(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/get_entidad.cgi", "username;"+username))
		fmt.Fprint(w, respuesta)
	}
}

//Envia el formulario para cargar los datos de una entidad concreta
func load_entidad(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/load_entidad.cgi", "edit_id;"+r.FormValue("load")))
		fmt.Fprint(w, respuesta)
	}
}

//Envia el formulario para modificar los datos de una entidad concreta
func edit_entidad(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		//Eliminamos puntos, dos puntos y puntos comas
		correct_res := libs.DeleteSplitsChars(r.FormValue("entidad"))
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/edit_entidad.cgi", "edit_id;"+r.FormValue("id"), "entidad;"+correct_res, "username;"+username))
		if respuesta == "OK" {
			good = "Entidad modificada correctamente"
			fmt.Fprintf(w, "<div class='form-group text-success'>%s</div>", good)
		} else {
			fmt.Fprint(w, respuesta)
		}
	}
}
