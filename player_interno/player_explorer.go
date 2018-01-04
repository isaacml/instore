package main

import (
	"bufio"
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"golang.org/x/text/encoding/charmap"
)

//Guarda la dirección donde se encuentra el explorador WIN
var directorio_actual string

//Función principal del explorador windows
func explorerMusic(w http.ResponseWriter, r *http.Request) {
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
	//DIRECTORIOS A MOSTRAR
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
			output = "<option value='' selected>[Seleccion de directorios]</option>"
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
	//EXPLORADOR DE DIRECTORIOS
	if r.FormValue("action") == "directorios" {
		//VOLVER UN DIRECTORIO ATRÁS
		if r.FormValue("directory") == "..." {
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
			db_mu.Lock()
			capacidad_arr = contador
			db_mu.Unlock()
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
				output = "<option value='' selected>[Seleccion de directorios]</option>"
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
				output = fmt.Sprintf("<option value='' selected>[Seleccion de directorios]</option><option value='...'>...</option>")
				for _, val := range directorios {
					if val.IsDir() {
						output += fmt.Sprintf("<option value='%s'>%s</option>", val.Name(), val.Name())

					}
				}
				output += ";<span style='color: #B8860B'>" + directorio_actual + "</span>"
				fmt.Fprint(w, output)
			}
			//En caso de que la ruta de directorio esté vacia.
		} else if r.FormValue("directory") == "" {
			//La capacidad del array no puede ser 0
			if capacidad_arr == 1 {
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
				output = fmt.Sprintf("<option value='' selected>[Seleccion de directorios]</option>")
				for _, val := range directorios {
					if val.IsDir() {
						output += fmt.Sprintf("<option value='%s'>%s</option>", val.Name(), val.Name())

					}
				}
				output += ";<span style='color: #B8860B'>" + directorio_actual + "</span>"
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
				output = fmt.Sprintf("<option value='' selected>[Seleccion de directorios]</option><option value='...'>...</option>")
				for _, val := range directorios {
					if val.IsDir() {
						output += fmt.Sprintf("<option value='%s'>%s</option>", val.Name(), val.Name())

					}
				}
				output += ";<span style='color: #B8860B'>" + directorio_actual + "</span>"
			}
			fmt.Fprint(w, output)
		} else {
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
				output = "<option value='' selected>[Seleccion de directorios]</option></option><option value='...'>...</option>"
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
			output = fmt.Sprintf("<option value='' selected>[Seleccion de directorios]</option><option value='...'>...</option>")
			for _, val := range directorios {
				if val.IsDir() {
					output += fmt.Sprintf("<option value='%s'>%s</option>", val.Name(), val.Name())

				}
			}
			output += ";<span style='color: #B8860B'>" + directorio_actual + "</span>"
			fmt.Fprint(w, output)
		}
	}
	//Se toman las carpetas enviadas por el formulario (addMusic.html) y se realiza la copia en el directorio de musica de la tienda.
	//Se puede seleccionar cualquier carpeta que contenga cualquier tipo de archivo, pero solo se copiaran los mp3.
	if r.FormValue("action") == "uploadDirs" {
		var bat string
		/*
		//Creamos un fichero(.bat) que va a guardar y ejecutar los comandos
		addMusic, err := os.Create("addMusic.bat")
		if err != nil {
			Error.Println(err)
			return
		}
		defer addMusic.Close()
		*/
		//Tomamos las carpetas añadidas por el usuario
		for clave, valor := range r.Form {
			for _, v := range valor {
				if clave == "directory" {
					if v == "" || v == "..." {
						output += "<span style='color: #800000'>Debes seleccionar mínimo un directorio</span>"
						fmt.Fprint(w, output)
						return
					}
					
					fmt.Println("directorios: ", v)
					//Ejemplo: "C:\Users\isaac\miMusica\ACDC\*.mp3"
					bat += "dir /s /b \""+directorio_actual + v + "\\*.mp3\"\r\n"
					/*
					//Mandamos el array de directorios y el directorio destino para realizar la copia
					libs.FileCopier(copy_arr, music_files)
					*/
				}
			}
		}
		err := exec.Command("cmd", "/c", "echo "+bat+" > addMusic.bat").Run()
		if err != nil {
			Error.Println(err)
			return
		}
		//Abrimos el txt que va a contener todos los ficheros de música correspondientes
		musica, err := os.Open("addMusic.txt")
		if err != nil {
			Error.Println(err)
			return
		}
		defer musica.Close()
		scanner := bufio.NewScanner(musica)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}
}

func decodeWindows1250(enc []byte) string {
    dec := charmap.Windows1250.NewDecoder()
    out, _ := dec.Bytes(enc)
    return string(out)
}

func encodeWindows1250(inp string) []byte {
    enc := charmap.Windows1250.NewEncoder()
    out, _ := enc.String(inp)
    return []byte(out)
}
