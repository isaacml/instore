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

var encript_files map[int]string = make(map[int]string) //Mapa que guarda los ficheros para encriptar

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
	//Toma los ficheros de la carpeta y nos genera la tabla con los ficheros MP3 correspondientes
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
	//Guarda en un mapa los archivos de encriptacion
	if r.FormValue("action") == "guardar" {
		cont := len(encript_files)
		for clave, valor := range r.Form {		
			for _, v := range valor {
				if clave == "enc_files" {
					db_mu.Lock()
					encript_files[cont] = directorio_actual + v //Se guardan los datos
					db_mu.Unlock()
					cont ++
				}
			}
		}
		//Muestra el numero de fichero que llevamos guardados para la encriptacion
		output = fmt.Sprintf("%d ficheros en la lista de encriptación", len(encript_files))
		fmt.Fprint(w, output)
	}
	//Borra todos los elementos del mapa de encriptacion
	if r.FormValue("action") == "borrar" {
		for k := range encript_files {
		    delete(encript_files, k)
		}
		output = fmt.Sprint("<span style='color: #2E8B57'>Se han borrado todos los ficheros</span>;")
		output += fmt.Sprintf("%d ficheros en la lista de encriptación", len(encript_files))
		fmt.Fprint(w, output)
	}
	//Muestra el listado de ficheros de música para encriptar
	if r.FormValue("action") == "mostrar_listado" {
		for _, v:= range encript_files {
			partir_direccion := strings.Split(v, "\\")
			mp3 := partir_direccion[len(partir_direccion)-1] // De aquí obtenemos el nombre de fichero y su extensión
		    output += fmt.Sprintf("<tr><td> %s </td></tr>", mp3)
		}
		fmt.Fprint(w, output)
	}
}
/*
func encriptacion(){
	cont := 0
	for _, v := range encript_files {
		cont ++
		partir_direccion := strings.Split(v, "\\")
		mp3 := partir_direccion[len(partir_direccion)-1] // De aquí obtenemos el nombre de fichero y su extensión
		del_ext := strings.Split(mp3, ".mp3")	//Quitamos la extensión y nos quedamos con el nombre
		cifrado := del_ext[0] + ".xxx"	//Le añadimos la extensión de encriptado(.xxx)
		//Generamos el fichero de encriptación
		libs.Cifrado(v, cifDir+cifrado, []byte{11, 22, 33, 44, 55, 66, 77, 88})
		fmt.Println(cont, len(encript_files))
		//fmt.Fprintf(w, "Cifrado %d de %d\n", cont, len(encript_files))
		time.Sleep(1*time.Second)
		transcurrido := time.Since(actual)
		fmt.Println(transcurrido)

	}
}
*/

func estado_encriptacion(w http.ResponseWriter, r *http.Request){

}