package main

import (
	"fmt"
	//"github.com/isaacml/instore/libs"
	"net/http"
	"os"
	"os/exec"
	"strings"
	//"time"
)


func encriptar_musica(w http.ResponseWriter, r *http.Request) {
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
			output = "<option value='' selected>[Selecciona un directorio]</option>"
			for _, val := range directorios {
				if val.IsDir() {
					output += fmt.Sprintf("<option value='%s'>%s</option>", val.Name(), val.Name())

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
			if !val.IsDir() {
				if strings.Contains(val.Name(), ".mp3"){
					output += fmt.Sprintf("<input type='checkbox' name='enc_files'  value='%s'> %s <br>", val.Name(), val.Name())
				}
			}
		}
		fmt.Fprint(w, output)
	}
	//TOMA LOS FICHEROS DEL FORMULARIO 2, Y LOS PROCESA
	//r.FormValue("type") == "publi", procedemos a insertar los datos en la tabla publi
	if r.FormValue("action") == "tomar_musica" {
		fmt.Println(r.Form)
		/*
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
		}*/
	}
}