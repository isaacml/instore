package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
)

//Gestion de regiones (regiones.html)
func regiones(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	accion := r.FormValue("accion")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		//Envio de datos al server_ext: Dar de alta una nueva region
		if accion == "region"{
			//Eliminamos puntos, dos puntos y puntos comas
			correct_res := libs.DeleteSplitsChars(r.FormValue("region"))
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/regiones.cgi", "accion;region", "pais;"+r.FormValue("pais"), "username;"+username, "region;"+correct_res))
			if respuesta == "OK" {
				good = "Región añadida correctamente"
				fmt.Fprintf(w, "<div class='form-group text-success'>%s</div>", good)
			} else {
				fmt.Fprint(w, respuesta)
			}
		}
		//Envio de datos al server_ext: Modificar una región concreta
		if accion == "edit_region"{
			//Eliminamos puntos, dos puntos y puntos comas
			correct_res := libs.DeleteSplitsChars(r.FormValue("region"))
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/regiones.cgi", "accion;edit_region", "edit_id;"+r.FormValue("id"), "region;"+correct_res, "username;"+username, "pais;"+r.FormValue("pais")))
			if respuesta == "OK" {
				good = "Región modificada correctamente"
				fmt.Fprintf(w, "<div class='form-group text-success'>%s</div>", good)
			} else {
				fmt.Fprint(w, respuesta)
			}
		}
		//Envio de datos al server_ext: Mostrar una tabla de las regiones
		if accion == "get_region"{
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/regiones.cgi", "accion;tabla_region", "username;"+username))
			fmt.Fprint(w, respuesta)
		}
		//Envio de datos al server_ext: Cargar una región concreta
		if accion == "load_region"{
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/regiones.cgi", "accion;load_region", "edit_id;"+r.FormValue("load")))
			fmt.Fprint(w, respuesta)
		}
		//Envio de datos al server_ext: Generar un select de paises para poder añadir una nueva region
		if accion == "region_pais"{
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/regiones.cgi", "accion;region_pais", "username;"+username))
			fmt.Fprint(w, respuesta)
		}
	}
}