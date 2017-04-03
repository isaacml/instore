package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/isaacml/instore/libs"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	//"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	Info              *log.Logger
	Warning           *log.Logger
	Error             *log.Logger
	db                *sql.DB
	db_mu             sync.RWMutex
	serverint         map[string]string = make(map[string]string) //Mapa que guarda la direccion del servidor interno
	domainint         map[string]string = make(map[string]string) //Mapa que guarda el dominio de la tienda
	username          string                                      //Variable de usuario y estado global
	directorio_actual string                                      //Va a contener en todo momento la dirección del explorador WIN(handles_publi.go)
)

// Inicializamos la conexion a BD y el log de errores
func init() {
	var err_db error
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Fallo al abrir el archivo de error:", err)
	}
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(file, os.Stderr), "ERROR :", log.Ldate|log.Ltime|log.Lshortfile)
	//Base de datos del admin externo
	db, err_db = sql.Open("sqlite3", "C:\\instore\\music.db") //WINDB: C:\\instore\\music.db
	if err_db != nil {
		Error.Println(err_db)
		log.Fatalln("Fallo al abrir el archivo de error:", err_db)
	}
	db.Exec("PRAGMA journal_mode=WAL;")
	loadSettings(serverRoot, serverint) // Se carga los valores del fichero playerint.reg

}

// funcion principal del programa
func main() {

	fmt.Printf("Golang HTTP Server starting at Port %s ...\n", http_port)
	go controlinternalsessions() // Controla la caducidad de la sesion
	go enviar_dominio_tienda()

	// handlers del servidor HTTP
	http.HandleFunc("/", root)
	http.HandleFunc(login_cgi, login)
	http.HandleFunc(logout_cgi, logout)
	// handler de configuracion de tienda
	http.HandleFunc("/check_config.cgi", check_config)
	http.HandleFunc("/get_orgs.cgi", get_orgs)
	http.HandleFunc("/send_orgs.cgi", send_orgs)
	//Bitmap Actions
	http.HandleFunc("/bitmaps.cgi", bitmaps)
	//Exploradores
	http.HandleFunc("/explorerMusic.cgi", explorerMusic)

	s := &http.Server{
		Addr:           ":" + http_port,
		Handler:        nil,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 13,
	}

	log.Fatal(s.ListenAndServe()) // servidor HTTP multihilo
}

func enviar_dominio_tienda() {
	for {
		var existe bool
		_, err := os.Stat(configShop)
		if err != nil {
			if os.IsNotExist(err) {
				existe = false
			}
		} else {
			existe = true
		}

		if existe == true {
			loadSettings(configShop, domainint)
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverint["serverinterno"]+"/send_domain.cgi", "dominio;"+domainint["shopdomain"]))
			fmt.Println(respuesta)
		}
		time.Sleep(1 * time.Minute)
	}
}

/*
loadSettings: esta función va a abrir un fichero, leer los datos que contiene y guardarlos en un mapa.
	filename: ruta completa donde se encuentra nuestro fichero(C:\instore\serverext.reg)
*/
func loadSettings(filename string, arr map[string]string) {
	fr, err := os.Open(filename)
	defer fr.Close()
	if err == nil {
		reader := bufio.NewReader(fr)
		for {
			linea, rerr := reader.ReadString('\n')
			if rerr != nil {
				break
			}
			linea = strings.TrimRight(linea, "\n")
			item := strings.Split(linea, " = ")
			if len(item) == 2 {
				arr[item[0]] = item[1]
			}
		}
	}
}
