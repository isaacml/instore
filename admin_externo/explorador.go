package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/isaacml/instore/libs"
)

//Variables para guardar el identificador anterior, en caso de no encontrar datos.
var last_entidad, last_almacen, last_pais, last_region, last_prov, last_tienda string

//Función principal del explorador windows
func explorer(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var output string //variable para imprimir los datos hacia JavaScript

	//Muestra por primera vez las Unidades de Disco que tiene el Sistema
	if r.FormValue("action") == "unidades" {
		drives, err := exec.Command("cmd", "/c", "fsutil fsinfo drives").Output()
		if err != nil {
			Error.Println(err)
			return
		}
		output = "<option value='' selected>[Selecciona una unidad]</option>"
		res := strings.Split(string(drives), ": ")
		limpiar := strings.TrimSpace(string(libs.LimpiarMatriz([]byte(res[1]))))
		unidades := strings.Split(limpiar, "\\")
		for _, v := range unidades {
			v = strings.TrimSpace(v)
			if v != "" {
				output += fmt.Sprintf("<option value='%s'>%s</option>", v, v)
			}
		}
		fmt.Fprint(w, output)
	}

	//EXPLORADOR DE DIRECTORIOS --> FORMULARIO 1(testform)
	if r.FormValue("action") == "dir_unidad" {
		if r.FormValue("unidades") != "" {
			db_mu.Lock()
			//Mostramos los directorios de la unidad seleccionada
			directorio_actual = r.FormValue("unidades") + "\\"
			db_mu.Unlock()
			file, err := os.Open(directorio_actual)
			defer file.Close()
			if err != nil {
				Error.Println(err)
				return
			}
			directorios, err := file.Readdir(0)
			if err != nil {
				Error.Println(err)
				return
			}
			output = "<option value='' selected>[Selecciona un directorio]</option>"
			for _, val := range directorios {
				if val.IsDir() {
					output += fmt.Sprintf("<option value='%s'>%s</option>", val.Name(), val.Name())

				}
			}
			output += ";<span style='color: #B8860B'>" + directorio_actual + "</span>"
			fmt.Fprint(w, output)

		} else {
			db_mu.Lock()
			directorio_actual = ""
			db_mu.Unlock()
			output += ";<span style='color: #B8860B'>" + directorio_actual + "</span>"
			fmt.Fprint(w, output)
		}
	}
	//EXPLORADOR DE DIRECTORIOS --> FORMULARIO 1(testform)
	if r.FormValue("action") == "directorios" {
		if r.FormValue("directory") != "" && r.FormValue("directory") != "..." {
			db_mu.Lock()
			directorio_actual = directorio_actual + r.FormValue("directory") + "\\"
			db_mu.Unlock()
			file, err := os.Open(directorio_actual)
			defer file.Close()
			if err != nil {
				// No se puede abrir el directorio, por falta de permisos
				Error.Println(err)
				//Volvemos a tomar el archivo anterior y lo abrimos
				old := strings.Split(directorio_actual, r.FormValue("directory")+"\\")
				db_mu.Lock()
				directorio_actual = old[0]
				db_mu.Unlock()
				file2, err := os.Open(old[0])
				defer file.Close()
				directorios, err := file2.Readdir(0)
				if err != nil {
					Error.Println(err)
					return
				}
				output = "<option value='' selected>[Selecciona un directorio]</option></option><option value='...'>...</option>"
				for _, val := range directorios {
					if val.IsDir() {
						output += fmt.Sprintf("<option value='%s'>%s</option>", val.Name(), val.Name())

					}
				}
				output += ";<span style='color: #800000'>Necesitas permisos para abrir ese directorio</span>"
				fmt.Fprint(w, output)
				return
			}
			directorios, err := file.Readdir(0)
			if err != nil {
				Error.Println(err)
				return
			}
			output = fmt.Sprintf("<option value='' selected>[Selecciona un directorio]</option><option value='...'>...</option>")
			for _, val := range directorios {
				if val.IsDir() {
					output += fmt.Sprintf("<option value='%s'>%s</option>", val.Name(), val.Name())

				}
			}
			output += ";<span style='color: #B8860B'>" + directorio_actual + "</span>"
			fmt.Fprint(w, output)

			//VOLVER UN DIRECTORIO ATRÁS
		} else if r.FormValue("directory") != "" && r.FormValue("directory") == "..." {
			var contenedor string
			var contador int
			//Array para guardar la ruta sin valores nulos y para guardar la nueva ruta generada
			var arr_sin_vacios, nueva_ruta []string
			//Separamos el directorio actual por su contrabarra, lo cual nos va a generar un array de directorios
			ruta := strings.Split(directorio_actual, "\\")
			for k, v := range ruta {
				if v == "" {
					//Borramos el valor nulo, y volvemos a formar un nuevo array
					arr_sin_vacios = libs.RemoveIndex(ruta, k)
				}
			}
			contador = len(arr_sin_vacios) - 1
			if contador == 1 {
				//Borramos la ultima posicion del array
				nueva_ruta = arr_sin_vacios[:contador]
				for _, v := range nueva_ruta {
					contenedor += v + "\\"
				}
				db_mu.Lock()
				//Guardamos la ruta que nos genera
				directorio_actual = contenedor
				db_mu.Unlock()
				//Abrimos el directorio y mostramos sus carpetas
				file, err := os.Open(directorio_actual)
				defer file.Close()
				if err != nil {
					Error.Println(err)
					return
				}
				directorios, err := file.Readdir(0)
				if err != nil {
					Error.Println(err)
					return
				}
				output = "<option value='' selected>[Selecciona un directorio]</option>"
				for _, val := range directorios {
					if val.IsDir() {
						output += fmt.Sprintf("<option value='%s'>%s</option>", val.Name(), val.Name())

					}
				}
				output += ";<span style='color: #B8860B'>" + directorio_actual + "</span>"
				fmt.Fprint(w, output)
			} else {
				//Borramos la ultima posicion del array
				nueva_ruta = arr_sin_vacios[:contador]
				for _, v := range nueva_ruta {
					contenedor += v + "\\"
				}
				db_mu.Lock()
				//Guardamos la ruta que nos genera
				directorio_actual = contenedor
				db_mu.Unlock()
				//Abrimos el directorio y mostramos sus carpetas
				file, err := os.Open(directorio_actual)
				defer file.Close()
				if err != nil {
					Error.Println(err)
					return
				}
				directorios, err := file.Readdir(0)
				if err != nil {
					Error.Println(err)
					return
				}
				output = fmt.Sprintf("<option value='' selected>[Selecciona un directorio]</option><option value='...'>...</option>")
				for _, val := range directorios {
					if val.IsDir() {
						output += fmt.Sprintf("<option value='%s'>%s</option>", val.Name(), val.Name())

					}
				}
				output += ";<span style='color: #B8860B'>" + directorio_actual + "</span>"
				fmt.Fprint(w, output)
			}
		} else {
			//Abrimos el directorio y mostramos sus carpetas
			file, err := os.Open(directorio_actual)
			defer file.Close()
			if err != nil {
				Error.Println(err)
				return
			}
			directorios, err := file.Readdir(0)
			if err != nil {
				Error.Println(err)
				return
			}
			output = fmt.Sprintf("<option value='' selected>[Selecciona un directorio]</option><option value='...'>...</option>")
			for _, val := range directorios {
				if val.IsDir() {
					output += fmt.Sprintf("<option value='%s'>%s</option>", val.Name(), val.Name())

				}
			}
			output += ";<span style='color: #B8860B'>" + directorio_actual + "</span>"
			fmt.Fprint(w, output)
		}
	}
	//EXPLORADOR DE FICHEROS --> Envia los ficheros al FORMULARIO 2 (testform2)
	if r.FormValue("action") == "ficheros" {
		file, err := os.Open(directorio_actual)
		defer file.Close()
		if err != nil {
			return
		}
		ficheros, err := file.Readdir(0)
		if err != nil {
			Error.Println(err)
			return
		}
		for _, val := range ficheros {
			if strings.Contains(val.Name(), ".mp3") || strings.Contains(val.Name(), ".wma") {
				output += fmt.Sprintf("<option style='color: #0000CC' value='%s'>%s</option>", val.Name(), val.Name())
			}
		}
		fmt.Fprint(w, output)
	}
	//TOMA LOS FICHEROS DEL FORMULARIO 2, Y LOS PROCESA
	//r.FormValue("type") == "publi", procedemos a insertar los datos en la tabla publi
	if r.FormValue("action") == "get_ficheros" && r.FormValue("type") == "publi" {

		//Variables
		f_inicio := r.FormValue("f_inicio")
		f_final := r.FormValue("f_fin")
		if f_inicio == "" || f_final == "" {
			output += ";<span style='color: #FF0303'>Los campos fecha no pueden estar vacíos</span>"
			fmt.Fprint(w, output)
		} else {
			dest := estado_destino
			timestamp := time.Now().Unix()
			gap := r.FormValue("gap")
			//trozeamos las fechas
			arr_inicio := strings.Split(f_inicio, "/")
			arr_final := strings.Split(f_final, "/")
			//establecemos el formato de fechas para la BD --> yyyymmdd
			fecha_SQL_inc := fmt.Sprintf("%s%s%s", arr_inicio[2], arr_inicio[1], arr_inicio[0])
			fecha_SQL_fin := fmt.Sprintf("%s%s%s", arr_final[2], arr_final[1], arr_final[0])

			for clave, valor := range r.Form {
				for _, v := range valor {
					if clave == "files" {
						//Insertamos datos en la tabla interna del admininistrador (programaciones.sql)
						stmt, err0 := db.Prepare("INSERT INTO publi (`ruta`, `fichero`, `fecha_inicio`, `fecha_final`, `destino`, `timestamp`, `gap`) VALUES (?,?,?,?,?,?,?)")
						if err0 != nil {
							Error.Println(err0)
						}
						db_mu.Lock()
						_, err1 := stmt.Exec(directorio_actual, v, fecha_SQL_inc, fecha_SQL_fin, dest, timestamp, gap)
						db_mu.Unlock()
						if err1 != nil {
							Error.Println(err1)
							output += ";<span style='color: #FF0303'>Fallo al subir los ficheros</span>"
							fmt.Fprint(w, output)
						} else {
							output += ";<span style='color: #2E8B57'>Archivo/os subido/os correctamente</span>"
							fmt.Fprint(w, output)
						}
					}
				}
			}
		}
	}
	//r.FormValue("type") == "msg", procedemos a insertar los datos en la tabla mensajes
	if r.FormValue("action") == "get_ficheros" && r.FormValue("type") == "msg" {
		//Variables
		f_inicio := r.FormValue("f_inicio")
		f_final := r.FormValue("f_fin")

		if f_inicio == "" || f_final == "" {
			output += ";<span style='color: #FF0303'>Los campos fecha no pueden estar vacíos</span>"
			fmt.Fprint(w, output)
		} else {
			dest := estado_destino
			timestamp := time.Now().Unix()
			hora := r.FormValue("hora")
			min := r.FormValue("minutos")

			//trozeamos las fechas
			arr_inicio := strings.Split(f_inicio, "/")
			arr_final := strings.Split(f_final, "/")
			//establecemos el formato de fechas para la BD --> yyyymmdd
			fecha_SQL_inc := fmt.Sprintf("%s%s%s", arr_inicio[2], arr_inicio[1], arr_inicio[0])
			fecha_SQL_fin := fmt.Sprintf("%s%s%s", arr_final[2], arr_final[1], arr_final[0])
			//formamos el campo playtime
			playtime := hora + ":" + min

			for clave, valor := range r.Form {
				for _, v := range valor {
					if clave == "files" {
						//Insertamos datos en la tabla interna del admininistrador (programaciones.sql)
						stmt, err0 := db.Prepare("INSERT INTO mensaje (`ruta`, `fichero`, `fecha_inicio`, `fecha_final`, `destino`, `timestamp`, `playtime`) VALUES (?,?,?,?,?,?,?)")
						if err0 != nil {
							Error.Println(err0)
						}
						db_mu.Lock()
						_, err1 := stmt.Exec(directorio_actual, v, fecha_SQL_inc, fecha_SQL_fin, dest, timestamp, playtime)
						db_mu.Unlock()
						if err1 != nil {
							Error.Println(err1)
							output += ";<span style='color: #FF0303'>Fallo al subir los ficheros</span>"
							fmt.Fprint(w, output)
						} else {
							output += ";<span style='color: #2E8B57'>Archivo/os subido/os correctamente</span>"
							fmt.Fprint(w, output)
						}
					}
				}
			}
		}
	}
}

