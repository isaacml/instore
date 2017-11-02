package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
)

//Gestion de tiendas (tiendas.html)
func tiendas(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		accion := r.FormValue("accion")
		loadSettings(serverRoot)
		updateExpires(sid)
		//Envio de datos al server_ext: Dar de alta una nueva tienda
		if accion == "tienda" {
			//Eliminamos puntos, dos puntos y puntos comas
			correct_res := libs.DeleteSplitsChars(r.FormValue("tienda"))
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/tiendas.cgi", "accion;tienda", "tienda;"+correct_res, "username;"+username, "provincia;"+r.FormValue("provincia"), "address;"+r.FormValue("address"), "phone;"+r.FormValue("phone"), "extra;"+r.FormValue("extra")))
			fmt.Fprint(w, respuesta)
		}
		//Envio de datos al server_ext: Modificar una tienda concreta
		if accion == "edit_tienda" {
			//Eliminamos puntos, dos puntos y puntos comas
			correct_res := libs.DeleteSplitsChars(r.FormValue("tienda"))
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/tiendas.cgi", "accion;edit_tienda", "edit_id;"+r.FormValue("id"), "tienda;"+correct_res, "provincia;"+r.FormValue("provincia"), "username;"+username, "address;"+r.FormValue("address"), "phone;"+r.FormValue("phone"), "extra;"+r.FormValue("extra")))
			fmt.Fprint(w, respuesta)
		}
		//Envio de datos al server_ext: Mostrar en una tabla las tiendas
		if accion == "get_tienda" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/tiendas.cgi", "accion;tabla_tienda", "username;"+username))
			fmt.Fprint(w, respuesta)
		}
		//Envio de datos al server_ext: Cargar una tienda concreta
		if accion == "load_tienda" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/tiendas.cgi", "accion;load_tienda", "edit_id;"+r.FormValue("load")))
			fmt.Fprint(w, respuesta)
		}
		//Envio de datos al server_ext: Generar un select de provincias para poder añadir una nueva tienda
		if accion == "show_prov" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/tiendas.cgi", "accion;show_prov", "reg;"+r.FormValue("region")))
			fmt.Fprint(w, respuesta)
		}
	}
}
