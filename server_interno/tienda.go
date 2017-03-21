package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
)

// enviamos user y pass de autenticación al servidor interno
func login_tienda(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	loadSettings(serverRoot)
	//SE PASAN LAS VARIABLES DE AUTENTICACION
	respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverexterno"]+"/login.cgi", "user;"+r.FormValue("user"), "pass;"+r.FormValue("pass")))
	fmt.Fprint(w, respuesta)
}

// enviamos el username al servidor interno
func send_user(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	loadSettings(serverRoot)
	respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverexterno"]+"/config_shop.cgi", "action;entidad", "username;"+r.FormValue("user")))
	fmt.Fprint(w, respuesta)
}

// enviamos la entidad al servidor interno
func send_ent(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	loadSettings(serverRoot)
	respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverexterno"]+"/config_shop.cgi", "action;almacen", "entidad;"+r.FormValue("entidad")))
	fmt.Println(respuesta)
	fmt.Fprint(w, respuesta)
}