//Explorador de destino: se obtiene una ruta de destino final para el usuario.
//El usuario se encarga de generar la ruta eligiendo las distintas organizaciones.
func dest_explorer(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		var output string
		var arr_entidad []string
		updateExpires(sid) //Actualizamos el tiempo de expiración de la clave
		if r.FormValue("action") == "destinos" {
			//Enviamos nombre de usuario recogido en el formulario hacia el server para evaluar los destinos
			resultado := libs.GenerateFORM(settings["serverroot"]+"/acciones.cgi", "accion;destinos", "internal_action;gent_ent", "userAdmin;"+username)
			res := strings.Split(resultado, "@@")
			//Usuario ADMIN: puede ver todas las organizaciones que ha creado
			if res[0] == "ADMIN" {
				div_ent := strings.Split(res[1], "::")
				for _, val := range div_ent {
					if val != "" {
						arr_entidad = strings.Split(val, ";")
						output += fmt.Sprintf("<option value='entidad:.:%s'>%s</option>", arr_entidad[0], arr_entidad[1])
						db_mu.Lock()
						back_org = "entidad"
						db_mu.Unlock()
					}
				}
				db_mu.Lock()
				estado_destino = "*"
				db_mu.Unlock()
				output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #2E8B57'></span>"
			} else { //Solo la organizacion a la que pertenece
				resultado := libs.GenerateFORM(settings["serverroot"]+"/acciones.cgi", "accion;destinos", "internal_action;gent_ent_noadmin", "userAdmin;"+username)
				res := strings.Split(resultado, "@@")
				div_ent := strings.Split(res[1], "::")
				for _, val := range div_ent {
					if val != "" {
						arr_entidad = strings.Split(val, ";")
						output += fmt.Sprintf("<option value='entidad:.:%s'>%s</option>", arr_entidad[0], arr_entidad[1])
						db_mu.Lock()
						back_org = "entidad"
						db_mu.Unlock()
					}
				}
				db_mu.Lock()
				estado_destino = "*"
				db_mu.Unlock()
				output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #2E8B57'></span>"
			}
			fmt.Fprint(w, output)
		}
		// Recogemos los datos al hacer ONCLICK en formulario de destinos en publi.html
		if r.FormValue("action") == "recoger_id" {
			var destino, ident string
			if r.FormValue("destinos") == "" {
				destino = back_org
			} else {
				valores := strings.Split(r.FormValue("destinos"), ":.:")
				//Revisamos que el array tenga más de un valor, sino da un panic
				destino = valores[0]
				ident = valores[1]
			}
			if destino == "entidad" {
				var st_entidad string //variable que va a contener el estado de la entidad
				//Enviamos nombre de usuario e id_entidad recogido en el formulario hacia el server para generar los destinos
				resultado := libs.GenerateFORM(settings["serverroot"]+"/acciones.cgi", "accion;destinos", "internal_action;entidades", "userAdmin;"+username, "id_entidad;"+ident)
				if resultado != "" {
					output, st_entidad = libs.GenerateSelectOrg(resultado, "almacen")
					db_mu.Lock()
					//Se guarda el identificador, para poder volver atrás
					last_entidad = ident
					//Se forma el nuevo estado
					estado_destino = st_entidad + ".*"
					back_org = "almacen"
					db_mu.Unlock()
					output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #2E8B57'></span>"
					fmt.Fprint(w, output)
				} else {
					//Si no hay resultado, volvemos a cargar las entidades
					resultado2 := libs.GenerateFORM(settings["serverroot"]+"/acciones.cgi", "accion;destinos", "internal_action;gent_ent", "userAdmin;"+username)
					if resultado2 != "" {
						res := strings.Split(resultado2, "@@")
						arr := strings.Split(res[1], "::")
						for _, val := range arr {
							if val != "" {
								arr_entidad = strings.Split(val, ";")
								output += fmt.Sprintf("<option value='entidad:.:%s'>%s</option>", arr_entidad[0], arr_entidad[1])
							}
						}
						db_mu.Lock()
						estado_destino = "*"
						back_org = "entidad"
						db_mu.Unlock()
						output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #FF0303'>No hay sub-organizaciones</span>"
						fmt.Fprint(w, output)
					}
				}
			}
			if destino == "almacen" {
				var st_almacen string //variable que va a contener el estado del almacen
				if ident == "0" {
					//Si no hay resultado, volvemos a cargar las entidades
					resultado2 := libs.GenerateFORM(settings["serverroot"]+"/acciones.cgi", "accion;destinos", "internal_action;gent_ent", "userAdmin;"+username)
					if resultado2 != "" {
						res := strings.Split(resultado2, "@@")
						arr := strings.Split(res[1], "::")
						for _, val := range arr {
							if val != "" {
								arr_entidad = strings.Split(val, ";")
								output += fmt.Sprintf("<option value='entidad:.:%s'>%s</option>", arr_entidad[0], arr_entidad[1])
							}
						}
						db_mu.Lock()
						estado_destino = "*"
						back_org = "entidad"
						db_mu.Unlock()
						output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #2E8B57'></span>"
						fmt.Fprint(w, output)
					}
				} else {
					//Enviamos nombre de usuario e id_almacen recogido en el formulario hacia el server para generar los destinos
					resultado := libs.GenerateFORM(settings["serverroot"]+"/acciones.cgi", "accion;destinos", "internal_action;almacenes", "userAdmin;"+username, "id_almacen;"+ident)
					if resultado != "" {
						output, st_almacen = libs.GenerateSelectOrg(resultado, "pais")
						db_mu.Lock()
						//Se guarda el identificador, para poder volver atrás
						last_almacen = ident
						//Se borra el asterisco(*)
						res := libs.BackDestOrg(estado_destino, 1)
						//Se forma el nuevo estado
						estado_destino = res + st_almacen + ".*"
						back_org = "pais"
						db_mu.Unlock()
						output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #2E8B57'></span>"
						fmt.Fprint(w, output)
					} else {
						//Si no hay resultado, volvemos a cargar los almacenes
						resultado2 := libs.GenerateFORM(settings["serverroot"]+"/acciones.cgi", "accion;destinos", "internal_action;entidades", "userAdmin;"+username, "id_entidad;"+last_entidad)
						if resultado2 != "" {
							output, st_almacen = libs.GenerateSelectOrg(resultado2, "almacen")
							db_mu.Lock()
							//Se forma el nuevo estado
							estado_destino = st_almacen + ".*"
							back_org = "almacen"
							db_mu.Unlock()
							output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #FF0303'>No hay sub-organizaciones</span>"
							fmt.Fprint(w, output)
						}
					}
				}
			}
			if destino == "pais" {
				var st_pais string //variable que va a contener el estado del pais
				if ident == "0" {
					//Si no hay resultado, volvemos a cargar los almacenes
					resultado2 := libs.GenerateFORM(settings["serverroot"]+"/acciones.cgi", "accion;destinos", "internal_action;entidades", "userAdmin;"+username, "id_entidad;"+last_entidad)
					if resultado2 != "" {
						output, st_pais = libs.GenerateSelectOrg(resultado2, "almacen")
						db_mu.Lock()
						estado_destino = st_pais + ".*"
						back_org = "almacen"
						db_mu.Unlock()
						output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #2E8B57'></span>"
						fmt.Fprint(w, output)
					}
				} else {
					//Enviamos nombre de usuario e id_pais recogido en el formulario hacia el server para generar los destinos
					resultado := libs.GenerateFORM(settings["serverroot"]+"/acciones.cgi", "accion;destinos", "internal_action;paises", "userAdmin;"+username, "id_pais;"+ident)
					if resultado != "" {
						output, st_pais = libs.GenerateSelectOrg(resultado, "region")
						db_mu.Lock()
						//Se guarda el identificador, para poder volver atrás
						last_pais = ident
						//Se borra el asterisco(*)
						res := libs.BackDestOrg(estado_destino, 1)
						//Se forma el nuevo estado
						estado_destino = res + st_pais + ".*"
						back_org = "region"
						db_mu.Unlock()
						output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #2E8B57'></span>"
						fmt.Fprint(w, output)
					} else {
						//Si no hay resultado, volvemos a cargar los paises
						resultado2 := libs.GenerateFORM(settings["serverroot"]+"/acciones.cgi", "accion;destinos", "internal_action;almacenes", "userAdmin;"+username, "id_almacen;"+last_almacen)
						if resultado2 != "" {
							output, st_pais = libs.GenerateSelectOrg(resultado2, "pais")
							db_mu.Lock()
							//Se borra el asterisco(*) y se retrocede en una ORG.
							res := libs.BackDestOrg(estado_destino, 2)
							//Se forma el nuevo estado
							estado_destino = res + st_pais + ".*"
							back_org = "pais"
							db_mu.Unlock()
							output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #FF0303'>No hay sub-organizaciones</span>"
							fmt.Fprint(w, output)
						}
					}
				}
			}
			if destino == "region" {
				var st_region string //variable que va a contener el estado de la región
				if ident == "0" {
					//Si no hay resultado, volvemos a cargar los paises
					resultado2 := libs.GenerateFORM(settings["serverroot"]+"/acciones.cgi", "accion;destinos", "internal_action;almacenes", "userAdmin;"+username, "id_almacen;"+last_almacen)
					if resultado2 != "" {
						output, st_region = libs.GenerateSelectOrg(resultado2, "pais")
						db_mu.Lock()
						res := libs.BackDestOrg(estado_destino, 3)
						//Se forma el nuevo estado
						estado_destino = res + st_region + ".*"
						back_org = "pais"
						db_mu.Unlock()
						output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #2E8B57'></span>"
						fmt.Fprint(w, output)
					}
				} else {
					//Enviamos nombre de usuario e id_region recogido en el formulario hacia el server para generar los destinos
					resultado := libs.GenerateFORM(settings["serverroot"]+"/acciones.cgi", "accion;destinos", "internal_action;regiones", "userAdmin;"+username, "id_region;"+ident)
					if resultado != "" {
						output, st_region = libs.GenerateSelectOrg(resultado, "provincia")
						db_mu.Lock()
						//Se guarda el identificador, para poder volver atrás
						last_region = ident
						//Se borra el asterisco(*)
						res := libs.BackDestOrg(estado_destino, 1)
						estado_destino = res + st_region + ".*"
						back_org = "provincia"
						db_mu.Unlock()
						output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #2E8B57'></span>"
						fmt.Fprint(w, output)
					} else {
						//Si no hay resultado, volvemos a cargar las regiones
						resultado2 := libs.GenerateFORM(settings["serverroot"]+"/acciones.cgi", "accion;destinos", "internal_action;paises", "userAdmin;"+username, "id_pais;"+last_pais)
						if resultado2 != "" {
							output, st_region = libs.GenerateSelectOrg(resultado2, "region")
							db_mu.Lock()
							//Se borra el asterisco(*) y se retrocede en una ORG.
							res := libs.BackDestOrg(estado_destino, 2)
							estado_destino = res + st_region + ".*"
							back_org = "region"
							db_mu.Unlock()
							output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #FF0303'>No hay sub-organizaciones</span>"
							fmt.Fprint(w, output)
						}
					}
				}
			}
			if destino == "provincia" {
				var st_provincia string //variable que va a contener el estado de la provincia
				if ident == "0" {
					//Si no hay resultado, volvemos a cargar las regiones
					resultado2 := libs.GenerateFORM(settings["serverroot"]+"/acciones.cgi", "accion;destinos", "internal_action;paises", "userAdmin;"+username, "id_pais;"+last_pais)
					if resultado2 != "" {
						output, st_provincia = libs.GenerateSelectOrg(resultado2, "region")
						db_mu.Lock()
						res := libs.BackDestOrg(estado_destino, 3)
						//Se forma el nuevo estado
						estado_destino = res + st_provincia + ".*"
						back_org = "region"
						db_mu.Unlock()
						output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #2E8B57'></span>"
						fmt.Fprint(w, output)
					}
				} else {
					//Enviamos nombre de usuario e id_provincia recogido en el formulario hacia el server para generar los destinos
					resultado := libs.GenerateFORM(settings["serverroot"]+"/acciones.cgi", "accion;destinos", "internal_action;provincias", "userAdmin;"+username, "id_provincia;"+ident)
					if resultado != "" {
						output, st_provincia = libs.GenerateSelectOrg(resultado, "tienda")
						db_mu.Lock()
						//Se guarda el identificador, para poder volver atrás
						last_prov = ident
						//Se borra el asterisco(*)
						res := libs.BackDestOrg(estado_destino, 1)
						//Se forma el nuevo estado
						estado_destino = res + st_provincia + ".*"
						back_org = "tienda"
						db_mu.Unlock()
						output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #2E8B57'></span>"
						fmt.Fprint(w, output)
					} else {
						//Si no hay resultado, volvemos a cargar las regiones
						resultado2 := libs.GenerateFORM(settings["serverroot"]+"/acciones.cgi", "accion;destinos", "internal_action;regiones", "userAdmin;"+username, "id_region;"+last_region)
						if resultado2 != "" {
							output, st_provincia = libs.GenerateSelectOrg(resultado2, "provincia")
							db_mu.Lock()
							//Se borra el asterisco(*) y se retrocede en una ORG.
							res := libs.BackDestOrg(estado_destino, 2)
							//Se forma el nuevo estado
							estado_destino = res + st_provincia + ".*"
							back_org = "provincia"
							db_mu.Unlock()
							output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #FF0303'>No hay sub-organizaciones</span>"
							fmt.Fprint(w, output)
						}
					}
				}
			}
			if destino == "tienda" {
				var st_tienda string //variable que va a contener el estado de la tienda
				if ident == "0" {
					//Si no hay resultado, volvemos a cargar las regiones
					resultado2 := libs.GenerateFORM(settings["serverroot"]+"/acciones.cgi", "accion;destinos", "internal_action;regiones", "userAdmin;"+username, "id_region;"+last_region)
					if resultado2 != "" {
						output, st_tienda = libs.GenerateSelectOrg(resultado2, "provincia")
						db_mu.Lock()
						res := libs.BackDestOrg(estado_destino, 3)
						//Se forma el nuevo estado
						estado_destino = res + st_tienda + ".*"
						back_org = "provincia"
						db_mu.Unlock()
						output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #2E8B57'></span>"
						fmt.Fprint(w, output)
					}
				} else {
					//Enviamos nombre de usuario e id_provincia recogido en el formulario hacia el server para generar los destinos
					resultado := libs.GenerateFORM(settings["serverroot"]+"/acciones.cgi", "accion;destinos", "internal_action;tiendas", "userAdmin;"+username, "id_tienda;"+ident)
					if resultado != "" {
						output = "<option value='destino_final:.:0'>...</option>"
						db_mu.Lock()
						//Se guarda el identificador, para poder volver atrás
						last_tienda = ident
						//Se borra el asterisco(*)
						res := libs.BackDestOrg(estado_destino, 1)
						estado_destino = res + resultado
						back_org = "destino_final"
						db_mu.Unlock()
						output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #2E8B57'></span>"
						fmt.Fprint(w, output)
					} else {
						//Si no hay resultado, volvemos a cargar las provincias
						var arr_tienda []string
						resultado2 := libs.GenerateFORM(settings["serverroot"]+"/acciones.cgi", "accion;destinos", "internal_action;provincias", "userAdmin;"+username, "id_provincia;"+last_prov)
						if resultado2 != "" {
							output = "<option value='tienda:.:0'>...</option>"
							arr := strings.Split(resultado2, "::")
							for _, val := range arr {
								if val != "" {
									arr_tienda = strings.Split(val, ";")
									output += fmt.Sprintf("<option value='tienda:.:%s'>%s</option>", arr_tienda[0], arr_tienda[1])
								}
							}
							db_mu.Lock()
							res := libs.BackDestOrg(estado_destino, 1)
							estado_destino = res + "*"
							back_org = "tienda"
							db_mu.Unlock()
							output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #FF0303'>No hay sub-organizaciones</span>"
							fmt.Fprint(w, output)
						}
					}
				}
			}
			if destino == "destino_final" {
				if ident == "0" {
					//value = 0 : volvemos a cargar las tiendas
					var arr_tienda []string
					resultado2 := libs.GenerateFORM(settings["serverroot"]+"/acciones.cgi", "accion;destinos", "internal_action;provincias", "userAdmin;"+username, "id_provincia;"+last_prov)
					if resultado2 != "" {
						output = "<option value='tienda:.:0'>...</option>"
						arr := strings.Split(resultado2, "::")
						for _, val := range arr {
							if val != "" {
								arr_tienda = strings.Split(val, ";")
								output += fmt.Sprintf("<option value='tienda:.:%s'>%s</option>", arr_tienda[0], arr_tienda[1])
							}
						}
						db_mu.Lock()
						res := libs.BackDestOrg(estado_destino, 1)
						estado_destino = res + "*"
						back_org = "tienda"
						db_mu.Unlock()
						output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #2E8B57'></span>"
						fmt.Fprint(w, output)
					}
				}
				//El destino ya se ha generado, solo hay opcion de volver atrás
				if ident == "" {
					db_mu.Lock()
					back_org = "destino_final"
					db_mu.Unlock()
					output = "<option title='Volver Atrás' value='destino_final:.:0'>...</option>;<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #FF0303'>El destino ya está formado</span>"
					fmt.Fprint(w, output)
				}
			}
		}
	}
}

