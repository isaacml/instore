package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
)

//Pasa el nombre de usuario al servidor externo
func bitmaps(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	respuesta := libs.GenerateFORM(serverext["serverexterno"]+"/bitmap_actions.cgi", "user;"+r.FormValue("user"))
	fmt.Fprint(w, respuesta)
}
