package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
	"time"
)

//Intermediario de organizaciones entre la tienda y el servidor externo
func transf_orgs(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var respuesta string
	accion := r.FormValue("action")
	//Enviamos el username al servidor interno
	if accion == "entidad" {
		respuesta = fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverexterno"]+"/config_shop.cgi", "action;entidad", "username;"+r.FormValue("user")))
	}
	//Enviamos la entidad al servidor interno
	if accion == "almacen" {
		respuesta = fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverexterno"]+"/config_shop.cgi", "action;almacen", "entidad;"+r.FormValue("entidad")))
	}
	//Enviamos el almacen al servidor interno
	if accion == "pais" {
		respuesta = fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverexterno"]+"/config_shop.cgi", "action;pais", "almacen;"+r.FormValue("almacen")))
	}
	//Enviamos el pais al servidor interno
	if accion == "region" {
		respuesta = fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverexterno"]+"/config_shop.cgi", "action;region", "pais;"+r.FormValue("pais")))
	}
	//Enviamos la región al servidor interno
	if accion == "provincia" {
		respuesta = fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverexterno"]+"/config_shop.cgi", "action;provincia", "region;"+r.FormValue("region")))
	}
	//Enviamos la provincia al servidor interno
	if accion == "tienda" {
		respuesta = fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverexterno"]+"/config_shop.cgi", "action;tienda", "provincia;"+r.FormValue("provincia")))
	}
	//Enviamos la tienda al servidor interno
	if accion == "cod_tienda" {
		respuesta = fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverexterno"]+"/config_shop.cgi", "action;cod_tienda", "tienda;"+r.FormValue("tienda")))
	}
	fmt.Fprint(w, respuesta)
}

//Acciones realizadas por parte del servidor interno
func acciones(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var output string
	accion := r.FormValue("action")
	//Envia user y pass hacia el server_externo
	if accion == "login_tienda" {
		//Se pasan las variables de autenticación
		output = fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverexterno"]+"/login.cgi", "user;"+r.FormValue("user"), "pass;"+r.FormValue("pass")))
	}
	//Pasa el nombre de usuario al servidor externo, nos devuelve los permisos para ese usuario
	if accion == "bitmaps" {
		output = libs.GenerateFORM(serverext["serverexterno"]+"/acciones.cgi", "accion;bitmap_perm", "user;"+r.FormValue("user"))
	}
	//Guardar dominios de la tienda en el fichero de configuracion.
	if accion == "save_domain" {
		output = libs.GenerateFORM(serverext["serverexterno"] + "/send_shop.cgi")
	}
	//Envia los dominios tomados de la tienda(configshop.reg) hacia el server externo
	if accion == "send_domains" {
		output = libs.GenerateFORM(serverext["serverexterno"]+"/recoger_dominio.cgi", "dominios;"+r.FormValue("dominios"))
	}
	//Envia la entidad tomada del dominio de la tienda, para obtener el estado de dicha entidad
	if accion == "check_entidad" {
		output = libs.GenerateFORM(serverext["serverexterno"]+"/acciones.cgi", "accion;check_entidad", "ent;"+r.FormValue("ent"))
	}
	fmt.Fprint(w, output)
}

//Recibe las peticiones de mensajes/publicidad por parte de la tienda
func publi_msg(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	accion := r.FormValue("action")
	//Guarda los mensajes en la base de datos interna
	if accion == "MsgFiles" {
		var existe string
		nombre_fichero := r.FormValue("fichero")
		timestamp := time.Now().Unix()
		//Se comprueba que la existencia en la tienda se corresponde con la existencia en el server interno
		db.QueryRow("SELECT existe FROM mensaje WHERE fichero=?", nombre_fichero).Scan(&existe)
		//existe = vacio --> No ha sido insertado nunca
		if existe == "" {
			//Por tanto, lo insertamos con el existe en N para que el bucle BuscarNuevosFicheros() pueda localizarlo
			publi, err := db.Prepare("INSERT INTO mensaje (`fichero`, `existe`, `timestamp`) VALUES (?,?,?)")
			if err != nil {
				Error.Println(err)
			}
			db_mu.Lock()
			_, err1 := publi.Exec(nombre_fichero, "N", timestamp)
			db_mu.Unlock()
			if err1 != nil {
				Error.Println(err1)
			}
		} else {
			//Se compara la existencia del fichero.
			if existe != r.FormValue("existencia") {
				//Se le manda una orden de descarga de fichero a la tienda
				fmt.Fprint(w, "Descarga")
			}
		}
	}
	//Guarda la publicidad en la base de datos interna
	if accion == "PubliFiles" {
		var existe string
		nombre_fichero := r.FormValue("fichero")
		gap := r.FormValue("gap")
		f_ini := r.FormValue("fecha_ini")
		timestamp := time.Now().Unix()
		//Se comprueba que la existencia del fichero en la tienda se corresponde con la existencia en el server interno
		db.QueryRow("SELECT existe FROM publi WHERE fichero=?", nombre_fichero).Scan(&existe)
		//existe = vacio --> No ha sido insertado nunca
		if existe == "" {
			//Por tanto, lo insertamos con el existe en N para que el bucle BuscarNuevosFicheros() pueda localizarlo
			publi, err := db.Prepare("INSERT INTO publi (`fichero`, `existe`, fecha_ini, `timestamp`, `gap`) VALUES (?,?,?,?,?)")
			if err != nil {
				Error.Println(err)
			}
			db_mu.Lock()
			_, err1 := publi.Exec(nombre_fichero, "N", f_ini, timestamp, gap)
			db_mu.Unlock()
			if err1 != nil {
				Error.Println(err1)
			}
		} else {
			//Se compara la existencia del fichero.
			if existe != r.FormValue("existencia") {
				//Se le manda una orden de descarga de fichero a la tienda
				fmt.Fprint(w, "Descarga")
			}
		}
	}
}
