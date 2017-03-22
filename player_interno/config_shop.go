package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
	"os"
)

//Comprueba si el fichero de configuracion de la tienda existe o no
func check_config(w http.ResponseWriter, r *http.Request) {
	var existe string
	_, err := os.Stat(configShop)
	if err != nil {
		if os.IsNotExist(err) {
			existe = "NOOK"
		}
	} else {
		existe = "OK"
	}
	fmt.Fprint(w, existe)
}

//Funcion que va a recoger los valores de los selects y mostrarlos
func get_orgs(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	accion := r.FormValue("action")

	if accion == "entidades" {
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverint["serverinterno"]+"/transf_orgs.cgi", "action;entidad", "user;"+username))
		fmt.Fprint(w, respuesta)
	}
	if accion == "almacenes" {
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverint["serverinterno"]+"/transf_orgs.cgi", "action;almacen", "entidad;"+r.FormValue("entidad")))
		fmt.Fprint(w, respuesta)
	}
	if accion == "paises" {
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverint["serverinterno"]+"/transf_orgs.cgi", "action;pais", "almacen;"+r.FormValue("almacen")))
		fmt.Fprint(w, respuesta)
	}
	if accion == "regiones" {
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverint["serverinterno"]+"/transf_orgs.cgi", "action;region", "pais;"+r.FormValue("pais")))
		fmt.Fprint(w, respuesta)
	}
	if accion == "provincias" {
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverint["serverinterno"]+"/transf_orgs.cgi", "action;provincia", "region;"+r.FormValue("region")))
		fmt.Fprint(w, respuesta)
	}
	if accion == "tiendas" {
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverint["serverinterno"]+"/transf_orgs.cgi", "action;tienda", "provincia;"+r.FormValue("provincia")))
		fmt.Fprint(w, respuesta)
	}
	if accion == "cod_tienda" {
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverint["serverinterno"]+"/transf_orgs.cgi", "action;cod_tienda", "tienda;"+r.FormValue("tienda")))
		fmt.Fprint(w, respuesta)
	}
}

//Funcion que que toma los valores del formulario (config_shop.html) tras ser enviados
func get_config_shop(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	entidad := r.FormValue("entidad")
	almacen := r.FormValue("almacen")
	pais := r.FormValue("pais")
	region := r.FormValue("region")
	provincia := r.FormValue("provincia")
	tienda := r.FormValue("tienda")
	dominio := r.FormValue("dominio")
	fmt.Println(entidad, almacen, pais, region, provincia, tienda, dominio)
	if entidad != "" || almacen != "" || pais != "" || region != "" || provincia != "" || tienda != "" {
		fmt.Println(dominio)
	} else {
		fmt.Fprint(w, "<span style='color: #FF0303'>Faltan campos por llenar</span>")
	}
}
