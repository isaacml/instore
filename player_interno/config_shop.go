package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
	"os"
	"strings"
	"time"
)

//Función que tramita el formulario de la página(config_shop.html)
//Toma los valores, los envia al server_externo, recoge la respuesta y la muestra
func get_orgs(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var respuesta string
	accion := r.FormValue("action")
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		libs.LoadSettingsWin(serverRoot, settings)
		if accion == "entidades" {
			respuesta = fmt.Sprintf("%s", libs.GenerateFORM(settings["serverinterno"]+"/transf_orgs.cgi", "action;entidad", "user;"+user[sid]))
		}
		if accion == "almacenes" {
			respuesta = fmt.Sprintf("%s", libs.GenerateFORM(settings["serverinterno"]+"/transf_orgs.cgi", "action;almacen", "entidad;"+r.FormValue("entidad")))
		}
		if accion == "paises" {
			respuesta = fmt.Sprintf("%s", libs.GenerateFORM(settings["serverinterno"]+"/transf_orgs.cgi", "action;pais", "almacen;"+r.FormValue("almacen")))
		}
		if accion == "regiones" {
			respuesta = fmt.Sprintf("%s", libs.GenerateFORM(settings["serverinterno"]+"/transf_orgs.cgi", "action;region", "pais;"+r.FormValue("pais")))
		}
		if accion == "provincias" {
			respuesta = fmt.Sprintf("%s", libs.GenerateFORM(settings["serverinterno"]+"/transf_orgs.cgi", "action;provincia", "region;"+r.FormValue("region")))
		}
		if accion == "tiendas" {
			respuesta = fmt.Sprintf("%s", libs.GenerateFORM(settings["serverinterno"]+"/transf_orgs.cgi", "action;tienda", "provincia;"+r.FormValue("provincia")))
		}
		if accion == "cod_tienda" {
			respuesta = fmt.Sprintf("%s", libs.GenerateFORM(settings["serverinterno"]+"/transf_orgs.cgi", "action;cod_tienda", "tienda;"+r.FormValue("tienda")))
		}
		fmt.Fprint(w, respuesta)
	}
}

//Acciones que se van a llevar acabo para configurar el dominio de la tienda
func config_shop(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	accion := r.FormValue("accion")
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		timestamp := time.Now().Unix()
		//Guarda el dominio principal(de un usuario) en el fichero de configuracion
		if accion == "gen_config_file" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(settings["serverinterno"]+"/acciones.cgi", "action;save_domain"))
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
								//Guardamos en la base de datos interna de la tienda, el dominio y la ultima conexion
								shop, err := db.Prepare("INSERT INTO tienda (`dominio`, `last_connect`) VALUES (?,?)")
								if err != nil {
									Error.Println(err)
								}
								db_mu.Lock()
								_, err1 := shop.Exec(domain, timestamp)
								db_mu.Unlock()
								if err1 != nil {
									Error.Println(err1)
								}
								block = false
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
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(settings["serverinterno"]+"/acciones.cgi", "action;save_domain"))
			//Partimos las respuesta para obtener: estado (OK o NOOK) y el dominio
			gen_domain := strings.Split(respuesta, ";")
			gen := gen_domain[0]
			domain := gen_domain[1]
			if gen == "OK" {
				domainint := make(map[string]string) //Mapa que guarda el dominio de la tienda
				libs.LoadDomains(configShop, domainint)
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
			//Checkeamos el sid
			for k, v := range r.Form {
				if k == "sid" {
					for _, sid := range v {
						_, ok := user[sid]
						if ok {
							//Redirigimos a la página de añadir dominios extra(adddomain.html)
							http.Redirect(w, r, "/dominios.html?"+sid, http.StatusSeeOther)
						}
					}
				}
			}
		}
	}
}
