package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"github.com/isaacml/instore/winamp"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

//Bitmap de música no cifrada: valores posibles 0 o 1
//0: solo se puede escuchar musica no cifrada
//1: tanto música cifrada como no cifrada
var st_music int

func acciones(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	accion := r.FormValue("accion")
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		//Obtener los valores de bitmap para el usuario de la tienda
		//Enviamos el nombre del usuario al server_interno y este lo pasará al server_externo
		if accion == "bitmaps" {
			respuesta := libs.GenerateFORM(settings["serverinterno"]+"/acciones.cgi", "action;bitmaps", "user;"+user[sid])
			bit := strings.Split(respuesta, ";")
			db_mu.Lock()
			if len(bit) > 1 {
				if libs.ToInt(bit[3]) == 0 {
					st_music = 0
				} else {
					st_music = 1
				}
			}
			db_mu.Unlock()
			fmt.Fprint(w, respuesta)
		}
		//Comprueba si el fichero de configuracion de la tienda existe o no
		if accion == "check_config" {
			var existe string
			existencia := libs.Existencia(configShop)
			if existencia == true {
				existe = "OK"
			} else {
				existe = "NOOK"
			}
			fmt.Fprint(w, existe)
		}
		//Muestra los datos de configuracion de la tienda: dominio principal y extras
		if accion == "dataConfig" {
			var dominios string
			domainint := make(map[string]string) //Mapa que guarda el dominio de la tienda
			libs.LoadDomains(configShop, domainint)
			for key, val := range domainint {
				if key == "shopdomain" {
					dominios += fmt.Sprintf("<tr><th>Dominio Principal:</th><td>&nbsp;</td><td>%s</td></tr>", val)
				} else {
					dominios += fmt.Sprintf("<tr><th>Dominio Secundario:</th><td>&nbsp;</td><td>%s</td></tr>", val)
				}
			}
			fmt.Fprint(w, dominios)
		}
		//Elimina el fichero de configuracion de la tienda para volver a ser configurada
		if accion == "reconfigure" {
			err := os.Remove(configShop)
			if err != nil {
				Error.Println(err)
			}
			db_mu.Lock()
			block = true
			db_mu.Unlock()
		}
		//Muestra cada segundo("setInterval") el estado de la tienda
		if accion == "estado_de_tienda" {
			output := "<tr><td>Conexión de la tienda: </td><td>&nbsp;</td>"
			if block == true {
				output += "<td class='text-danger'> Desactivada</td></tr>"
			} else {
				output += "<td class='text-success'> Activada</td></tr>"
			}
			fmt.Fprint(w, output)
		}
		//Recoge de SettingsShop.reg la IP del servidor y la muestra en el html
		if accion == "send_ip" {
			var ip1, ip2, ip3, ip4, port int
			libs.LoadSettingsWin(serverRoot, settings)
			fmt.Sscanf(settings["serverinterno"], "http://%d.%d.%d.%d:%d", &ip1, &ip2, &ip3, &ip4, &port)
			output := fmt.Sprintf("%d;%d;%d;%d;%d", ip1, ip2, ip3, ip4, port)
			fmt.Fprintf(w, output)
		}
		//Recoge del html la direccion ip de la tienda y la modifica en SettingsShop.reg
		if accion == "edit_ip" {
			r.ParseForm()
			input, err := ioutil.ReadFile(serverRoot)
			if err != nil {
				Error.Println(err)
			}
			lines := strings.Split(string(input), "\r\n")
			for i, line := range lines {
				if strings.Contains(line, "serverinterno") {
					lines[i] = fmt.Sprintf("serverinterno = http://%s.%s.%s.%s:%s", r.FormValue("ip1"), r.FormValue("ip2"), r.FormValue("ip3"), r.FormValue("ip4"), r.FormValue("port"))
				}
			}
			output := strings.Join(lines, "\r\n")
			err = ioutil.WriteFile(serverRoot, []byte(output), 0644)
			if err != nil {
				Error.Println(err)
			}
			fmt.Fprint(w, "<div class='text-success'>La IP del servidor se ha modificado</div>")
		}
		//Muestra las horas en el fichero de programacion de música
		if accion == "mostrar_horas" {
			var h_ini, h_fin string
			var output string
			err := db.QueryRow("SELECT hora_inicial, hora_final FROM horario").Scan(&h_ini, &h_fin)
			if err != nil {
				Error.Println(err)
			}
			hora_ini := libs.MostrarHoras(h_ini)
			mins_ini := libs.MostrarMinutos(h_ini)
			hora_fin := libs.MostrarHoras(h_fin)
			mins_fin := libs.MostrarMinutos(h_fin)
			output = fmt.Sprintf("%s;%s;%s;%s", hora_ini, mins_ini, hora_fin, mins_fin)
			fmt.Fprint(w, output)
		}
		//Recoge la hora del formulario de programación de música
		if accion == "recoger_horas" {
			r.ParseForm()
			var h_ini, h_fin string
			var hora_inicial, hora_final string
			//Formamos la hora inicial y hora final
			hora_inicial = r.FormValue("hora1") + ":" + r.FormValue("min1")
			hora_final = r.FormValue("hora2") + ":" + r.FormValue("min2")
			db.QueryRow("SELECT hora_inicial, hora_final FROM horario").Scan(&h_ini, &h_fin)
			//Se comprueba las variable en base de datos para insertar o actualizar
			if h_ini == "" && h_fin == "" {
				//Guardamos en base de datos
				stm, err := db.Prepare("INSERT INTO horario (`hora_inicial`, `hora_final`) VALUES (?,?)")
				if err != nil {
					Error.Println(err)
				}
				db_mu.Lock()
				_, err1 := stm.Exec(hora_inicial, hora_final)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
				}
			} else {
				//Actualizamos el horario en base de datos
				ok, err := db.Prepare("UPDATE horario SET hora_inicial=?, hora_final=?")
				if err != nil {
					Error.Println(err)
				}
				db_mu.Lock()
				_, err1 := ok.Exec(hora_inicial, hora_final)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
				}
			}
		}
	}
}

