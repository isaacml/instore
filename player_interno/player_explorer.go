package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"bufio"
	"github.com/isaacml/instore/libs"
)
//Guarda la capacidad que tiene el array que guarda la ruta de directorio
var capacidad_arr int
//Contenedor que va a guardar los ficheros que van a ser copiados a "C:\instore\\Music\"
var copy_arr []string

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
	//DIRECTORIOS A MOSTRAR
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
			output = "<option value='' selected>[Seleccion de directorios]</option>"
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
					arr_sin_vacios = RemoveIndex(ruta, k)
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
		}else if r.FormValue("directory") == "" {
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
			}else{
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
		}else{
			directorio_actual = directorio_actual + r.FormValue("directory") + "\\"
			file, err := os.Open(directorio_actual)
			defer file.Close()
			if err != nil {
				// No se puede abrir el directorio, por falta de permisos
				Error.Println(err)
				//Volvemos a tomar el archivo anterior y lo abrimos
				old := strings.Split(directorio_actual, r.FormValue("directory")+"\\")
				directorio_actual = old[0]
				fmt.Println(directorio_actual)
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
	//Se toman las carpetas enviadas por el formulario (musicNoCif.html) y se realiza la copia en el directorio de musica de la tienda.
	//Se puede seleccionar cualquier carpeta que contenga cualquier tipo de archivo, pero solo se copiaran los mp3.
	if r.FormValue("action") == "noCifDirs" {
		for clave, valor := range r.Form {
			for _, v := range valor {
				if clave == "directory" {
					if v == "" {
						output += "<span style='color: #FF0303'>Debes seleccionar mínimo un directorio</span>"
						fmt.Fprint(w, output)
						return
					}
					//Ejemplo: C:\Users\isaac\miMusica\
					selected_dir := directorio_actual+"\\"+v+"\\"
					//Va a listar todos los ficheros mp3, tanto del directorio padre como sus hijos.
					cmd := exec.Command("cmd", "/c", "dir /s /b "+selected_dir+"*.mp3")
				    //Comienza la ejecucion del pipe
					stdoutRead, _ := cmd.StdoutPipe()
					reader := bufio.NewReader(stdoutRead)
					cmd.Start()
				    for {
						line, err := reader.ReadString('\n')
						if err != nil {
							break
						}
						db_mu.Lock()
						//Guardamos los ficheros en el array de copia
						copy_arr = append(copy_arr, strings.TrimRight(line, "\r\n"))
						db_mu.Unlock()
					}
					//Mandamos el array de directorios y el directorio destino para realizar la copia
					libs.FileCopier(copy_arr, music_files)
				}
			}
		}
	}
	//Se toman las carpetas enviadas por el formulario (musicCif.html) y se realiza la copia en el directorio de musica de la tienda.
	//Se puede seleccionar cualquier carpeta que contenga cualquier tipo de archivo, pero solo se copiaran los ficheros encriptados (.xxx)
	if r.FormValue("action") == "cifDirs" {
		for clave, valor := range r.Form {
			for _, v := range valor {
				if clave == "directory" {
					if v == "" {
						output += "<span style='color: #FF0303'>Debes seleccionar mínimo un directorio</span>"
						fmt.Fprint(w, output)
						return
					}
					//Ejemplo: C:\Users\isaac\miMusica\
					selected_dir := directorio_actual+"\\"+v+"\\"
					//Va a listar todos los ficheros .xxx, tanto del directorio padre como sus hijos.
					cmd := exec.Command("cmd", "/c", "dir /s /b "+selected_dir+"*.xxx")
				    //Comienza la ejecucion del pipe
					stdoutRead, _ := cmd.StdoutPipe()
					reader := bufio.NewReader(stdoutRead)
					cmd.Start()
				    for {
						line, err := reader.ReadString('\n')
						if err != nil {
							break
						}
						db_mu.Lock()
						//Guardamos los ficheros en el array de copia
						copy_arr = append(copy_arr, strings.TrimRight(line, "\r\n"))
						db_mu.Unlock()
					}
					//Mandamos el array de directorios y el directorio destino para realizar la copia
					libs.FileCopier(copy_arr, music_files)
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

//Función para borrar espacios vacios dentro de un array
func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}
