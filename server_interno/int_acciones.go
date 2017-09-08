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
	var existe string
	nombre_fichero := r.FormValue("fichero")
	timestamp := time.Now().Unix()
	//Se comprueba que la existencia en la tienda se corresponde con la existencia en el server interno
	db.QueryRow("SELECT existe FROM mensaje WHERE fichero=?", nombre_fichero).Scan(&existe)
	//existe = vacio --> No ha sido insertado nunca
	if existe == "" {
		//Por tanto, lo insertamos con el existe en N para que el bucle BuscarNuevosFicheros() pueda localizarlo
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
	} else {
		//Se compara la existencia del fichero.
		if existe != r.FormValue("existencia") {
			//Se le manda una orden de descarga de fichero a la tienda
			fmt.Fprint(w, "Descarga")
		}
	}
}

//Recibe la peticion de publicidad por parte del player_interno(tienda)
func downloadPubliFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var existe string
	nombre_fichero := r.FormValue("fichero")
	gap := r.FormValue("gap")
	f_ini := r.FormValue("fecha_ini")
	timestamp := time.Now().Unix()
	//Se comprueba que la existencia del fichero en la tienda se corresponde con la existencia en el server interno
	db.QueryRow("SELECT existe FROM publi WHERE fichero=?", nombre_fichero).Scan(&existe)
	//existe = vacio --> No ha sido insertado nunca
	if existe == "" {
		//Por tanto, lo insertamos con el existe en N para que el bucle BuscarNuevosFicheros() pueda localizarlo
		publi, err := db.Prepare("INSERT INTO publi (`fichero`, `existe`, fecha_ini, `timestamp`, `gap`) VALUES (?,?,?,?,?)")
		if err != nil {
			Error.Println(err)
		}
		db_mu.Lock()
		_, err1 := publi.Exec(nombre_fichero, "N", f_ini, timestamp, gap)
		db_mu.Unlock()
		if err1 != nil {
			Error.Println(err1)
		}
	} else {
		//Se compara la existencia del fichero.
		if existe != r.FormValue("existencia") {
			//Se le manda una orden de descarga de fichero a la tienda
			fmt.Fprint(w, "Descarga")
		}
	}
}