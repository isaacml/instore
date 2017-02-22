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
	// handlers del administrador externo
	http.HandleFunc("/user_admin.cgi", user_admin)
	//Usuario
	http.HandleFunc("/edit_own_user.cgi", edit_own_user)
	http.HandleFunc("/alta_users.cgi", alta_users)
	http.HandleFunc("/get_users.cgi", get_users)
	http.HandleFunc("/load_user.cgi", load_user)
	http.HandleFunc("/edit_user.cgi", edit_user)
	http.HandleFunc("/user_entidad.cgi", user_entidad)
	http.HandleFunc("/user_permiso.cgi", user_permiso)
	//Entidad
	http.HandleFunc("/entidad.cgi", entidad)
	http.HandleFunc("/get_entidad.cgi", get_entidad)
	http.HandleFunc("/load_entidad.cgi", load_entidad)
	http.HandleFunc("/edit_entidad.cgi", edit_entidad)
	//Almacen
	http.HandleFunc("/almacen_entidad.cgi", almacen_entidad)
	http.HandleFunc("/almacen.cgi", almacen)
	http.HandleFunc("/get_almacen.cgi", get_almacen)
	http.HandleFunc("/load_almacen.cgi", load_almacen)
	http.HandleFunc("/edit_almacen.cgi", edit_almacen)
	//Pais
	http.HandleFunc("/pais_almacen.cgi", pais_almacen)
	http.HandleFunc("/pais.cgi", pais)
	http.HandleFunc("/get_pais.cgi", get_pais)
	http.HandleFunc("/load_pais.cgi", load_pais)
	http.HandleFunc("/edit_pais.cgi", edit_pais)
	//Region
	http.HandleFunc("/region_pais.cgi", region_pais)
	http.HandleFunc("/region.cgi", region)
	http.HandleFunc("/get_region.cgi", get_region)
	http.HandleFunc("/load_region.cgi", load_region)
	http.HandleFunc("/edit_region.cgi", edit_region)
	//Provincia
	http.HandleFunc("/provincia_region.cgi", provincia_region)
	http.HandleFunc("/provincia.cgi", provincia)
	http.HandleFunc("/get_provincia.cgi", get_provincia)
	http.HandleFunc("/load_provincia.cgi", load_provincia)
	http.HandleFunc("/edit_provincia.cgi", edit_provincia)
	//Tienda
	http.HandleFunc("/tienda_provincia.cgi", tienda_provincia)
	http.HandleFunc("/tienda.cgi", tienda)
	http.HandleFunc("/get_tienda.cgi", get_tienda)
	http.HandleFunc("/load_tienda.cgi", load_tienda)
	http.HandleFunc("/edit_tienda.cgi", edit_tienda)
	//Publicidad
	http.HandleFunc("/explorer.cgi", explorer)
	http.HandleFunc("/recoger_destinos.cgi", recoger_destinos)

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
		var ruta, fichero, fecha_inicio, fecha_final, destino string
		var id int
		query, err := db.Query("SELECT id, ruta, fichero, fecha_inicio, fecha_final, destino FROM publi WHERE id = (SELECT MIN(id) FROM publi)")
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			err = query.Scan(&id, &ruta, &fichero, &fecha_inicio, &fecha_final, &destino)
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
			v.Add("ownUser", serverext["usuarioPropietario"])
			_, errUp := libs.ClienteUpload(rutaCompleta, serverext["serverroot"]+"/get_files.cgi?"+v.Encode(), 1000, 0)
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
		time.Sleep(20 * time.Second)
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
