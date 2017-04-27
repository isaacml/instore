package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
	"os"
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

//Recibe la peticion de publicidad por parte del player_interno(tienda)
func downloadPubliFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	nombre_fichero := r.FormValue("fichero")
	timestamp := time.Now().Unix()
	if r.FormValue("existencia") == "N" {
		var existe bool
		var cont int
		//Comprobar que existe el fichero en el directorio de publicidad del server interno
		_, err := os.Stat(publi_files_location + nombre_fichero)
		if err != nil {
			if os.IsNotExist(err) {
				existe = false
			}
		} else {
			existe = true
		}
		//Comprobamos si el fichero está insertado en BD. cont = 0 --> No ha sido insertado nunca
		err = db.QueryRow("SELECT count(id) FROM publi WHERE fichero=?", nombre_fichero).Scan(&cont)
		if err != nil {
			Error.Println(err)
		}
		//Fichero de publicidad no existe en directorio del server interno
		if existe == false {
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
		} else { //Fichero de publicidad existe en directorio del server interno
			if cont == 0 {
				//Lo insertamos con el existe en Y
				publi, err := db.Prepare("INSERT INTO publi (`fichero`, `existe`, `timestamp`) VALUES (?,?,?)")
				if err != nil {
					Error.Println(err)
				}
				db_mu.Lock()
				_, err1 := publi.Exec(nombre_fichero, "Y", timestamp)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
				}
			}
		}
	}

}
