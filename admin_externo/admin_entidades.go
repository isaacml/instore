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
	_, ok := user[sid]
	if ok {
		accion := r.FormValue("accion")
		updateExpires(sid)
		//Envio de datos al server_ext: Dar de alta una nueva entidad
		if accion == "entidad" {
			correct_res := libs.DeleteSplitsChars(r.FormValue("entidad"))
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(settings["serverroot"]+"/entidades.cgi", "accion;entidad", "entidad;"+correct_res, "username;"+username))
			fmt.Fprint(w, respuesta)
		}
		//Envio de datos al server_ext: Modificar los datos de una entidad concreta
		if accion == "edit_entidad" {
			//Eliminamos puntos, dos puntos y puntos comas
			correct_res := libs.DeleteSplitsChars(r.FormValue("entidad"))
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(settings["serverroot"]+"/entidades.cgi", "accion;edit_entidad", "edit_id;"+r.FormValue("id"), "entidad;"+correct_res, "username;"+username))
			fmt.Fprint(w, respuesta)
		}
		//Envio de datos al server_ext: Mostrar en una tabla los datos de una entidad
		if accion == "get_entidad" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(settings["serverroot"]+"/entidades.cgi", "accion;tabla_entidad", "username;"+username))
			fmt.Fprint(w, respuesta)
		}
		//Envio de datos al server_ext: Cargar los datos de una entidad concreta
		if accion == "load_entidad" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(settings["serverroot"]+"/entidades.cgi", "accion;load_entidad", "edit_id;"+r.FormValue("load")))
			fmt.Fprint(w, respuesta)
		}
		//Modifica el estado de la entidad(ON/OFF)
		if accion == "edit_status" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(settings["serverroot"]+"/entidades.cgi", "accion;edit_status", "edit_id;"+r.FormValue("id"), "st_ent;"+r.FormValue("st_ent")))
			fmt.Fprint(w, respuesta)
		}
	}
}