//Vista: acciones sobre la publicidad y los mensajes
func vista(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		accion := r.FormValue("accion")
		if accion == "first_show" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(settings["serverroot"]+"/modo_vista.cgi", "accion;"+accion, "tabla;"+r.FormValue("tabla"), "username;"+username))
			fmt.Fprint(w, respuesta)
		}
		if accion == "mostrar" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(settings["serverroot"]+"/modo_vista.cgi", "accion;"+accion, "tabla;"+r.FormValue("tabla"), "search;"+r.FormValue("buscar_por"), "username;"+username))
			fmt.Fprint(w, respuesta)
		}
		if accion == "borrar" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(settings["serverroot"]+"/modo_vista.cgi", "accion;"+accion, "tabla;"+r.FormValue("tabla"), "borrar;"+r.FormValue("borrar")))
			fmt.Fprint(w, respuesta)
		}
		if accion == "load" {
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(settings["serverroot"]+"/modo_vista.cgi", "accion;"+accion, "tabla;"+r.FormValue("tabla"), "edit_id;"+r.FormValue("load")))
			fmt.Fprint(w, respuesta)
		}
		if accion == "modificar" {
			var recojo_destino string
			if r.FormValue("tomar_dest") == "SI" {
				recojo_destino = estado_destino
			}
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(settings["serverroot"]+"/modo_vista.cgi", "accion;"+accion,
				"tabla;"+r.FormValue("tabla"), "id;"+r.FormValue("id"), "f_ini;"+r.FormValue("f_inicio"), "f_fin;"+r.FormValue("f_fin"),
				"destino;"+recojo_destino, "gap;"+r.FormValue("gap"), "origen;"+r.FormValue("origen"), "hora;"+r.FormValue("hora"), "minutos;"+r.FormValue("minutos")))
			fmt.Fprint(w, respuesta)
		}
	}
}
