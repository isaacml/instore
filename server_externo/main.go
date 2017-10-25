package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	Info          *log.Logger
	Warning       *log.Logger
	Error         *log.Logger
	db            *sql.DB
	db_mu         sync.RWMutex
	bad, empty    string //Variables de estado global
	status_dom    string //Variable que va a guardar el dominio de la tienda
	enviar_estado bool   //Variable de estado para saber si podemos guardar el dominio de la tienda o NO
)

//MASCARAS PARA BITMAP
const (
	PROG_PUB = 1 << iota
	PROG_MUS
	PROG_MSG
	ADD_MUS
	MSG_AUTO
	MSG_NORMAL
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
	db, err_db = sql.Open("sqlite3", "/home/isaac/INSTORE_SQL/pruebas/servext.db")
	if err_db != nil {
		Error.Println(err_db)
		log.Fatalln("Fallo al abrir el archivo de error:", err_db)
	}
	db.Exec("PRAGMA journal_mode=WAL;")
}

// funcion principal del programa
func main() {
	fmt.Printf("Golang HTTP Server starting at Port %s ...\n", http_port)
	go BorrarFicherosAntiguos()
	// handlers del servidor
	http.HandleFunc("/login.cgi", login)
	http.HandleFunc("/", root)
	//ACCIONES
	http.HandleFunc("/acciones.cgi", acciones)
	//USUARIOS
	http.HandleFunc("/usuarios.cgi", usuarios)
	//ENTIDADES
	http.HandleFunc("/entidades.cgi", entidades)
	//ALMACENES
	http.HandleFunc("/almacenes.cgi", almacenes)
	//PAISES
	http.HandleFunc("/paises.cgi", paises)
	//REGIONES
	http.HandleFunc("/regiones.cgi", regiones)
	//PROVINCIAS
	http.HandleFunc("/provincias.cgi", provincias)
	//TIENDAS
	http.HandleFunc("/tiendas.cgi", tiendas)
	//RECOGER LOS FICHEROS
	http.HandleFunc("/publi_files.cgi", publi_files)
	http.HandleFunc("/msg_files.cgi", msg_files)
	//SELECTS PARA EL PLAYER INTERNO
	http.HandleFunc("/config_shop.cgi", config_shop)
	http.HandleFunc("/send_shop.cgi", send_shop)
	http.HandleFunc("/recoger_dominio.cgi", recoger_dominio)

	http.HandleFunc("/organizaciones.cgi", organizaciones)

	s := &http.Server{
		Addr:           ":" + http_port,
		Handler:        nil,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 13,
	}

	log.Fatal(s.ListenAndServe()) // servidor HTTP multihilo
}

//Borrado de ficheros de publicidad y mensajes
func BorrarFicherosAntiguos() {
	for {
		//fecha de ahora
		yy, mm, dd := time.Now().Date()
		now := fmt.Sprintf("%4d%02d%2d", yy, int(mm), dd)
		//tiempo limite = 1 mes 2592000
		limit_time := time.Now().Unix() - 86400
		//PUBLICIDAD
		publi, errP := db.Query("SELECT id, fichero FROM publi WHERE fecha_final < ? AND timestamp < ? ", now, limit_time)
		if errP != nil {
			Error.Println(errP)
		}
		for publi.Next() {
			var id int
			var fichero string
			//Tomamos el nombre del fichero publicidad
			err := publi.Scan(&id, &fichero)
			if err != nil {
				Error.Println(err)
			}
			//Borramos id y fichero desde la ruta interna para el borrado
			err = os.Remove(publi_files_location + fichero)
			if err != nil {
				Error.Println(err)
			}
			//Borramos de la base de datos los ficheros de publicidad
			db_mu.Lock()
			db.Exec("DELETE FROM publi WHERE id = ?", id)
			db_mu.Unlock()
		}
		//MENSAJES
		msg, errM := db.Query("SELECT id, fichero FROM mensaje WHERE fecha_final < ? AND timestamp < ? ", now, limit_time)
		if errM != nil {
			Error.Println(errM)
		}
		for msg.Next() {
			var id int
			var fichero string
			//Tomamos id y nombre del fichero mensaje para el borrado
			err := msg.Scan(&id, &fichero)
			if err != nil {
				Error.Println(err)
			}
			//Borramos el fichero desde la ruta interna
			err = os.Remove(msg_files_location + fichero)
			if err != nil {
				Error.Println(err)
			}
			//Borramos de la base de datos los ficheros de mensajes
			db_mu.Lock()
			db.Exec("DELETE FROM mensaje WHERE id = ?", id)
			db_mu.Unlock()
		}
		time.Sleep(2 * time.Minute) //Cada 2 minutos se revisa en busca de nuevos ficheros (publi/msg) para borrar
	}
}
