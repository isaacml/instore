package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
)

//Envia el formulario para dar de alta una nueva tienda
func tienda(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/tienda.cgi", "tienda;"+r.FormValue("tienda"), "username;"+username, "provincia;"+r.FormValue("provincia"), "address;"+r.FormValue("address"), "phone;"+r.FormValue("phone"), "extra;"+r.FormValue("extra")))
		if respuesta == "OK" {
			good = "Tienda añadida correctamente"
			fmt.Fprintf(w, "<div class='form-group text-success'>%s</div>", good)
		} else {
			fmt.Fprint(w, respuesta)
		}
	}
}

//Envia el formulario para mostrar en una tabla las tiendas
func get_tienda(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/get_tienda.cgi", "username;"+username))
		fmt.Fprint(w, respuesta)
	}
}

//Envia el formulario para cargar una tienda concreta
func load_tienda(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/load_tienda.cgi", "edit_id;"+r.FormValue("load")))
		fmt.Fprint(w, respuesta)
	}
}

//Envia el formulario para modificar una tienda concreta
func edit_tienda(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/edit_tienda.cgi", "edit_id;"+r.FormValue("id"), "tienda;"+r.FormValue("tienda"), "provincia;"+r.FormValue("provincia"), "username;"+username, "address;"+r.FormValue("address"), "phone;"+r.FormValue("phone"), "extra;"+r.FormValue("extra")))
		if respuesta == "OK" {
			good = "Tienda modificada correctamente"
			fmt.Fprintf(w, "<div class='form-group text-success'>%s</div>", good)
		} else {
			fmt.Fprint(w, respuesta)
		}
	}
}
