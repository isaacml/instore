package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
)

// enviamos user y pass de autenticación al servidor interno
func login_tienda(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//SE PASAN LAS VARIABLES DE AUTENTICACION
	respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverexterno"]+"/login.cgi", "user;"+r.FormValue("user"), "pass;"+r.FormValue("pass")))
	fmt.Fprint(w, respuesta)
}

func transf_orgs(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	accion := r.FormValue("action")
	//Enviamos el username al servidor interno
	if accion == "entidad" {
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverexterno"]+"/config_shop.cgi", "action;entidad", "username;"+r.FormValue("user")))
		fmt.Fprint(w, respuesta)
	}
	//Enviamos la entidad al servidor interno
	if accion == "almacen" {
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverexterno"]+"/config_shop.cgi", "action;almacen", "entidad;"+r.FormValue("entidad")))
		fmt.Fprint(w, respuesta)
	}
	//Enviamos el almacen al servidor interno
	if accion == "pais" {
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverexterno"]+"/config_shop.cgi", "action;pais", "almacen;"+r.FormValue("almacen")))
		fmt.Fprint(w, respuesta)
	}
	//Enviamos el pais al servidor interno
	if accion == "region" {
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverexterno"]+"/config_shop.cgi", "action;region", "pais;"+r.FormValue("pais")))
		fmt.Fprint(w, respuesta)
	}
	//Enviamos la región al servidor interno
	if accion == "provincia" {
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverexterno"]+"/config_shop.cgi", "action;provincia", "region;"+r.FormValue("region")))
		fmt.Fprint(w, respuesta)
	}
	//Enviamos la provincia al servidor interno
	if accion == "tienda" {
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverexterno"]+"/config_shop.cgi", "action;tienda", "provincia;"+r.FormValue("provincia")))
		fmt.Fprint(w, respuesta)
	}
	//Enviamos la tienda al servidor interno
	if accion == "tienda" {
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverexterno"]+"/config_shop.cgi", "action;cod_tienda", "tienda;"+r.FormValue("tienda")))
		fmt.Fprint(w, respuesta)
	}
}
