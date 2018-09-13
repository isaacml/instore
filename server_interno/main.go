package main

import (
	"database/sql"
	"fmt"
	"github.com/isaacml/instore/libs"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	Info      *log.Logger
	Warning   *log.Logger
	Error     *log.Logger
	db        *sql.DB
	db_mu     sync.Mutex
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
	db, err_db = sql.Open("sqlite3", sql_file)
	if err_db != nil {
		Error.Println(err_db)
		log.Fatalln("Fallo al abrir el archivo de error:", err_db)
	}
	db.Exec("PRAGMA journal_mode=WAL;")
	libs.LoadSettingsLin(serverRoot, serverext) // Se carga los valores del fichero serverint.reg
}

// funcion principal del programa
func main() {
	fmt.Printf("Golang HTTP Server starting at Port %s ...\n", serverext["puerto_interno"])
	go BuscarNuevosFicheros()
	go BorrarFicherosAntiguos()
	//servidor de ficheros
	http.HandleFunc("/", root)
	//Actions del servidor externo
	http.HandleFunc("/auth.cgi", auth)
	http.HandleFunc("/acciones.cgi", acciones)
	http.HandleFunc("/transf_orgs.cgi", transf_orgs)
	http.HandleFunc("/publi_msg.cgi", publi_msg)

	s := &http.Server{
		Addr:           ":" + serverext["puerto_interno"],
		Handler:        nil,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 13,
	}

	log.Fatal(s.ListenAndServe()) // servidor HTTP multihilo
}

//Esta funcion esta consultando al servidor externo cada cierto tiempo en busca de ficheros publi/msg que bajarse
func BuscarNuevosFicheros() {
	for {
		//Buscamos todos lo ficheros de publicidad que no tenemos en la BD
		db_mu.Lock()
		publicidad, errP := db.Query("SELECT fichero FROM publi WHERE existe='N'")
		db_mu.Unlock()
		if errP != nil {
			Error.Println(errP)
			return
		}
		for publicidad.Next() {
			var fichero string
			//Tomamos el nombre del fichero publicidad
			err := publicidad.Scan(&fichero)
			if err != nil {
				Error.Println(err)
				continue
			}
			//Descargamos el fichero del servidor externo
			bytes, err := libs.DownloadFile(serverext["serverexterno"]+"/"+fichero+"?accion=publicidad", publi_files_location+fichero, 0, 1000)
			//bytes igual a 0 o error diferente de nulo: la descarga ha ido mal
			if err != nil || bytes == 0 {
				Error.Println(err)
				continue
			}
			//bytes distintos de 0 o error igual a nulo: la descarga se ha realizado correctamente.
			if bytes != 0 || err == nil {
				//Cambiamos el estado del fichero de publicidad en BD, a existe.
				ok, err := db.Prepare("UPDATE publi SET existe=? WHERE fichero = ?")
				if err != nil {
					Error.Println(err)
					continue
				}
				db_mu.Lock()
				_, err1 := ok.Exec("Y", fichero)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
					continue
				}
			}
		}
		//Buscamos todos lo ficheros de mensajes que no tenemos en la BD
		db_mu.Lock()
		mensajes, errM := db.Query("SELECT fichero FROM mensaje WHERE existe='N'")
		db_mu.Unlock()
		if errM != nil {
			Error.Println(errM)
			return
		}
		for mensajes.Next() {
			var fichero string
			//Tomamos el nombre del fichero mensaje
			err := mensajes.Scan(&fichero)
			if err != nil {
				Error.Println(err)
				continue
			}
			//Descargamos el fichero del servidor externo
			bytes, err := libs.DownloadFile(serverext["serverexterno"]+"/"+fichero+"?accion=mensaje", msg_files_location+fichero, 0, 1000)
			//bytes igual a 0 o error diferente de nulo: la descarga ha ido mal
			if err != nil || bytes == 0 {
				Error.Println(err)
				continue
			}
			//bytes distintos de 0 o error igual a nulo: la descarga se ha realizado correctamente.
			if bytes != 0 || err == nil {
				//Cambiamos el estado del fichero de publicidad en BD, a existe.
				ok, err := db.Prepare("UPDATE mensaje SET existe=? WHERE fichero = ?")
				if err != nil {
					Error.Println(err)
					continue
				}
				db_mu.Lock()
				_, err1 := ok.Exec("Y", fichero)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
					continue
				}
			}
		}
		time.Sleep(2 * time.Minute) //Cada 2 minutos se revisa en busca de nuevos ficheros (publi/msg)
	}
}

//Borrado de ficheros de publicidad y mensajes
func BorrarFicherosAntiguos() {
	var del_err error
	for {
		//tiempo limite = 1 mes
		limit_time := time.Now().Unix() - 2592000
		//PUBLICIDAD
		db_mu.Lock()
		publi, errP := db.Query("SELECT fichero FROM publi WHERE timestamp < ?", limit_time)
		db_mu.Unlock()
		if errP != nil {
			Error.Println(errP)
			continue
		}
		for publi.Next() {
			var fichero string
			//Tomamos el nombre del fichero mensaje
			err := publi.Scan(&fichero)
			if err != nil {
				Error.Println(err)
				continue
			}
			//Borramos el fichero desde la ruta interna
			err = os.Remove(publi_files_location + fichero)
			if err != nil {
				Error.Println(err)
				continue
			}
		}
		//Borramos de la base de datos los ficheros de publicidad
		db_mu.Lock()
		_, del_err = db.Exec("DELETE FROM publi WHERE timestamp < ?", limit_time)
		db_mu.Unlock()
		if del_err != nil {
			Error.Println(del_err)
			continue
		}
		//MENSAJES
		msg, errM := db.Query("SELECT fichero FROM mensaje WHERE timestamp < ?", limit_time)
		if errM != nil {
			Error.Println(errM)
			continue
		}
		for msg.Next() {
			var fichero string
			//Tomamos el nombre del fichero mensaje
			err := msg.Scan(&fichero)
			if err != nil {
				Error.Println(err)
				continue
			}
			//Borramos el fichero desde la ruta interna
			err = os.Remove(msg_files_location + fichero)
			if err != nil {
				Error.Println(err)
				continue
			}
		}
		//Borramos de la base de datos los ficheros de mensajes
		db_mu.Lock()
		_, del_err = db.Exec("DELETE FROM mensaje WHERE timestamp < ?", limit_time)
		db_mu.Unlock()
		if del_err != nil {
			Error.Println(del_err)
			continue
		}
		time.Sleep(5 * time.Minute) //Cada 5 minutos se revisa en busca de nuevos ficheros (publi/msg) para borrar
	}
}
