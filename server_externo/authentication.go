package main

import (
	"fmt"
	"net/http"
)

// This function could be used to access to a Database for user/pass authentication procedure
func authentication(user, pass string) bool {
	var username, password string
	db_mu.RLock()
	query2, err := db.Query("SELECT user, pass FROM usuarios WHERE user = ?", user)
	db_mu.RUnlock()
	if err != nil {
		Error.Println(err)
	}
	for query2.Next() {
		err = query2.Scan(&username, &password)
		if err != nil {
			Error.Println(err)
		}
	}
	if user == username && pass == password {
		return true
	} else {
		return false
	}
}

//Función que tramita el login correcto o erroneo
func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	usuario := r.FormValue("user")
	pass := r.FormValue("pass")

	if authentication(usuario, pass) {
		fmt.Fprintf(w, "OK")
	} else {
		fmt.Println("Login incorrecto")
		fmt.Fprintf(w, "NOOK")
	}
}
