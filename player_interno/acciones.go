package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
	"github.com/isaacml/instore/winamp"
	"os"
	"strings"
)

//Pasamos el nombre de usuario al servidor interno para que lo transfiera al server externo
func bitmaps(w http.ResponseWriter, r *http.Request) {
	respuesta := libs.GenerateFORM(serverint["serverinterno"]+"/bitmaps.cgi", "user;"+username)
	fmt.Fprint(w, respuesta)
}

//Mensajes Instantaneos
func mensajesInstantaneos(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var output string //variable para imprimir los datos hacia JavaScript
	
	//Generar un listado de Mensajes
	if r.FormValue("action") == "mensajes" {
		var msg_instantaneo string
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
				if key == 1 {
					msg_instantaneo = val.Name()
				}
			}
		}
		//Mostrar msg: muestra el nombre del mensaje por primera vez
		output += fmt.Sprintf("#:#<span style='color: #006400'>Mensaje seleccionado: </span>"  + msg_instantaneo)
	}
	//Estado de mensaje
	if r.FormValue("action") == "status" {
		output = fmt.Sprintf("<span style='color: #006400'>Mensaje seleccionado: </span>"  + r.FormValue("instantaneos"))
	}
	//Recibe el mensaje instantaneo y lo procesa
	if r.FormValue("action") == "send" {
		var win winamp.Winamp
		//Reproducimos el mensaje instantaneo
		win.PlayFFplay(msg_files_location + r.FormValue("instantaneos"))
		//Mensaje de reproduccion acabada
		output = fmt.Sprintf("<span style='color: #FF0303'>Reproducción Finalizada</span>")
	}
	fmt.Fprint(w, output)
}