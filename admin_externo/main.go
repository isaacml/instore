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
	"net/url"
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
	serverext         map[string]string = make(map[string]string) //Mapa que guarda la direccion del servidor externo
	username, good    string                                      //Variable de usuario y estado global
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
	db, err_db = sql.Open("sqlite3", "C:\\instore\\programaciones.db") //WINDB: C:\\instore\\programaciones.db
	if err_db != nil {
		Error.Println(err_db)
		log.Fatalln("Fallo al abrir el archivo de error:", err_db)
	}
	db.Exec("PRAGMA journal_mode=WAL;")
}

// funcion principal del programa
func main() {

	fmt.Printf("Golang HTTP Server starting at Port %s ...\n", http_port)
	go controlinternalsessions() // Controla la caducidad de la sesion
	go mantenimiento()

	// handlers del servidor HTTP
	http.HandleFunc("/", root)
	http.HandleFunc(login_cgi, login)
	http.HandleFunc(logout_cgi, logout)
	//ACCIONES
	http.HandleFunc("/acciones.cgi", acciones)
	//Usuarios
	http.HandleFunc("/usuarios.cgi", usuarios)
	//Entidades
	http.HandleFunc("/entidades.cgi", entidades)
	//Almacenes
	http.HandleFunc("/almacenes.cgi", almacenes)
	//Paises
	http.HandleFunc("/paises.cgi", paises)
	//Regiones
	http.HandleFunc("/regiones.cgi", regiones)
	//Provincias
	http.HandleFunc("/provincias.cgi", provincias)
	//Tiendas
	http.HandleFunc("/tiendas.cgi", tiendas)
	//Publicidad
	http.HandleFunc("/explorer.cgi", explorer)
	http.HandleFunc("/dest_explorer.cgi", dest_explorer)
	//Encriptar Musica de JR
	http.HandleFunc("/encriptar_musica.cgi", encriptar_musica)
	http.HandleFunc("/orgs.cgi", orgs)
	http.HandleFunc("/selected_org.cgi", selected_org)
	http.HandleFunc("/estado_encriptacion.cgi", estado_encriptacion)

	s := &http.Server{
		Addr:           ":" + http_port,
		Handler:        nil,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 13,
	}

	log.Fatal(s.ListenAndServe()) // servidor HTTP multihilo
}

//MANTENIMIENTO
func mantenimiento() {
	loadSettings(serverRoot)
	for {
		var ruta, fichero, fecha_inicio, fecha_final, destino, playtime, gap string
		var id int
		//ZONA PUBLICIDAD
		publicidad, err := db.Query("SELECT id, ruta, fichero, fecha_inicio, fecha_final, destino, gap FROM publi WHERE id = (SELECT MIN(id) FROM publi)")
		if err != nil {
			Error.Println(err)
		}
		for publicidad.Next() {
			err = publicidad.Scan(&id, &ruta, &fichero, &fecha_inicio, &fecha_final, &destino, &gap)
			if err != nil {
				Error.Println(err)
			}
			rutaCompleta := ruta + fichero
			//parámetros pasados por URL
			v := url.Values{}
			v.Add("fichero", fichero)
			v.Add("f_inicio", fecha_inicio)
			v.Add("f_final", fecha_final)
			v.Add("destino", destino)
			v.Add("gap", gap)
			v.Add("ownUser", serverext["usuarioPropietario"])
			_, errUp := libs.ClienteUpload(rutaCompleta, serverext["serverroot"]+"/publi_files.cgi?"+v.Encode(), 1000, 0)
			if errUp == nil {
				db_mu.Lock()
				_, err1 := db.Exec("DELETE FROM publi WHERE id=?", id)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
				}
				continue
			}
		}
		//ZONA MENSAJES
		mensajes, err := db.Query("SELECT id, ruta, fichero, fecha_inicio, fecha_final, destino, playtime FROM mensaje WHERE id = (SELECT MIN(id) FROM mensaje)")
		if err != nil {
			Error.Println(err)
		}
		for mensajes.Next() {
			err = mensajes.Scan(&id, &ruta, &fichero, &fecha_inicio, &fecha_final, &destino, &playtime)
			if err != nil {
				Error.Println(err)
			}
			rutaCompleta := ruta + fichero
			//parámetros pasados por URL
			v := url.Values{}
			v.Add("fichero", fichero)
			v.Add("f_inicio", fecha_inicio)
			v.Add("f_final", fecha_final)
			v.Add("destino", destino)
			v.Add("playtime", playtime)
			v.Add("ownUser", serverext["usuarioPropietario"])
			_, errUp := libs.ClienteUpload(rutaCompleta, serverext["serverroot"]+"/msg_files.cgi?"+v.Encode(), 1000, 0)
			if errUp == nil {
				db_mu.Lock()
				_, err1 := db.Exec("DELETE FROM mensaje WHERE id=?", id)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
				}
				continue
			}
		}
		time.Sleep(30 * time.Second)
	}
}

/*
loadSettings: esta función va a abrir un fichero, leer los datos que contiene y guardarlos en un mapa.
	filename: ruta completa donde se encuentra nuestro fichero(C:\instore\serverext.reg)
*/
func loadSettings(filename string) {
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
				serverext[item[0]] = item[1]
			}
		}
	}
}
