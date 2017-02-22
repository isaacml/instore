package main

import (
	"fmt"
	"github.com/todostreaming/realip"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"strings"
	"time"
)

// sirve todos los ficheros estáticos de la web html,css,js,graficos,etc
func root(w http.ResponseWriter, r *http.Request) {
	var namefile string
	namefile = strings.TrimRight(rootdir+r.URL.Path[1:], "/")
	//fmt.Println("... Buscando fichero: ", namefile)
	fileinfo, err := os.Stat(namefile)
	if err != nil {
		// fichero no existe
		//fmt.Println("404 - Fichero no encontrado: ",namefile)
		http.NotFound(w, r)
		return
	} else if fileinfo.IsDir() {
		// es un directorio, luego le añadimos index.html
		namefile = namefile + "/" + first_page + ".html"
		//fmt.Println("/ - Entramos en el Directorio Buscando el fichero: ",namefile)
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
	defer fr.Close()

	if strings.Contains(r.URL.String(), "?err") {
		fmt.Println("Ruta: " + namefile + " contiene ?err")
		// sustituimos <span id="loginerr"></span> por un texto de error a mostrar
		buf, _ := ioutil.ReadAll(fr) // leemos el HTML template de una sola vez
		html := string(buf)
		// Vamos a meter los campos options creados en bmdinfo() en el HTML
		html = strings.Replace(html, spanHTMLlogerr, ErrorText, -1)
		w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
		fmt.Fprint(w, html)
	} else {
		file := strings.Split(namefile, ".")
		if (file[1] != "html") || strings.Contains(namefile, "/"+first_page+".html") {
			http.ServeContent(w, r, namefile, fileinfo.ModTime(), fr)
		} else {
			//Se comprueba que el link contiene la parte {{sid}}
			if strings.Contains(r.URL.String(), "?") {
				sid := strings.Split(r.URL.String(), "?")
				for k, _ := range user {
					dir_ip := realip.RealIP(r)
					if k == sid[1] && dir_ip == ip[sid[1]] {
						var nivel int
						buf, _ := ioutil.ReadAll(fr) // leemos el HTML template de una sola vez
						html := string(buf[5:])
						fmt.Sscanf(string(buf[:5]), "[[%d]]", &nivel) //Obtengo el nivel del HTML
						// Si el nivel[0-9] del mapa de usuario es superior al del html, no dejo acceso
						if level[sid[1]] > nivel {
							http.Redirect(w, r, "/"+first_page+".html", http.StatusFound)
							return
						}
						//Aumento el tiempo de expiración
						expires := time.Now().Unix() + int64(session_timeout)
						tiempo[sid[1]] = expires
						//Reemplazamos todos los enlaces para que tengan el mismo {{sid}}
						html = strings.Replace(html, "?{{sid}}", fmt.Sprintf("?%s", sid[1]), -1)
						//Reemplazamos todos los cgi para que tengan el mismo {{sid}}
						html = strings.Replace(html, "?sid=sid", fmt.Sprintf("?sid=%s", sid[1]), -1)
						w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
						fmt.Fprint(w, html)
						return
					}
				}
			}
			//Sino el SID no se corresponde, lo tiramos para fuera
			http.Redirect(w, r, "/"+first_page+".html", http.StatusFound)
		}
	}
	r.Body.Close()
}

//revisa que el SID para los CGI sea el correcto
func checkCGI(r *http.Request) bool {
	r.ParseForm()
	var estado bool
	sid := r.FormValue("sid")
	for k, _ := range user {
		dir_ip := realip.RealIP(r)
		if k == sid && dir_ip == ip[sid] {
			expires := time.Now().Unix() + int64(session_timeout)
			tiempo[sid] = expires
			estado = true
		} else {
			estado = false
		}
	}
	return estado

}