//Reproductor de Mensajes Instantaneos
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
			if strings.Contains(val.Name(), ".mp3") {
				//Formamos el select
				output += fmt.Sprintf("<option value='%s'>%s</option>", val.Name(), val.Name())
				//Guardamos el nombre del primer mensaje
				if key == 0 {
					msg_instantaneo = val.Name()
				}
			}
		}
		//Mostrar msg: muestra el nombre del mensaje por primera vez
		output += fmt.Sprintf("#:#<span style='color: #006400'>Mensaje seleccionado: </span>" + msg_instantaneo)
	}
	//Estado de mensaje
	if r.FormValue("action") == "status" {
		msg_instantaneo = r.FormValue("instantaneos")
		output = fmt.Sprintf("<span style='color: #006400'>Mensaje seleccionado: </span>" + msg_instantaneo)
	}
	//Recibe el mensaje instantaneo y lo procesa
	if r.FormValue("action") == "send" {
		var win winamp.Winamp
		//Reproducimos el mensaje instantaneo
		win.PlayFFplay(msg_files_location + r.FormValue("instantaneos"))
	}
	fmt.Fprint(w, output)
}

//Programar Musica para la Tienda
func programarMusica(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var output string
	accion := r.FormValue("accion")
	//Muestra los directorios de musica de la tienda
	if accion == "show_dirs" {
		//Abrimos el directorio (C:\instore\Music\)
		file, err := os.Open(music_files)
		defer file.Close()
		if err != nil {
			Error.Println(err)
			return
		}
		musicDirs, err := file.Readdir(0)
		if err != nil {
			Error.Println(err)
			return
		}
		for _, val := range musicDirs {
			if !strings.Contains(val.Name(), ".txt") {
				output += fmt.Sprintf("<tr><td><input type='checkbox' name='musicDirs' value='%s'></td><td>&nbsp;</td><td>%s</td>", val.Name(), val.Name())
			}
		}
	}
	//Se recogen los datos de formulario (prog.html)
	if accion == "enviar" {
		cont := 0
		if len(programmedMusic) == 0 {
			for clave, valor := range r.Form {
				for _, v := range valor {
					if clave == "musicDirs" {
						db_mu.Lock()
						//Guardamos cada una de los directorios seleccionadas
						programmedMusic[cont] = v
						statusProgammedMusic = "Inicial"
						db_mu.Unlock()
					}
					cont++
				}
			}
		} else {
			//Borramos los antiguos directorios programados
			for k, _ := range programmedMusic {
				delete(programmedMusic, k)
			}
			//Añadimos los nuevos directorios
			for clave, valor := range r.Form {
				for _, v := range valor {
					if clave == "musicDirs" {
						db_mu.Lock()
						//Guardamos cada una de las carpetas seleccionadas
						programmedMusic[cont] = v
						if statusProgammedMusic == "Actualizada" {
							statusProgammedMusic = "Modificar"
						} else {
							statusProgammedMusic = "Actualizada"
						}
						db_mu.Unlock()
					}
					cont++
				}
			}
		}
	}
	fmt.Fprint(w, output)
}
