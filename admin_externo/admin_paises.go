package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
)

//Gestion de paises (paises.html)
func paises(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		accion := r.FormValue("accion")
		loadSettings(serverRoot)
		updateExpires(sid)
		//Envio de datos al server_ext: Dar de alta un nuevo país
		if accion == "pais" {
			//Eliminamos puntos, dos puntos y puntos comas
			correct_res := libs.DeleteSplitsChars(r.FormValue("pais"))
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/paises.cgi", "accion;pais", "pais;"+correct_res, "username;"+username, "almacen;"+r.FormValue("almacen")))
			fmt.Fprint(w, respuesta)
		}
		//Envio de datos al server_ext: Modificar un país concreto
		if accion == "edit_pais" {
			//Eliminamos puntos, dos puntos y puntos comas
			correct_res := libs.DeleteSplitsChars(r.FormValue("pais"))
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/paises.cgi", "accion;edit_pais", "edit_id;"+r.FormValue("id"), "almacen;"+r.FormValue("almacen"), "username;"+username, "pais;"+correct_res))
			fmt.Fprint(w, respuesta)
		}
		//Envio de datos al server_ext: Mostrar una tabla con los paises
		if accion == "get_pais" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/paises.cgi", "accion;tabla_pais", "username;"+username))
			fmt.Fprint(w, respuesta)
		}
		//Envio de datos al server_ext: Cargar un país concreto
		if accion == "load_pais" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/paises.cgi", "accion;load_pais", "edit_id;"+r.FormValue("load")))
			fmt.Fprint(w, respuesta)
		}
		//Borrar un país concreto
		if accion == "del_pais" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/paises.cgi", "accion;del_pais", "borrar;"+r.FormValue("borrar")))
			fmt.Fprint(w, respuesta)
		}
		//Envio de datos al server_ext: Generar un select de almacenes para poder añadir un nuevo pais
		if accion == "show_almacen" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/paises.cgi", "accion;show_almacen", "username;"+username))
			fmt.Fprint(w, respuesta)
		}
	}
}
