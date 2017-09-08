package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
	"github.com/isaacml/instore/winamp"
	"os"
	"strings"
)

func acciones(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	accion := r.FormValue("accion")
	//Obtener los valores de bitmap para el usuario de la tienda
	//Enviamos el nombre del usuario al server_interno y este lo pasará al server_externo
	if accion == "bitmaps"{
		respuesta := libs.GenerateFORM(serverint["serverinterno"]+"/bitmaps.cgi", "user;"+username)
		fmt.Fprint(w, respuesta)
	}
	//Comprueba si el fichero de configuracion de la tienda existe o no
	if accion == "check_config"{
		var existe string
		existencia := libs.Existencia(configShop)
		if existencia == true {
			existe = "OK"
		}else{
			existe = "NOOK"
		}
		fmt.Fprint(w, existe)	
	}
}

//Mensajes Instantaneos
func mensajesInstantaneos(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var output string //variable para imprimir los datos hacia JavaScript
	var msg_instantaneo string
	//Generar un listado de Mensajes
	if r.FormValue("action") == "mensajes" {
		//Abrimos el directorio de mensajes(MessagesShop) 
		file, err := os.Open(msg_files_location)
		defer file.Close()
		if err != nil {
			Error.Println(err)
			return
		}
		ficheros, err := file.Readdir(0)
		if err != nil {
			Error.Println(err)
			return
		}
		for key, val := range ficheros {
			//Tomamos solamente ficheros MP3
			if strings.Contains(val.Name(), ".mp3"){
				//Formamos el select
				output += fmt.Sprintf("<option value='%s'>%s</option>", val.Name(), val.Name())
				//Guardamos el nombre del primer mensaje
				if key == 0 {
					msg_instantaneo = val.Name()
				}
			}
		}
		//Mostrar msg: muestra el nombre del mensaje por primera vez
		output += fmt.Sprintf("#:#<span style='color: #006400'>Mensaje seleccionado: </span>"  + msg_instantaneo)
	}
	//Estado de mensaje
	if r.FormValue("action") == "status" {
		msg_instantaneo = r.FormValue("instantaneos")
		output = fmt.Sprintf("<span style='color: #006400'>Mensaje seleccionado: </span>"  + msg_instantaneo)
	}
	//Recibe el mensaje instantaneo y lo procesa
	if r.FormValue("action") == "send" {
		var win winamp.Winamp
		//Reproducimos el mensaje instantaneo
		win.PlayFFplay(msg_files_location + r.FormValue("instantaneos"))
	}
	fmt.Fprint(w, output)
}