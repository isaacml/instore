package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
	"os"
	//"strings"
)

var redirect bool

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

	if accion == "enviar_sid" {
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverint["serverinterno"]+"/transf_orgs.cgi", "action;enviar_sid", "sid_id;"+r.FormValue("sid")))
		fmt.Fprint(w, respuesta)
	}
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

/*
//Función que envía todos los campos POST (config_shop.html)
func send_orgs(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	fmt.Println(sid)
	http.Redirect(w, r, "/"+enter_page+"?"+sid, http.StatusMovedPermanently)
	fmt.Fprint(w, "Quiero salir de aqui")

		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverint["serverinterno"]+"/send_orgs.cgi"))
		//Partimos las respuesta para obtener: estado (OK o NOOK) y el dominio
		gen_domain := strings.Split(respuesta, ";")
		gen := gen_domain[0]
		domain := gen_domain[1]
		if gen == "OK" {
			config_file, err := os.Create(configShop)
			if err != nil {
				Error.Println(err)
			}
			config_file.WriteString("shop_domain = " + domain)
			http.Redirect(w, r, "/"+enter_page+"?"+sid, http.StatusFound)
			fmt.Fprint(w, "Quiero salir de aqui")
		} else {
			//output := "<span style='color: #FF0303'>Faltan campos por llenar</span>"
			//fmt.Fprint(w, output)
		}

}
*/

// función que tramita el logout de la session
func send_orgs(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	sid := r.FormValue("sid")
	for k, v := range r.Form {
		fmt.Println(k, v)
	}
	fmt.Println(sid)
	_, ok := user[sid]
	fmt.Println(ok)
	if ok {
		http.Redirect(w, r, "/"+enter_page+"?"+sid, http.StatusSeeOther)
	}
}
