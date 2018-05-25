package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/isaacml/instore/libs"
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
			libs.LoadSettingsLin(serverRoot, settings)
			fmt.Fprintf(w, settings["serverinterno"])
		}
		//Recoge del html la direccion ip de la tienda y la modifica en SettingsShop.reg
		if accion == "edit_ip" {
			r.ParseForm()
			input, err := ioutil.ReadFile(serverRoot)
			if err != nil {
				Error.Println(err)
			}
			lines := strings.Split(string(input), "\n")
			for i, line := range lines {
				if strings.Contains(line, "serverinterno") {
					lines[i] = fmt.Sprintf("serverinterno = %s", r.FormValue("ip"))
				}
			}
			output := strings.Join(lines, "\n")
			err = ioutil.WriteFile(serverRoot, []byte(output), 0755)
			if err != nil {
				Error.Println(err)
			}
			fmt.Fprint(w, "<div class='text-success'>La dirección del servidor se ha modificado</div>")
		}
		//Muestra las horas en el fichero de programacion de música
		if accion == "mostrar_horas" {
			var h_ini, h_fin string
			var output string
			db.QueryRow("SELECT hora_inicial, hora_final FROM horario").Scan(&h_ini, &h_fin)
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
			var err error
			//Formamos la hora inicial y hora final (para usuario)
			hora_inicial := r.FormValue("hora1") + ":" + r.FormValue("min1")
			hora_final := r.FormValue("hora2") + ":" + r.FormValue("min2")
			db.QueryRow("SELECT hora_inicial, hora_final FROM horario").Scan(&h_ini, &h_fin)
			//Se comprueba las variable en base de datos para insertar o actualizar
			if h_ini != "" && h_fin != "" {
				//Borramos los datos que habían anteriormente
				db_mu.Lock()
				_, err = db.Exec("DELETE FROM horario")
				db_mu.Unlock()
				if err != nil {
					Error.Println(err)
				}
			}
			//Guardamos en base de datos horario
			stm, err1 := db.Prepare("INSERT INTO horario (`hora_inicial`, `hora_final`) VALUES (?,?)")
			if err1 != nil {
				Error.Println(err1)
			}
			db_mu.Lock()
			_, err1 = stm.Exec(hora_inicial, hora_final)
			db_mu.Unlock()
			if err1 != nil {
				Error.Println(err1)
			}
			//Hora inicial y final en entero (para nuestro programa)
			h_ini_int := libs.Hour2min(libs.ToInt(r.FormValue("hora1")), libs.ToInt(r.FormValue("min1")))
			h_fin_int := libs.Hour2min(libs.ToInt(r.FormValue("hora2")), libs.ToInt(r.FormValue("min2")))
			//Se comprueba la existencia de datos en la tabla
			db.QueryRow("SELECT hora_inicial, hora_final FROM aux").Scan(&h_ini, &h_fin)
			//Comprobamos si la hora inicial es mayor que la final
			if h_ini_int > h_fin_int {
				//Se comprueba la existencia de datos en la tabla
				if h_ini != "" && h_fin != "" {
					//Borramos los datos que habían anteriormente
					db_mu.Lock()
					_, err = db.Exec("DELETE FROM aux")
					db_mu.Unlock()
					if err != nil {
						Error.Println(err)
					}
				}
				db_mu.Lock()
				_, err = db.Exec("INSERT INTO aux (`hora_inicial`, `hora_final`) VALUES (?,?)", h_ini_int, 1439)
				db_mu.Unlock()
				if err != nil {
					Error.Println(err)
				}
				db_mu.Lock()
				_, err = db.Exec("INSERT INTO aux (`hora_inicial`, `hora_final`) VALUES (?,?)", 0, h_fin_int)
				db_mu.Unlock()
				if err != nil {
					Error.Println(err)
				}
			} else {
				//Se comprueba la existencia de datos en la tabla
				if h_ini != "" && h_fin != "" {
					//Borramos los datos que habían anteriormente
					db_mu.Lock()
					_, err = db.Exec("DELETE FROM aux")
					db_mu.Unlock()
					if err != nil {
						Error.Println(err)
					}
				}
				db_mu.Lock()
				//Metemos el dato nuevo
				_, err = db.Exec("INSERT INTO aux (`hora_inicial`, `hora_final`) VALUES (?,?)", h_ini_int, h_fin_int)
				db_mu.Unlock()
				if err != nil {
					Error.Println(err)
				}
			}
		}
	}
}

//Generar un listado de Mensajes Instantaneos
func instantaneos(w http.ResponseWriter, r *http.Request) {
	var output string //variable para imprimir los datos hacia JavaScript
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
	for _, val := range ficheros {
		//Tomamos solamente ficheros MP3
		if strings.Contains(val.Name(), ".mp3") || strings.Contains(val.Name(), ".wma") {
			//Formamos el select
			output += fmt.Sprintf("<option value='%s'>%s</option>", val.Name(), val.Name())
		}
	}
	fmt.Fprint(w, output)
}

