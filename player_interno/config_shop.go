package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
	"os"
	"strings"
)

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

//Función que tramita el submit de formulario para la página(config_shop.html)
func send_orgs(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverint["serverinterno"]+"/send_orgs.cgi"))
	//Partimos las respuesta para obtener: estado (OK o NOOK) y el dominio
	gen_domain := strings.Split(respuesta, ";")
	gen := gen_domain[0]
	domain := gen_domain[1]
	if gen == "OK" {
		//OK: dominio de configuración correcto, se genera el fichero de configuración de la tienda
		config_file, err := os.Create(configShop)
		if err != nil {
			Error.Println(err)
		}
		config_file.WriteString("shopdomain = " + domain + "\n")
		//Aquí tomamos el SID que nos proporciona el formulario (action="/send_orgs.cgi?sid={{sid}}")
		for k, v := range r.Form {
			if k == "sid" {
				for _, sid := range v {
					//Una vez generado el fichero configuracion de la tienda, redirigimos a menu.html con el SID correspondiente
					_, ok := user[sid]
					if ok {
						http.Redirect(w, r, "/"+enter_page+"?"+sid, http.StatusSeeOther)
					}
				}
			}
		}
	} else {
		//NOOK: el dominio de configuración no está correcto
		for k, v := range r.Form {
			if k == "sid" {
				for _, sid := range v {
					//Redirigimos a la página de configuración(config_shop.html) con el SID correspondiente
					_, ok := user[sid]
					if ok {
						http.Redirect(w, r, "/"+shop_config_page+"?"+sid, http.StatusSeeOther)
					}
				}
			}
		}
	}
}

func additional_domains(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverint["serverinterno"]+"/send_orgs.cgi"))
	domain := strings.Split(respuesta, ";")
	dom := domain[1]
	fmt.Println(dom)
	fmt.Println(domainint)
}
