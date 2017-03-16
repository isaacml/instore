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

//Mostrar las entidades para un usuario concreto
func entidades(w http.ResponseWriter, r *http.Request) {
	respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverint["serverinterno"]+"/send_user.cgi", "user;"+username))
	fmt.Fprint(w, respuesta)
}
