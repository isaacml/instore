package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
)

//Función que tramita el login correcto o erroneo
func auth(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	var usuario, contrasenia, salida string
	user := r.FormValue("user")
	pass := r.FormValue("pass")
	//Primero probamos la autenticacion en el servidor_interno
	db.QueryRow("SELECT user, pass FROM usuarios WHERE user = ?", user).Scan(&usuario, &contrasenia)
	//Si en la base de datos del servidor_interno estan vacios...
	if usuario == "" && contrasenia == "" {
		//Hacemos la autenticacion en el servidor_externo
		res := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverexterno"]+"/login.cgi", "user;"+user, "pass;"+pass))
		//Si la auth es correcta, guardamos los datos de usuario en el servidor_interno
		if res == "OK" {
			db_mu.Lock()
			_, err1 := db.Exec("INSERT INTO usuarios (`user`, `pass`) VALUES (?, ?)", user, pass)
			db_mu.Unlock()
			if err1 != nil {
				Error.Println(err1)
			}
		}
		salida = res
	} else {
		//Si estan bien, autenticamos con los datos en la BD del server_interno
		if user == usuario && pass == contrasenia {
			salida = "OK"
		} else {
			salida = "NOOK"
		}
	}
	//Enviamos el resultado de la autenticacion
	fmt.Fprintf(w, salida)
}
