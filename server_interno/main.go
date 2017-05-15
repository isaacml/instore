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
	"os"
	"strings"
	"sync"
	"time"
)

var (
	Info      *log.Logger
	Warning   *log.Logger
	Error     *log.Logger
	db        *sql.DB
	db_mu     sync.RWMutex
	serverext map[string]string = make(map[string]string) //Mapa que guarda la direccion del servidor externo
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
	db, err_db = sql.Open("sqlite3", "/home/isaac/INSTORE_SQL/pruebas/servint.db")
	if err_db != nil {
		Error.Println(err_db)
		log.Fatalln("Fallo al abrir el archivo de error:", err_db)
	}
	db.Exec("PRAGMA journal_mode=WAL;")
	loadSettings(serverRoot) // Se carga los valores del fichero serverint.reg
}

// funcion principal del programa
func main() {

	fmt.Printf("Golang HTTP Server starting at Port %s ...\n", http_port)
	go BuscarNuevosFicheros()

	// handlers de la tienda
	http.HandleFunc("/login_tienda.cgi", login_tienda)
	http.HandleFunc("/transf_orgs.cgi", transf_orgs)
	http.HandleFunc("/send_orgs.cgi", send_orgs)
	//Actions
	http.HandleFunc("/bitmaps.cgi", bitmaps)
	http.HandleFunc("/send_domain.cgi", send_domain)
	http.HandleFunc("/downloadPubliFile.cgi", downloadPubliFile)
	http.HandleFunc("/downloadMsgFile.cgi", downloadMsgFile)

	s := &http.Server{
		Addr:           ":" + http_port,
		Handler:        nil,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 13,
	}

	log.Fatal(s.ListenAndServe()) // servidor HTTP multihilo
}

func BuscarNuevosFicheros() {
	for {
		//Buscamos todos lo ficheros de publicidad que no tenemos en la BD
		publicidad, errP := db.Query("SELECT fichero FROM publi WHERE existe='N'")
		if errP != nil {
			Error.Println(errP)
		}
		for publicidad.Next() {
			var fichero string
			//Tomamos el nombre del fichero publicidad
			err := publicidad.Scan(&fichero)
			if err != nil {
				Error.Println(err)
			}
			//Descargamos el fichero del servidor externo
			bytes, err := libs.DownloadFile(serverext["serverexterno"]+"/"+fichero+"?accion=publicidad", publi_files_location+fichero, 0, 1000)
			//bytes igual a 0 o error diferente de nulo: la descarga ha ido mal
			if err != nil || bytes == 0 {
				Error.Println(err)
			}
			//bytes distintos de 0 o error igual a nulo: la descarga se ha realizado correctamente.
			if bytes != 0 || err == nil {
				//Cambiamos el estado del fichero de publicidad en BD, a existe.
				ok, err := db.Prepare("UPDATE publi SET existe=? WHERE fichero = ?")
				if err != nil {
					Error.Println(err)
				}
				db_mu.Lock()
				_, err1 := ok.Exec("Y", fichero)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
				}
			}
		}
		//Buscamos todos lo ficheros de mensajes que no tenemos en la BD
		mensajes, errM := db.Query("SELECT fichero FROM mensaje WHERE existe='N'")
		if errM != nil {
			Error.Println(errM)
		}
		for mensajes.Next() {
			var fichero string
			//Tomamos el nombre del fichero mensaje
			err := mensajes.Scan(&fichero)
			if err != nil {
				Error.Println(err)
			}
			//Descargamos el fichero del servidor externo
			bytes, err := libs.DownloadFile(serverext["serverexterno"]+"/"+fichero+"?accion=mensaje", msg_files_location+fichero, 0, 1000)
			//bytes igual a 0 o error diferente de nulo: la descarga ha ido mal
			if err != nil || bytes == 0 {
				Error.Println(err)
			}
			//bytes distintos de 0 o error igual a nulo: la descarga se ha realizado correctamente.
			if bytes != 0 || err == nil {
				//Cambiamos el estado del fichero de publicidad en BD, a existe.
				ok, err := db.Prepare("UPDATE mensaje SET existe=? WHERE fichero = ?")
				if err != nil {
					Error.Println(err)
				}
				db_mu.Lock()
				_, err1 := ok.Exec("Y", fichero)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
				}
			}
		}
		time.Sleep(2 * time.Minute) //Cada 2 minutos se revisa en busca de nuevos ficheros (publi/msg)
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
