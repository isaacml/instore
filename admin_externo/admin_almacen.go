package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
)

//Gestion de almacenes (almacenes.html)
func almacenes(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		accion := r.FormValue("accion")
		loadSettings(serverRoot)
		updateExpires(sid)
		//Envio de datos al server_ext: Dar de alta un nuevo almacen
		if accion == "almacen" {
			//Eliminamos puntos, dos puntos y puntos comas
			correct_res := libs.DeleteSplitsChars(r.FormValue("almacen"))
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/almacenes.cgi", "accion;almacen", "almacen;"+correct_res, "username;"+username, "entidad;"+r.FormValue("entidad")))
			fmt.Fprint(w, respuesta)
		}
		//Envio de datos al server_ext: Modificar los datos de un almacen concreto
		if accion == "edit_almacen" {
			//Eliminamos puntos, dos puntos y puntos comas
			correct_res := libs.DeleteSplitsChars(r.FormValue("almacen"))
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/almacenes.cgi", "accion;edit_almacen", "edit_id;"+r.FormValue("id"), "almacen;"+correct_res, "username;"+username, "entidad;"+r.FormValue("entidad")))
			fmt.Fprint(w, respuesta)
		}
		//Envio de datos al server_ext:  Mostrar en una tabla los datos de los almacenes
		if accion == "get_almacen" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/almacenes.cgi", "accion;tabla_almacen", "username;"+username))
			fmt.Fprint(w, respuesta)
		}
		//Envio de datos al server_ext:  Cargar los datos de un almacen concreto
		if accion == "load_almacen" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/almacenes.cgi", "accion;load_almacen", "edit_id;"+r.FormValue("load")))
			fmt.Fprint(w, respuesta)
		}
		//Envio de datos al server_ext:  Generar un select de entidades para poder añadir un nuevo almacen
		if accion == "almacen_entidad" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/almacenes.cgi", "accion;almacen_entidad", "username;"+username))
			fmt.Fprint(w, respuesta)
		}
	}
}
