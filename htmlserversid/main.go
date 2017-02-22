package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

// Inicializamos la conexion a BD y el log de errores
func init() {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalln("Fallo al abrir el archivo de error:", err)
	}

	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(file, os.Stderr), "ERROR :", log.Ldate|log.Ltime|log.Lshortfile)

}

// funcion principal del programa
func main() {

	fmt.Printf("Golang HTTP Server starting at Port %s ...\n", http_port)
	go controlinternalsessions() // Controla la caducidad de la sesion

	// handlers del servidor HTTP
	http.HandleFunc("/", root)
	http.HandleFunc(login_cgi, login)
	http.HandleFunc(logout_cgi, logout)
	// handlers del encoder

	http.HandleFunc("/newPage.cgi", newPage)

	s := &http.Server{
		Addr:           ":" + http_port,
		Handler:        nil,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 13,
	}

	log.Fatal(s.ListenAndServe()) // servidor HTTP multihilo
}

func newPage(w http.ResponseWriter, r *http.Request) {
	if checkCGI(r) == true {
		r.ParseForm()
		fmt.Printf("Usuario: %s - Contraseña: %s \n", r.FormValue("usuario"), r.FormValue("password"))
	} else {
		http.Redirect(w, r, "/"+first_page+".html", http.StatusFound)
	}
}
