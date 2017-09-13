package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
)

//Gestion de provincias (provincias.html)
func provincias(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	accion := r.FormValue("accion")
	_, ok := user[sid]
	fmt.Println(sid, accion)
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		//Envio de datos al server_ext: Dar de alta una nueva provincia
		if accion == "provincia" {
			//Eliminamos puntos, dos puntos y puntos comas
			correct_res := libs.DeleteSplitsChars(r.FormValue("provincia"))
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/provincias.cgi", "accion;provincia", "provincia;"+correct_res, "username;"+username, "region;"+r.FormValue("region")))
			if respuesta == "OK" {
				good = "Provincia añadida correctamente"
				fmt.Fprintf(w, "<div class='form-group text-success'>%s</div>", good)
			} else {
				fmt.Fprint(w, respuesta)
			}
		}
		//Envio de datos al server_ext: Modificar los datos de una provincia concreta
		if accion == "edit_provincia" {
			//Eliminamos puntos, dos puntos y puntos comas
			correct_res := libs.DeleteSplitsChars(r.FormValue("provincia"))
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/provincias.cgi", "accion;edit_provincia", "edit_id;"+r.FormValue("id"), "provincia;"+correct_res, "username;"+username, "region;"+r.FormValue("region")))
			if respuesta == "OK" {
				good = "Provincia modificada correctamente"
				fmt.Fprintf(w, "<div class='form-group text-success'>%s</div>", good)
			} else {
				fmt.Fprint(w, respuesta)
			}
		}
		//Envio de datos al server_ext: Mostrar una tabla de provincias
		if accion == "get_provincia" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/provincias.cgi", "accion;tabla_provincia", "username;"+username))
			fmt.Fprint(w, respuesta)
		}
		//Envio de datos al server_ext: Cargar los datos de una provincia concreta
		if accion == "load_provincia" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/provincias.cgi", "accion;load_provincia", "edit_id;"+r.FormValue("load")))
			fmt.Fprint(w, respuesta)
		}
		//Envio de datos al server_ext: Generar un select de regiones para poder añadir una nueva provincia
		if accion == "provincia_region" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/provincias.cgi", "accion;provincia_region", "username;"+username))
			fmt.Fprint(w, respuesta)
		}
	}
}
