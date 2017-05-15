package main

import (
	"net/http"
	"os"
	"strings"
)

//sirve todos los ficheros mp3 tanto de publicidad como de mensajes
func root(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var publicidad, mensajes string
	accion := r.FormValue("accion") //De aqui tomamos la opcion de servir un fichero publi o msg
	//ZONA DE PUBLICIDAD
	if accion == "publicidad"{
		publicidad = strings.TrimRight(publi_files_location+r.URL.Path[1:], "/")
		filepubliinfo, err := os.Stat(publicidad)
		if err != nil {
			http.NotFound(w, r)
			return
		} else if filepubliinfo.IsDir() {
			_, err2 := os.Stat(publicidad)
			if err2 != nil {
				http.NotFound(w, r)
				return
			}
		}
		fr, errn := os.Open(publicidad)
		if errn != nil {
			http.Error(w, "Internal Server Error", 500)
			return
		}
		defer fr.Close()
		if publicidad != publi_files_location {
			//Se sirven todos los ficheros de publicidad
			http.ServeContent(w, r, publicidad, filepubliinfo.ModTime(), fr)
		}
	}else{
		//ZONA DE MENSAJES
		mensajes = strings.TrimRight(msg_files_location+r.URL.Path[1:], "/")
		filemsginfo, err := os.Stat(mensajes)
		if err != nil {
			http.NotFound(w, r)
			return
		} else if filemsginfo.IsDir() {
			_, err2 := os.Stat(mensajes)
			if err2 != nil {
				http.NotFound(w, r)
				return
			}
		}
		fm, errn := os.Open(mensajes)
		if errn != nil {
			http.Error(w, "Internal Server Error", 500)
			return
		}
		defer fm.Close()
		if mensajes != msg_files_location {
			//Se sirven todos los ficheros de publicidad
			http.ServeContent(w, r, mensajes, filemsginfo.ModTime(), fm)
		}
	}
}