func playInstantaneos(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	mensaje := r.FormValue("instantaneos")
	go winplayer.PlayFFplay(msg_files_location + mensaje)
}

//Programar Musica para la Tienda
func programarMusica(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	accion := r.FormValue("accion")
	//Muestra los directorios de musica de la tienda
	if accion == "show_dirs" {
		var numcarp int
		var salida string
		err := db.QueryRow("SELECT count(carpeta) FROM musica").Scan(&numcarp)
		if err != nil {
			Error.Println(err)
		}
		if numcarp == 0 {
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
					salida += fmt.Sprintf("<tr><td><input type='checkbox' name='musicDirs' value='%s'></td><td>&nbsp;</td><td>%s</td>", val.Name(), val.Name())
				}
			}
			fmt.Fprint(w, salida)
		} else {
			var salida string
			var cadena_salida string
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
			query, err := db.Query("SELECT carpeta FROM musica")
			if err != nil {
				Error.Println(err)
			}
			for query.Next() {
				var carpeta string
				err = query.Scan(&carpeta)
				if err != nil {
					Error.Println(err)
				}
				cadena_salida += carpeta + ";"
			}
			for _, val := range musicDirs {
				if !strings.Contains(val.Name(), ".txt") {
					if strings.Contains(cadena_salida, val.Name()) {
						salida += fmt.Sprintf("<tr><td><input type='checkbox' name='musicDirs' value='%s' checked></td><td>&nbsp;</td><td>%s</td>", val.Name(), val.Name())
					} else {
						salida += fmt.Sprintf("<tr><td><input type='checkbox' name='musicDirs' value='%s'></td><td>&nbsp;</td><td>%s</td>", val.Name(), val.Name())
					}
				}
			}
			fmt.Fprint(w, salida)
		}
	}
	//Se recogen los datos de formulario (prog.html)
	if accion == "enviar" {
		var num_carp int
		var err error
		err = db.QueryRow("SELECT count(carpeta) FROM musica").Scan(&num_carp)
		if err != nil {
			Error.Println(err)
		}
		if num_carp == 0 {
			for clave, valor := range r.Form {
				for _, v := range valor {
					if clave == "musicDirs" {
						db_mu.Lock()
						//Guardamos cada una de los directorios seleccionadas en la BD
						_, err = db.Exec("INSERT INTO musica (`carpeta`) VALUES (?)", v)
						db_mu.Unlock()
						if err != nil {
							Error.Println(err)
						}
					}
				}
			}
			db_mu.Lock()
			_, err = db.Exec("INSERT INTO st_prog_music (`estado`) VALUES (?)", "PrimerCambio")
			db_mu.Unlock()
			if err != nil {
				Error.Println(err)
			}
		} else {
			//Borramos las carpetas que habían anteriormente
			db_mu.Lock()
			_, err = db.Exec("DELETE FROM musica")
			db_mu.Unlock()
			if err != nil {
				Error.Println(err)
			} else {
				for clave, valor := range r.Form {
					for _, v := range valor {
						if clave == "musicDirs" {
							db_mu.Lock()
							//Guardamos cada una de los directorios seleccionadas en la BD
							_, err = db.Exec("INSERT INTO musica (`carpeta`) VALUES (?)", v)
							db_mu.Unlock()
							if err != nil {
								Error.Println(err)
							}
						}
					}
				}
				var st_prog string
				err = db.QueryRow("SELECT estado FROM st_prog_music").Scan(&st_prog)
				if err != nil {
					Error.Println(err)
				}
				if st_prog == "PrimerCambio" {
					db_mu.Lock()
					_, err = db.Exec("UPDATE st_prog_music SET estado=?", "SegundoCambio")
					db_mu.Unlock()
					if err != nil {
						Error.Println(err)
					}
				} else if st_prog == "SegundoCambio" {
					db_mu.Lock()
					_, err := db.Exec("UPDATE st_prog_music SET estado=?", "PrimerCambio")
					db_mu.Unlock()
					if err != nil {
						Error.Println(err)
					}
				}
			}
		}
	}
}

//Muestra o NO el boton de reproduccion de mensajes instantaneos
func mostrar_boton(w http.ResponseWriter, r *http.Request) {
	var output string
	mostrar := winplayer.Status().MostrarButton
	if mostrar == false {
		output = "Deshabilitamos"
	} else {
		output = "Habilitamos"
	}
	fmt.Fprint(w, output)
}

//Establece el volumen que va a tener el winamp
func volumen_global(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var salida string
	accion := r.FormValue("accion")
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		//Muestra los directorios de musica de la tienda
		if accion == "mostrar_volumen" {
			salida = fmt.Sprintf("<tr><td><input type='text' name='volumen' value='%d'></td></tr>", winplayer.Status().Volume)
		}
		if accion == "guardar_volumen" {
			vol := libs.ToInt(r.FormValue("volumen"))
			winplayer.SetVolume(vol)
		}
		fmt.Fprint(w, salida)
	}
}
