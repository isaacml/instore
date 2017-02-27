package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
)

//Envia el formulario para dar de alta una nueva provincia
func provincia(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		//Eliminamos puntos, dos puntos y puntos comas
		correct_res := libs.DeleteSplitsChars(r.FormValue("provincia"))
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/provincia.cgi", "provincia;"+correct_res, "username;"+username, "region;"+r.FormValue("region")))
		if respuesta == "OK" {
			good = "Provincia añadida correctamente"
			fmt.Fprintf(w, "<div class='form-group text-success'>%s</div>", good)
		} else {
			fmt.Fprint(w, respuesta)
		}
	}
}

//Envia el formulario para mostrar una tabla de provincias
func get_provincia(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/get_provincia.cgi", "username;"+username))
		fmt.Fprint(w, respuesta)
	}
}

//Envia el formulario para cargar los datos de una provincia concreta
func load_provincia(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/load_provincia.cgi", "edit_id;"+r.FormValue("load")))
		fmt.Fprint(w, respuesta)
	}
}

//Envia el formulario para modificar los datos de una provincia concreta
func edit_provincia(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		//Eliminamos puntos, dos puntos y puntos comas
		correct_res := libs.DeleteSplitsChars(r.FormValue("provincia"))
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/edit_provincia.cgi", "edit_id;"+r.FormValue("id"), "provincia;"+correct_res, "username;"+username, "region;"+r.FormValue("region")))
		if respuesta == "OK" {
			good = "Provincia modificada correctamente"
			fmt.Fprintf(w, "<div class='form-group text-success'>%s</div>", good)
		} else {
			fmt.Fprint(w, respuesta)
		}
	}
}
