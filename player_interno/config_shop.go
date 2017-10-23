package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	//"io/ioutil"
	//"bufio"
	"net/http"
	"os"
	"strings"
)

//Acciones que se van a llevar acabo para configurar el dominio de la tienda
func config_shop(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	accion := r.FormValue("accion")
	//Guarda el dominio principal(de un usuario) en el fichero de configuracion
	if accion == "gen_config_file" {
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverint["serverinterno"]+"/acciones.cgi", "action;save_domain"))
		//Partimos las respuesta para obtener: estado (OK o NOOK) y el dominio
		gen_domain := strings.Split(respuesta, ";")
		gen := gen_domain[0]
		domain := gen_domain[1]
		if gen == "OK" {
			//OK: dominio de configuración correcto, se genera el fichero.
			config_file, err := os.Create(configShop)
			if err != nil {
				Error.Println(err)
			}
			defer config_file.Close()
			//Guardamos el dominio
			config_file.WriteString("shopdomain = " + domain + "\n")
			//Checkeamos el SID que nos proporciona el formulario
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
	//Añade nuevos dominios adiccionales a la tienda
	if accion == "extra_domains" {
		cont := 0
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverint["serverinterno"]+"/acciones.cgi", "action;save_domain"))
		//Partimos las respuesta para obtener: estado (OK o NOOK) y el dominio
		gen_domain := strings.Split(respuesta, ";")
		gen := gen_domain[0]
		domain := gen_domain[1]
		if gen == "OK" {
			domainint := make(map[string]string) //Mapa que guarda el dominio de la tienda
			loadSettings(configShop, domainint)
			for _, val := range domainint {
				if val == domain {
					cont++
				}
			}
			//contador = 0: No existe dominio en el fichero
			if cont == 0 {
				fr, err := os.OpenFile(configShop, os.O_APPEND, 0666)
				defer fr.Close()
				if err == nil {
					//escribimos nuestro dominio extra
					fr.WriteString("extradomain = " + domain + "\n")
				}
			}
		}
	}
}

//Función que tramita el submit de formulario para la página(config_shop.html)
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
