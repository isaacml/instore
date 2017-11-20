package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//Mapa que guarda los directorios que queremos cifrar
var encripts_dirs map[int]string = make(map[int]string)

//Estado de encriptacion: muestra los directorios que se están cifrando
var estado_encript string

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
		limpiar := strings.TrimSpace(string(libs.LimpiarMatriz([]byte(res[1]))))
		unidades := strings.Split(limpiar, "\\")
		for _, v := range unidades {
			v = strings.TrimSpace(v)
			if v != "" {
				output += fmt.Sprintf("<option value='%s'>%s</option>", v, v)
			}
		}
		output += fmt.Sprintf(";<span style='color: #2E8B57'>%d directorios para la encriptación</span>", len(encripts_dirs))
		fmt.Fprint(w, output)
	}

	//EXPLORADOR DE DIRECTORIOS (PRIMERA EJECUCION)
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
					output += fmt.Sprintf("<option value='%s'>%s</option>", val.Name(), val.Name())
				}
			}
			output += ";<span style='color: #B8860B'>" + directorio_actual + "</span>"
		} else {
			directorio_actual = ""
			output += ";<span style='color: #B8860B'>" + directorio_actual + "</span>"
		}
		fmt.Fprint(w, output)
	}
	//EXPLORADOR DE DIRECTORIOS (AL HACER DOBLE CLICK)
	if r.FormValue("action") == "directorios" {
		if r.FormValue("directory") != "" && r.FormValue("directory") != "..." {
			directorio_actual = directorio_actual + r.FormValue("directory") + "\\"
			file, err := os.Open(directorio_actual)
			defer file.Close()
			//No se puede abrir el directorio, por falta de permisos
			if err != nil {
				//Tenemos que volver a cargar el directorio anterior
				old := strings.Split(directorio_actual, "\\")
				//Puede ocurrir que al retroceder en un directorio nos encontremos con la unidad por tanto:
				if len(old) == 3 { //longitud = 3 --> Unidad de disco + directorio + contrabarra
					back_dir := old[:len(old)-2] //Quito el directorio y la contrabarra
					//Guardamos la nueva ruta
					directorio_actual = back_dir[0] + "\\"
					//Abrimos el directorio y mostramos todas sus carpetas
					file2, err := os.Open(directorio_actual)
					defer file.Close()
					directorios, err := file2.Readdir(0)
					if err != nil {
						Error.Println(err)
						return
					}
					for _, val := range directorios {
						if val.IsDir() {
							output += fmt.Sprintf("<option value='%s'>%s</option>", val.Name(), val.Name())
						}
					}
					output += ";<span style='color: #800000'>Necesitas permisos para abrir ese directorio</span>"
					fmt.Fprint(w, output)
					return
					//Otro directorio distinto a la unidad
				} else {
					var back_dir string
					old = old[:len(old)-2] //Quito el directorio y la contrabarra
					for _, v := range old {
						back_dir += v + "\\"
					}
					//Guardamos la nueva ruta
					directorio_actual = back_dir
					//Abrimos el directorio y mostramos todas sus carpetas
					file2, err := os.Open(directorio_actual)
					defer file.Close()
					directorios, err := file2.Readdir(0)
					if err != nil {
						Error.Println(err)
						return
					}
					output = "<option value='...'>...</option>"
					for _, val := range directorios {
						if val.IsDir() {
							output += fmt.Sprintf("<option value='%s'>%s</option>", val.Name(), val.Name())

						}
					}
					output += ";<span style='color: #800000'>Necesitas permisos para abrir ese directorio</span>"
					fmt.Fprint(w, output)
					return
				}
			}
			//Abrimos el directorio y mostramos sus carpetas
			directorios, err := file.Readdir(0)
			if err != nil {
				Error.Println(err)
				return
			}
			output = fmt.Sprintf("<option value='...'>...</option>")
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
				output = fmt.Sprintf("<option value='...'>...</option>")
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
			output = fmt.Sprintf("<option value='...'>...</option>")
			for _, val := range directorios {
				if val.IsDir() {
					output += fmt.Sprintf("<option value='%s'>%s</option>", val.Name(), val.Name())

				}
			}
			output += ";<span style='color: #B8860B'>" + directorio_actual + "</span>"
			fmt.Fprint(w, output)

		}
	}
	//Guarda en un mapa los directorios de encriptacion
	if r.FormValue("action") == "guardar" {
		cont := len(encripts_dirs)
		for clave, valor := range r.Form {
			for _, v := range valor {
				if clave == "directory" {
					db_mu.Lock()
					encripts_dirs[cont] = directorio_actual + v //Se guardan los datos
					db_mu.Unlock()
					cont++
				}
			}
		}
		//Muestra el numero de directorios que llevamos guardados para la encriptacion
		output = fmt.Sprintf("<span style='color: #2E8B57'>%d directorios para la encriptación</span>", len(encripts_dirs))
		fmt.Fprint(w, output)
	}
	//Borra todos los directorios del mapa de encriptacion
	if r.FormValue("action") == "borrar" {
		for k := range encripts_dirs {
			delete(encripts_dirs, k)
		}
		output = fmt.Sprintf("<span style='color: #2E8B57'>%d directorios para la encriptación</span>", len(encripts_dirs))
		fmt.Fprint(w, output)
	}
	//Muestra el listado de directorios de música para encriptar
	if r.FormValue("action") == "mostrar_listado" {
		for k, v := range encripts_dirs {
			partir_direccion := strings.Split(v, "\\")
			mp3 := partir_direccion[len(partir_direccion)-1] // De aquí obtenemos el nombre de fichero y su extensión
			output += fmt.Sprintf("<tr><td> %s </td><td><button title='Borra este archivo del listado' class='btn btn-md btn-warning' onclick=\"delfile(%d)\"><i class='fa fa-trash-o'></i></button></td></tr>", mp3, k)
		}
		fmt.Fprint(w, output)
	}
	//Borrar un directorio del listado de encriptacion
	if r.FormValue("action") == "borrar_de_listado" {
		id := r.FormValue("id_directory")
		intID, err := strconv.Atoi(id)
		if err != nil {
			Error.Println(err)
			return
		}
		db_mu.Lock()
		delete(encripts_dirs, intID)
		db_mu.Unlock()
	}
	//Encripta un listado completo de directorios y subdirectorios
	if r.FormValue("action") == "encriptar" {
		for k, v := range encripts_dirs {
			//recorre recursivamente un directorio, se obtienen todos los subdirectorios y ficheros que cuelgan de el
			filepath.Walk(v, func(path string, f os.FileInfo, err error) error {
				dir := f.IsDir()
				//Obtenemos unicamente los ficheros
				if dir == false {
					//solo MP3
					if strings.Contains(path, ".mp3") {
						partir_direccion := strings.Split(path, "\\")
						mp3 := partir_direccion[len(partir_direccion)-1] // De aquí obtenemos el nombre de fichero y su extensión
						del_ext := strings.Split(mp3, ".mp3")            //Quitamos la extensión y nos quedamos con el nombre
						cifrado := del_ext[0] + ".xxx"                   //Le añadimos la extensión de encriptado(.xxx)
						//Generamos el fichero de encriptación
						libs.Cifrado(path, cifDir+cifrado, []byte{11, 22, 33, 44, 55, 66, 77, 88})
					}
				}
				return nil
			})
			//Borramos del listado la carpeta, una vez encriptada
			db_mu.Lock()
			delete(encripts_dirs, k)
			//Enviamos un estado indicando la carpeta que se ha cifrado
			estado_encript += fmt.Sprintf("<tr><td> %s </td><td><img src='img/Ok.png' title='Carpeta Cifrada' alt='Carpeta Cifrada'></td></tr>", v)
			db_mu.Unlock()
			time.Sleep(1 * time.Second)
			if len(encripts_dirs) == 0 {
				//Mensaje final, cuando termina el cifrado
				db_mu.Lock()
				estado_encript = "<span style='color: #2E8B57'>Todas las carpetas han sido cifradas</span>"
				db_mu.Unlock()
			}
		}
		time.Sleep(1 * time.Second)
		estado_encript = ""
	}
}

//Muestra los directorios cifrados cada 1seg
func estado_encriptacion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, estado_encript)
}
