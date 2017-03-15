package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
)

// funcion q tramita el login correcto o erroneo
func login_tienda(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	loadSettings(serverRoot)
	//SE PASAN LAS VARIABLES POST AL SERVIDOR EXTERNO PARA LA AUTENTICACION
	respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverexterno"]+"/login.cgi", "user;"+r.FormValue("user"), "pass;"+r.FormValue("pass")))
	fmt.Fprint(w, respuesta)
}
