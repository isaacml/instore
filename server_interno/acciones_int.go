package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
	"time"
)

//Pasa el nombre de usuario al servidor externo
func bitmaps(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	respuesta := libs.GenerateFORM(serverext["serverexterno"]+"/bitmap_actions.cgi", "user;"+r.FormValue("user"))
	fmt.Fprint(w, respuesta)
}

//Intermediario para enviar el dominio de la tienda
func send_domain(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	respuesta := libs.GenerateFORM(serverext["serverexterno"]+"/recoger_dominio.cgi", "dominio;"+r.FormValue("dominio"))
	fmt.Fprint(w, respuesta)
}

//Recibe la peticion de mensaje por parte del player_interno(tienda)
func downloadMsgFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	nombre_fichero := r.FormValue("fichero")
	timestamp := time.Now().Unix()
	if r.FormValue("existencia") == "N" {
		var cont int
		//Comprobamos si el fichero está insertado en BD.
		err := db.QueryRow("SELECT count(id) FROM mensaje WHERE fichero=?", nombre_fichero).Scan(&cont)
		if err != nil {
			Error.Println(err)
		}
		//cont = 0 --> No ha sido insertado nunca
		if cont == 0 {
			//Lo insertamos con el existe en N
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
		}
	}
}

//Recibe la peticion de publicidad por parte del player_interno(tienda)
func downloadPubliFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	nombre_fichero := r.FormValue("fichero")
	timestamp := time.Now().Unix()
	if r.FormValue("existencia") == "N" { //existencia del fichero publi, mandado por el player_interno
		var cont int
		//Comprobamos si el fichero está insertado en BD.
		err := db.QueryRow("SELECT count(id) FROM publi WHERE fichero=?", nombre_fichero).Scan(&cont)
		if err != nil {
			Error.Println(err)
		}
		//cont = 0 --> No ha sido insertado nunca
		if cont == 0 {
			//Lo insertamos con el existe en N
			publi, err := db.Prepare("INSERT INTO publi (`fichero`, `existe`, `timestamp`) VALUES (?,?,?)")
			if err != nil {
				Error.Println(err)
			}
			db_mu.Lock()
			_, err1 := publi.Exec(nombre_fichero, "N", timestamp)
			db_mu.Unlock()
			if err1 != nil {
				Error.Println(err1)
			}
		}
	}
}
