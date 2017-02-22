package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

//Variable que guarda el estado del destino
var estado_destino string

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
		output = "<option value='' selected>Selecciona una unidad</option>"
		res := strings.Split(string(drives), ": ")
		limpiar := strings.TrimSpace(string(limpiar_matriz([]byte(res[1]))))
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
			//Mostramos los directorios de la unidad seleccionada
			directorio_actual = r.FormValue("unidades") + "\\"
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
			for _, val := range directorios {
				if val.IsDir() {
					output += fmt.Sprintf("<option style='color: #B8860B' value='%s'>%s</option>", val.Name(), val.Name())

				}
			}
			output += ";<span style='color: #B8860B'>" + directorio_actual + "</span>"
			fmt.Fprint(w, output)
		} else {
			directorio_actual = ""
			output += ";<span style='color: #B8860B'>" + directorio_actual + "</span>"
			fmt.Fprint(w, output)
		}
	}
	//EXPLORADOR DE DIRECTORIOS --> FORMULARIO 1(testform)
	if r.FormValue("action") == "directorios" {
		if r.FormValue("directory") != "" && r.FormValue("directory") != "..." {
			directorio_actual = directorio_actual + r.FormValue("directory") + "\\"
			file, err := os.Open(directorio_actual)
			defer file.Close()
			if err != nil {
				// No se puede abrir el directorio, por falta de permisos
				Error.Println(err)
				//Volvemos a tomar el archivo anterior y lo abrimos
				old := strings.Split(directorio_actual, r.FormValue("directory")+"\\")
				directorio_actual = old[0]
				file2, err := os.Open(old[0])
				defer file.Close()
				directorios, err := file2.Readdir(0)
				if err != nil {
					Error.Println(err)
					return
				}
				for _, val := range directorios {
					if val.IsDir() {
						output += fmt.Sprintf("<option style='color: #B8860B' value='%s'>%s</option>", val.Name(), val.Name())

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
			output = fmt.Sprintf("<option style='color: #B8860B' value='...'>...</option>")
			for _, val := range directorios {
				if val.IsDir() {
					output += fmt.Sprintf("<option style='color: #B8860B' value='%s'>%s</option>", val.Name(), val.Name())

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
					arr_sin_vacios = RemoveIndex(ruta, k)
				}
			}
			contador = len(arr_sin_vacios) - 1
			if contador == 1 {
				//Borramos la ultima posicion del array
				nueva_ruta = arr_sin_vacios[:contador]
				for _, v := range nueva_ruta {
					contenedor += v + "\\"
				}
				//Guardamos la ruta que nos genera
				directorio_actual = contenedor
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
				for _, val := range directorios {
					if val.IsDir() {
						output += fmt.Sprintf("<option style='color: #B8860B' value='%s'>%s</option>", val.Name(), val.Name())

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
				//Guardamos la ruta que nos genera
				directorio_actual = contenedor
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
				output = fmt.Sprintf("<option style='color: #B8860B' value='...'>...</option>")
				for _, val := range directorios {
					if val.IsDir() {
						output += fmt.Sprintf("<option style='color: #B8860B' value='%s'>%s</option>", val.Name(), val.Name())

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
			output = fmt.Sprintf("<option style='color: #B8860B' value='...'>...</option>")
			for _, val := range directorios {
				if val.IsDir() {
					output += fmt.Sprintf("<option style='color: #B8860B' value='%s'>%s</option>", val.Name(), val.Name())

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
			if !val.IsDir() {
				output += fmt.Sprintf("<option style='color: #0000CC' value='%s'>%s</option>", val.Name(), val.Name())
			}
		}
		fmt.Fprint(w, output)
	}
	//TOMA LOS FICHEROS DEL FORMULARIO 2, Y LOS PROCESA
	if r.FormValue("action") == "get_ficheros" {
		//Variables
		f_inicio := r.FormValue("f_inicio")
		f_final := r.FormValue("f_fin")
		dest := r.FormValue("destino")
		timestamp := time.Now().Unix()
		//trozeamos las fechas
		arr_inicio := strings.Split(f_inicio, "/")
		arr_final := strings.Split(f_final, "/")
		//establecemos el formato de fechas para la BD --> yyyymmdd
		fecha_SQL_inc := fmt.Sprintf("%s%s%s", arr_inicio[2], arr_inicio[1], arr_inicio[0])
		fecha_SQL_fin := fmt.Sprintf("%s%s%s", arr_final[2], arr_final[1], arr_final[0])

		for clave, valor := range r.Form {
			for _, v := range valor {
				if clave == "files" {
					stmt, err0 := db.Prepare("INSERT INTO publi (`ruta`, `fichero`, `fecha_inicio`, `fecha_final`, `destino`, `timestamp`) VALUES (?,?,?,?,?,?)")
					if err0 != nil {
						Error.Println(err0)
					}
					db_mu.Lock()
					_, err1 := stmt.Exec(directorio_actual, v, fecha_SQL_inc, fecha_SQL_fin, dest, timestamp)
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

//Funcion que muestra los datos en los respectivos selects de destino
func recoger_destinos(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		var output string
		loadSettings(serverRoot)
		updateExpires(sid) //Actualizamos el tiempo de expiración de la clave
		if r.FormValue("action") == "destinos" {
			var arr_entidad []string
			//Enviamos nombre de usuario recogido en el formulario hacia el server para generar los destinos
			resultado := libs.GenerateFORM(serverext["serverroot"]+"/destino.cgi", "action;destinos", "userAdmin;"+username)
			arr := strings.Split(resultado, "::")
			for _, val := range arr {
				if val != "" {
					arr_entidad = strings.Split(val, ";")
					output += fmt.Sprintf("<option value='entidad:.:%s'>%s</option>", arr_entidad[0], arr_entidad[1])
				}
			}
			estado_destino = "*"
			output += ";<span style='color: #1A5276'>" + estado_destino + "</span>"
			fmt.Fprint(w, output)
		}
		// Recogemos los datos al hacer ONCLICK en formulario de destinos en publi.html
		if r.FormValue("action") == "recoger_id" {
			valores := strings.Split(r.FormValue("destinos"), ":.:")
			destino := valores[0]
			ident := valores[1]

			if destino == "entidad" {
				//Enviamos nombre de usuario e id_entidad recogido en el formulario hacia el server para generar los destinos
				resultado := libs.GenerateFORM(serverext["serverroot"]+"/destino.cgi", "action;entidades", "userAdmin;"+username, "id_entidad;"+ident)
				if resultado != "" {
					var arr_almacen []string
					output = "<option value='almacen:.:0'>...</option>"
					arr := strings.Split(resultado, "::")
					for _, val := range arr {
						if val != "" {
							arr_almacen = strings.Split(val, ";")
							output += fmt.Sprintf("<option value='almacen:.:%s'>%s</option>", arr_almacen[0], arr_almacen[1])
						}
					}
					last_entidad = ident
					estado_destino = arr_almacen[2] + ".*"
					output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #2E8B57'></span>"
					fmt.Fprint(w, output)
				} else {
					//Si no hay resultado, volvemos a cargar las entidades
					var arr_entidad []string
					resultado2 := libs.GenerateFORM(serverext["serverroot"]+"/destino.cgi", "action;destinos", "userAdmin;"+username)
					if resultado2 != "" {
						arr := strings.Split(resultado2, "::")
						for _, val := range arr {
							if val != "" {
								arr_entidad = strings.Split(val, ";")
								output += fmt.Sprintf("<option value='entidad:.:%s'>%s</option>", arr_entidad[0], arr_entidad[1])
							}
						}
						estado_destino = "*"
						output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #FF0303'>No hay sub-organizaciones</span>"
						fmt.Fprint(w, output)
					}
				}
			}
			if destino == "almacen" {
				if ident == "0" {
					//Si no hay resultado, volvemos a cargar las entidades
					var arr_entidad []string
					resultado2 := libs.GenerateFORM(serverext["serverroot"]+"/destino.cgi", "action;destinos", "userAdmin;"+username)
					if resultado2 != "" {
						arr := strings.Split(resultado2, "::")
						for _, val := range arr {
							if val != "" {
								arr_entidad = strings.Split(val, ";")
								output += fmt.Sprintf("<option value='entidad:.:%s'>%s</option>", arr_entidad[0], arr_entidad[1])
							}
						}
						estado_destino = "*"
						output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #2E8B57'></span>"
						fmt.Fprint(w, output)
					}
				} else {
					//Enviamos nombre de usuario e id_almacen recogido en el formulario hacia el server para generar los destinos
					resultado := libs.GenerateFORM(serverext["serverroot"]+"/destino.cgi", "action;almacenes", "userAdmin;"+username, "id_almacen;"+ident)
					if resultado != "" {
						var arr_pais []string
						var res string
						output = "<option value='pais:.:0'>...</option>"
						arr := strings.Split(resultado, "::")
						for _, val := range arr {
							if val != "" {
								arr_pais = strings.Split(val, ";")
								output += fmt.Sprintf("<option value='pais:.:%s'>%s</option>", arr_pais[0], arr_pais[1])
							}
						}
						last_almacen = ident
						separator := strings.Split(estado_destino, ".")
						arr_alm := separator[:len(separator)-1]
						for _, v := range arr_alm {
							res += v + "."
						}
						estado_destino = res + arr_pais[2] + ".*"
						output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #2E8B57'></span>"
						fmt.Fprint(w, output)
					} else {
						//Si no hay resultado, volvemos a cargar los almacenes
						var arr_almacen []string
						resultado2 := libs.GenerateFORM(serverext["serverroot"]+"/destino.cgi", "action;entidades", "userAdmin;"+username, "id_entidad;"+last_entidad)
						if resultado2 != "" {
							output = "<option value='almacen:.:0'>...</option>"
							arr := strings.Split(resultado2, "::")
							for _, val := range arr {
								if val != "" {
									arr_almacen = strings.Split(val, ";")
									output += fmt.Sprintf("<option value='almacen:.:%s'>%s</option>", arr_almacen[0], arr_almacen[1])
								}
							}
							estado_destino = arr_almacen[2] + ".*"
							output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #FF0303'>No hay sub-organizaciones</span>"
							fmt.Fprint(w, output)
						}
					}
				}
			}
			if destino == "pais" {
				if ident == "0" {
					//Si no hay resultado, volvemos a cargar los almacenes
					var arr_almacen []string
					resultado2 := libs.GenerateFORM(serverext["serverroot"]+"/destino.cgi", "action;entidades", "userAdmin;"+username, "id_entidad;"+last_entidad)
					if resultado2 != "" {
						output = "<option value='almacen:.:0'>...</option>"
						arr := strings.Split(resultado2, "::")
						for _, val := range arr {
							if val != "" {
								arr_almacen = strings.Split(val, ";")
								output += fmt.Sprintf("<option value='almacen:.:%s'>%s</option>", arr_almacen[0], arr_almacen[1])
							}
						}
						estado_destino = arr_almacen[2] + ".*"
						output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #2E8B57'></span>"
						fmt.Fprint(w, output)
					}
				} else {
					//Enviamos nombre de usuario e id_pais recogido en el formulario hacia el server para generar los destinos
					resultado := libs.GenerateFORM(serverext["serverroot"]+"/destino.cgi", "action;paises", "userAdmin;"+username, "id_pais;"+ident)
					if resultado != "" {
						var arr_region []string
						var res string
						output = "<option value='region:.:0'>...</option>"
						arr := strings.Split(resultado, "::")
						for _, val := range arr {
							if val != "" {
								arr_region = strings.Split(val, ";")
								output += fmt.Sprintf("<option value='region:.:%s'>%s</option>", arr_region[0], arr_region[1])
							}
						}
						last_pais = ident
						separator := strings.Split(estado_destino, ".")
						arr_pais := separator[:len(separator)-1]
						for _, v := range arr_pais {
							res += v + "."
						}
						estado_destino = res + arr_region[2] + ".*"
						output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #2E8B57'></span>"
						fmt.Fprint(w, output)
					} else {
						//Si no hay resultado, volvemos a cargar los paises
						var arr_pais []string
						var res string
						resultado2 := libs.GenerateFORM(serverext["serverroot"]+"/destino.cgi", "action;almacenes", "userAdmin;"+username, "id_almacen;"+last_almacen)
						if resultado2 != "" {
							output = "<option value='pais:.:0'>...</option>"
							arr := strings.Split(resultado2, "::")
							for _, val := range arr {
								if val != "" {
									arr_pais = strings.Split(val, ";")
									output += fmt.Sprintf("<option value='pais:.:%s'>%s</option>", arr_pais[0], arr_pais[1])
								}
							}
							separator := strings.Split(estado_destino, ".")
							arr_alm := separator[:len(separator)-2]
							for _, v := range arr_alm {
								res += v + "."
							}
							estado_destino = res + arr_pais[2] + ".*"
							output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #FF0303'>No hay sub-organizaciones</span>"
							fmt.Fprint(w, output)
						}
					}
				}
			}
			if destino == "region" {
				if ident == "0" {
					//Si no hay resultado, volvemos a cargar los paises
					var arr_pais []string
					var res string
					resultado2 := libs.GenerateFORM(serverext["serverroot"]+"/destino.cgi", "action;almacenes", "userAdmin;"+username, "id_almacen;"+last_almacen)
					if resultado2 != "" {
						output = "<option value='pais:.:0'>...</option>"
						arr := strings.Split(resultado2, "::")
						for _, val := range arr {
							if val != "" {
								arr_pais = strings.Split(val, ";")
								output += fmt.Sprintf("<option value='pais:.:%s'>%s</option>", arr_pais[0], arr_pais[1])
							}
						}
						separator := strings.Split(estado_destino, ".")
						arr_alm := separator[:len(separator)-3]
						for _, v := range arr_alm {
							res += v + "."
						}
						estado_destino = res + arr_pais[2] + ".*"
						output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #2E8B57'></span>"
						fmt.Fprint(w, output)
					}
				} else {
					//Enviamos nombre de usuario e id_region recogido en el formulario hacia el server para generar los destinos
					resultado := libs.GenerateFORM(serverext["serverroot"]+"/destino.cgi", "action;regiones", "userAdmin;"+username, "id_region;"+ident)
					if resultado != "" {
						var arr_provincia []string
						var res string
						output = "<option value='provincia:.:0'>...</option>"
						arr := strings.Split(resultado, "::")
						for _, val := range arr {
							if val != "" {
								arr_provincia = strings.Split(val, ";")
								output += fmt.Sprintf("<option value='provincia:.:%s'>%s</option>", arr_provincia[0], arr_provincia[1])
							}
						}
						last_region = ident
						separator := strings.Split(estado_destino, ".")
						arr_reg := separator[:len(separator)-1]
						for _, v := range arr_reg {
							res += v + "."
						}
						estado_destino = res + arr_provincia[2] + ".*"
						output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #2E8B57'></span>"
						fmt.Fprint(w, output)
					} else {
						//Si no hay resultado, volvemos a cargar las regiones
						var arr_region []string
						var res string
						resultado2 := libs.GenerateFORM(serverext["serverroot"]+"/destino.cgi", "action;paises", "userAdmin;"+username, "id_pais;"+last_pais)
						if resultado2 != "" {
							output = "<option value='region:.:0'>...</option>"
							arr := strings.Split(resultado2, "::")
							for _, val := range arr {
								if val != "" {
									arr_region = strings.Split(val, ";")
									output += fmt.Sprintf("<option value='region:.:%s'>%s</option>", arr_region[0], arr_region[1])
								}
							}
							separator := strings.Split(estado_destino, ".")
							arr_pais := separator[:len(separator)-2]
							for _, v := range arr_pais {
								res += v + "."
							}
							estado_destino = res + arr_region[2] + ".*"
							output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #FF0303'>No hay sub-organizaciones</span>"
							fmt.Fprint(w, output)
						}
					}
				}
			}
			if destino == "provincia" {
				if ident == "0" {
					//Si no hay resultado, volvemos a cargar las regiones
					var arr_region []string
					var res string
					resultado2 := libs.GenerateFORM(serverext["serverroot"]+"/destino.cgi", "action;paises", "userAdmin;"+username, "id_pais;"+last_pais)
					if resultado2 != "" {
						output = "<option value='region:.:0'>...</option>"
						arr := strings.Split(resultado2, "::")
						for _, val := range arr {
							if val != "" {
								arr_region = strings.Split(val, ";")
								output += fmt.Sprintf("<option value='region:.:%s'>%s</option>", arr_region[0], arr_region[1])
							}
						}
						separator := strings.Split(estado_destino, ".")
						arr_pais := separator[:len(separator)-3]
						for _, v := range arr_pais {
							res += v + "."
						}
						estado_destino = res + arr_region[2] + ".*"
						output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #2E8B57'></span>"
						fmt.Fprint(w, output)
					}
				} else {
					//Enviamos nombre de usuario e id_provincia recogido en el formulario hacia el server para generar los destinos
					resultado := libs.GenerateFORM(serverext["serverroot"]+"/destino.cgi", "action;provincias", "userAdmin;"+username, "id_provincia;"+ident)
					if resultado != "" {
						var arr_tienda []string
						var res string
						output = "<option value='tienda:.:0'>...</option>"
						arr := strings.Split(resultado, "::")
						for _, val := range arr {
							if val != "" {
								arr_tienda = strings.Split(val, ";")
								output += fmt.Sprintf("<option value='tienda:.:%s'>%s</option>", arr_tienda[0], arr_tienda[1])
							}
						}
						last_prov = ident
						separator := strings.Split(estado_destino, ".")
						arr_prov := separator[:len(separator)-1]
						for _, v := range arr_prov {
							res += v + "."
						}
						estado_destino = res + arr_tienda[2] + ".*"
						output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #2E8B57'></span>"
						fmt.Fprint(w, output)
					} else {
						//Si no hay resultado, volvemos a cargar las regiones
						var arr_provincia []string
						var res string
						resultado2 := libs.GenerateFORM(serverext["serverroot"]+"/destino.cgi", "action;regiones", "userAdmin;"+username, "id_region;"+last_region)
						if resultado2 != "" {
							output = "<option value='provincia:.:0'>...</option>"
							arr := strings.Split(resultado2, "::")
							for _, val := range arr {
								if val != "" {
									arr_provincia = strings.Split(val, ";")
									output += fmt.Sprintf("<option value='provincia:.:%s'>%s</option>", arr_provincia[0], arr_provincia[1])
								}
							}
							separator := strings.Split(estado_destino, ".")
							arr_reg := separator[:len(separator)-2]
							for _, v := range arr_reg {
								res += v + "."
							}
							estado_destino = res + arr_provincia[2] + ".*"
							output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #FF0303'>No hay sub-organizaciones</span>"
							fmt.Fprint(w, output)
						}
					}
				}
			}
			if destino == "tienda" {
				if ident == "0" {
					//Si no hay resultado, volvemos a cargar las regiones
					var arr_provincia []string
					var res string
					resultado2 := libs.GenerateFORM(serverext["serverroot"]+"/destino.cgi", "action;regiones", "userAdmin;"+username, "id_region;"+last_region)
					if resultado2 != "" {
						output = "<option value='provincia:.:0'>...</option>"
						arr := strings.Split(resultado2, "::")
						for _, val := range arr {
							if val != "" {
								arr_provincia = strings.Split(val, ";")
								output += fmt.Sprintf("<option value='provincia:.:%s'>%s</option>", arr_provincia[0], arr_provincia[1])
							}
						}
						separator := strings.Split(estado_destino, ".")
						arr_reg := separator[:len(separator)-3]
						for _, v := range arr_reg {
							res += v + "."
						}
						estado_destino = res + arr_provincia[2] + ".*"
						output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #2E8B57'></span>"
						fmt.Fprint(w, output)
					}
				} else {
					//Enviamos nombre de usuario e id_provincia recogido en el formulario hacia el server para generar los destinos
					resultado := libs.GenerateFORM(serverext["serverroot"]+"/destino.cgi", "action;tiendas", "userAdmin;"+username, "id_tienda;"+ident)
					if resultado != "" {
						var res string
						output = "<option value='destino_final:.:0'>...</option>"
						last_tienda = ident
						separator := strings.Split(estado_destino, ".")
						arr_prov := separator[:len(separator)-1]
						for _, v := range arr_prov {
							res += v + "."
						}
						estado_destino = res + resultado
						output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #2E8B57'></span>"
						fmt.Fprint(w, output)
					} else {
						//Si no hay resultado, volvemos a cargar las provincias
						var arr_tienda []string
						var res string
						resultado2 := libs.GenerateFORM(serverext["serverroot"]+"/destino.cgi", "action;provincias", "userAdmin;"+username, "id_provincia;"+last_tienda)
						if resultado2 != "" {
							output = "<option value='tienda:.:0'>...</option>"
							arr := strings.Split(resultado2, "::")
							for _, val := range arr {
								if val != "" {
									arr_tienda = strings.Split(val, ";")
									output += fmt.Sprintf("<option value='tienda:.:%s'>%s</option>", arr_tienda[0], arr_tienda[1])
								}
							}
							separator := strings.Split(estado_destino, ".")
							rest_org := separator[:len(separator)-1]
							for _, v := range rest_org {
								res += v + "."
							}
							estado_destino = res + "*"
							output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #FF0303'>No hay sub-organizaciones</span>"
							fmt.Fprint(w, output)
						}
					}
				}
			}
			if destino == "destino_final" {
				if ident == "0" {
					//Si no hay resultado, volvemos a cargar las provincias
					var arr_tienda []string
					var res string
					resultado2 := libs.GenerateFORM(serverext["serverroot"]+"/destino.cgi", "action;provincias", "userAdmin;"+username, "id_provincia;"+last_prov)
					if resultado2 != "" {
						output = "<option value='tienda:.:0'>...</option>"
						arr := strings.Split(resultado2, "::")
						for _, val := range arr {
							if val != "" {
								arr_tienda = strings.Split(val, ";")
								output += fmt.Sprintf("<option value='tienda:.:%s'>%s</option>", arr_tienda[0], arr_tienda[1])
							}
						}
						separator := strings.Split(estado_destino, ".")
						rest_org := separator[:len(separator)-1]
						for _, v := range rest_org {
							res += v + "."
						}
						estado_destino = res + "*"
						output += ";<span style='color: #1A5276'>" + estado_destino + "</span>;<span style='color: #2E8B57'></span>"
						fmt.Fprint(w, output)
					}
				}
			}
		}
	}
}

//Función que limpia de carácteres nulos la matriz salida de windows
func limpiar_matriz(matriz []byte) []byte {
	var matriz_limpiada []byte
	for _, v := range matriz {
		if v != 0 {
			matriz_limpiada = append(matriz_limpiada, v)
		}
	}
	return matriz_limpiada
}

//Función para borrar un slice vacio
func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}
