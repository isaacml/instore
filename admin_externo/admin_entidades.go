package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
)

//GESTION DE ENTIDADES (entidades.html)
func entidades(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	accion := r.FormValue("accion")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		//Envio de datos al server_ext: Dar de alta una nueva entidad
		if accion == "entidad" {
			correct_res := libs.DeleteSplitsChars(r.FormValue("entidad"))
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/entidades.cgi", "accion;entidad", "entidad;"+correct_res, "username;"+username))
			if respuesta == "OK" {
				fmt.Fprint(w, "<div class='form-group text-success'>Entidad añadida correctamente</div>")
			} else {
				fmt.Fprint(w, respuesta)
			}
		}
		//Envio de datos al server_ext: Modificar los datos de una entidad concreta
		if accion == "edit_entidad" {
			//Eliminamos puntos, dos puntos y puntos comas
			correct_res := libs.DeleteSplitsChars(r.FormValue("entidad"))
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/entidades.cgi", "accion;edit_entidad", "edit_id;"+r.FormValue("id"), "entidad;"+correct_res, "username;"+username))
			if respuesta == "OK" {
				fmt.Fprint(w, "<div class='form-group text-success'>Entidad modificada correctamente</div>")
			} else {
				fmt.Fprint(w, respuesta)
			}
		}
		//Envio de datos al server_ext: Mostrar en una tabla los datos de una entidad
		if accion == "get_entidad" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/entidades.cgi", "accion;tabla_entidad", "username;"+username))
			fmt.Fprint(w, respuesta)
		}
		//Envio de datos al server_ext: Cargar los datos de una entidad concreta
		if accion == "load_entidad" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/entidades.cgi", "accion;load_entidad", "edit_id;"+r.FormValue("load")))
			fmt.Fprint(w, respuesta)
		}
		//Organizaciones
		if accion == "orgs" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/organizaciones.cgi", "username;"+username))
			fmt.Fprint(w, respuesta)
		}
		//Organizaciones
		if accion == "almacenes" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/organizaciones.cgi", "username;"+username))
			fmt.Fprint(w, respuesta)
		}
	}
}
