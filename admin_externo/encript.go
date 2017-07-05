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
					output += fmt.Sprintf("<tr><td><input type='checkbox' name='enc_files' value='%s'></td><td> %s </td></tr>", val.Name(), val.Name())
				}
			}
		}
		fmt.Fprint(w, output)
	}
	//TOMAMOS LOS FICHEROS DE MUSICA DEL FORMULARIO 2
	if r.FormValue("action") == "tomar_musica" {
		cont := 0
		actual := time.Now()
		for clave, valor := range r.Form {		
			for _, v := range valor {
				if clave == "enc_files" {
					cont ++
					del_ext := strings.Split(v, ".mp3")
					cifrado := del_ext[0] + ".xxx"
					libs.Cifrado(directorio_actual+v, cifDir+cifrado, []byte{11, 22, 33, 44, 55, 66, 77, 88})
					transcurrido := time.Since(actual)
					fmt.Println(transcurrido)
					output = fmt.Sprintf("Cifrado %d de %d ", cont, len(valor))
					
				}
			}
		}
		fmt.Fprint(w, output)
	}
}