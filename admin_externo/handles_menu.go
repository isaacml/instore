package main

import (
	"net/http"
	"github.com/isaacml/instore/libs"
	"fmt"
)

func zona_publi(w http.ResponseWriter, r *http.Request) {
	respuesta := libs.GenerateFORM(serverext["serverroot"]+"/check_actions.cgi", "action;create_pub", "user;"+username)
	if respuesta == "0" {
		fmt.Fprint(w, "NOOK")
	}
}
func zona_music(w http.ResponseWriter, r *http.Request) {
	respuesta := libs.GenerateFORM(serverext["serverroot"]+"/check_actions.cgi", "action;prog_mus", "user;"+username)
	if respuesta == "0" {
		fmt.Fprint(w, "NOOK")
	}
}
