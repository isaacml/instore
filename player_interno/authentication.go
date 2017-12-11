package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
)

func authentication(username, password string) string {
	var usuario, contrasenia, salida string
	//Primero se hace una autenticacion interna
	err := db.QueryRow("SELECT user, pass FROM usuarios WHERE user = ?", username).Scan(&usuario, &contrasenia)
	if err != nil {
		Error.Println(err)
	}
	//Si en la base de datos de la tienda estan vacios...
	if usuario == "" && contrasenia == "" {
		libs.LoadSettingsWin(serverRoot, settings)
		//Probamos a autenticar con los datos del server interno de la tienda
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(settings["serverinterno"]+"/auth.cgi", "user;"+username, "pass;"+password))
		//Si la auth es correcta, guardamos los datos de usuario en la BD de la tienda
		if respuesta == "OK" {
			db_mu.Lock()
			_, err1 := db.Exec("INSERT INTO usuarios (`user`, `pass`) VALUES (?, ?)", username, password)
			db_mu.Unlock()
			if err1 != nil {
				Error.Println(err1)
			}
		}
		salida = respuesta
	} else {
		//Se procede a la auth en la propia tienda
		if username == usuario && password == contrasenia {
			salida = "OK"
		} else {
			salida = "NOOK"
		}
	}
	return salida
}
