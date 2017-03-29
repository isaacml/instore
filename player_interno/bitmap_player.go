package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
)

//Pasamos el nombre de usuario al servidor interno para que lo transfiera al server externo
func bitmaps(w http.ResponseWriter, r *http.Request) {
	respuesta := libs.GenerateFORM(serverint["serverinterno"]+"/bitmaps.cgi", "user;"+username)
	fmt.Fprint(w, respuesta)
	fmt.Println(respuesta)
}
