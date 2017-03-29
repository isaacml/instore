package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
)

//Pasa los destinos al servidor externo
func destino(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	accion := r.FormValue("action")

	if accion == "destinos" {
		respuesta := libs.GenerateFORM(serverext["serverexterno"]+"/destino.cgi", "action;destinos", "userAdmin;"+r.FormValue("userAdmin"))
		fmt.Fprint(w, respuesta)
	}
	if accion == "entidades" {
		respuesta := libs.GenerateFORM(serverext["serverexterno"]+"/destino.cgi", "action;entidades", "id_entidad;"+r.FormValue("id_entidad"), "userAdmin;"+r.FormValue("userAdmin"))
		fmt.Fprint(w, respuesta)
	}
	if accion == "almacenes" {
		respuesta := libs.GenerateFORM(serverext["serverexterno"]+"/destino.cgi", "action;almacenes", "id_almacen;"+r.FormValue("id_almacen"), "userAdmin;"+r.FormValue("userAdmin"))
		fmt.Fprint(w, respuesta)
	}
	if accion == "paises" {
		respuesta := libs.GenerateFORM(serverext["serverexterno"]+"/destino.cgi", "action;paises", "id_pais;"+r.FormValue("id_pais"), "userAdmin;"+r.FormValue("userAdmin"))
		fmt.Fprint(w, respuesta)
	}
	if accion == "regiones" {
		respuesta := libs.GenerateFORM(serverext["serverexterno"]+"/destino.cgi", "action;regiones", "id_region;"+r.FormValue("id_region"), "userAdmin;"+r.FormValue("userAdmin"))
		fmt.Fprint(w, respuesta)
	}
	if accion == "provincias" {
		respuesta := libs.GenerateFORM(serverext["serverexterno"]+"/destino.cgi", "action;provincias", "id_provincia;"+r.FormValue("id_provincia"), "userAdmin;"+r.FormValue("userAdmin"))
		fmt.Fprint(w, respuesta)
	}
	if accion == "tiendas" {
		respuesta := libs.GenerateFORM(serverext["serverexterno"]+"/destino.cgi", "action;tiendas", "id_tienda;"+r.FormValue("id_tienda"), "userAdmin;"+r.FormValue("userAdmin"))
		fmt.Fprint(w, respuesta)
	}
}
