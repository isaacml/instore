package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
)

//GESTION DE ENTIDADES (entidades.html)
func orgs(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		accion := r.FormValue("accion")
		loadSettings(serverRoot)
		updateExpires(sid)
		if accion == "entidades" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/organizaciones.cgi", "accion;ent_org", "username;"+username))
			fmt.Fprint(w, respuesta)
		}
		if accion == "almacenes" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/organizaciones.cgi", "accion;alm_org", "username;"+username))
			fmt.Fprint(w, respuesta)
		}
		if accion == "paises" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/organizaciones.cgi", "accion;pais_org", "username;"+username))
			fmt.Fprint(w, respuesta)
		}
		if accion == "regiones" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/organizaciones.cgi", "accion;reg_org", "username;"+username))
			fmt.Fprint(w, respuesta)
		}
	}
}

func selected_org(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		var list string
		accion := r.FormValue("accion")
		loadSettings(serverRoot)
		updateExpires(sid)
		if accion == "ent" {
			list = "<div class='panel-heading'>Añadir Almacen</div><div class='panel-body'><form id='testform2' action=\"\"><div class='form-group'>"
			list += "<input name='ent_id' value='" + r.FormValue("entidad") + "' type='hidden'><input class='form-control' id='almacen' placeholder='Nombre de almacen' name='almacen' type='text' autofocus></div></form>"
			list += "<button id='enviar' class='btn btn-lg btn-success btn-block' name='enviar'>Enviar</button>"
			fmt.Println(r.Form)
		}
		if accion == "alm" {

		}
		fmt.Fprint(w, list)
	}
}
