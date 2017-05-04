package main

import (
	"fmt"
	//"github.com/todostreaming/realip"
	//"io/ioutil"
	//"mime"
	"net/http"
	"os"
	"strings"
	//"time"
)

// sirve todos los ficheros estáticos de la web html,css,js,graficos,etc
func root(w http.ResponseWriter, r *http.Request) {
	var namefile string
	namefile = strings.TrimRight(publi_files_location+r.URL.Path[1:], "/")
	fmt.Println("... Buscando fichero: ", namefile)
	fileinfo, err := os.Stat(namefile)
	if err != nil {
		http.NotFound(w, r)
		return
	} else if fileinfo.IsDir() {
		_, err2 := os.Stat(namefile)
		if err2 != nil {
			//fmt.Println("404 - Fichero no encontrado: ",namefile)
			http.NotFound(w, r)
			return
		}
	}
	fr, errn := os.Open(namefile)
	if errn != nil {
		//fmt.Printf("[root(4)] - Error de acceso al fichero: %s\n",namefile)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	http.ServeFile(w, r, namefile)
}
