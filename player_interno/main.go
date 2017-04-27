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
	db, err_db = sql.Open("sqlite3", "C:\\instore\\shop.db") //WINDB: C:\\instore\\shop.db
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
	go saveListInBD()
	go bajadoDeFicheros()

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

//Comprueba cada 5 min listados nuevos y los guarda en la base de datos interna
func saveListInBD() {
	for {
		//Comprobar que existe el fichero de configuracion de la tienda
		var existe bool
		_, err := os.Stat(configShop)
		if err != nil {
			if os.IsNotExist(err) {
				existe = false
			}
		} else {
			existe = true
		}
		fecha_actual := time.Now()
		//Formato de la fecha actual --> 20070405
		string_fecha := fmt.Sprintf("%4d%02d%02d", fecha_actual.Year(), int(fecha_actual.Month()), fecha_actual.Day())
		//Si el fichero de configuracion existe, enviamos el dominio de la tienda
		if existe == true {
			loadSettings(configShop, domainint)
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverint["serverinterno"]+"/send_domain.cgi", "dominio;"+domainint["shopdomain"]))
			//Si la respuesta NO está vacía, comprobamos la respuesta.
			if respuesta != "" {
				//De la respuesta obtenemos el listado de mensajes y publicidad
				separar_publi := strings.Split(respuesta, "[publi]")
				//Hay ficheros de publicidad
				if len(separar_publi) > 1 {
					//Se comprueba si el listado contiene mensajes
					tiene_msg := strings.Contains(separar_publi[1], "[mensaje]")
					if tiene_msg != true {
						//SOLO ARCHIVOS DE PUBLICIDAD
						arch_publi := strings.Split(separar_publi[1], ";")
						for _, publi := range arch_publi {
							var cont int
							//Comprobamos si existen los ficheros de publi en la BD interna
							publicidad, errS := db.Query("SELECT * FROM publi WHERE fichero=?", publi)
							if errS != nil {
								Error.Println(errS)
							}
							//Si existe, el contador incrementará
							for publicidad.Next() {
								cont++
							}
							//Contador = 0 --> La BD interna no tiene el fichero publi
							if cont == 0 {
								//Se comprueba si el player_interno tiene el fichero publi.
								_, err := os.Stat(publi_files_location + publi)
								if err != nil {
									//NO lo tiene, se guarda en la BD de player con el estado en N.
									if os.IsNotExist(err) {
										nook, err := db.Prepare("INSERT INTO publi (`fichero`, `existe`, `fecha`) VALUES (?,?,?)")
										if err != nil {
											Error.Println(err)
										}
										db_mu.Lock()
										_, err1 := nook.Exec(publi, "N", string_fecha)
										db_mu.Unlock()
										if err1 != nil {
											Error.Println(err1)
										}
									}
								} else {
									//SI lo tiene, se guarda en la BD de player con el estado en Y.
									ok, err := db.Prepare("INSERT INTO publi (`fichero`, `existe`, `fecha`) VALUES (?,?,?)")
									if err != nil {
										Error.Println(err)
									}
									db_mu.Lock()
									_, err1 := ok.Exec(publi, "Y", string_fecha)
									db_mu.Unlock()
									if err1 != nil {
										Error.Println(err1)
									}
								}
							}
						}
					}
					separar_msg := strings.Split(separar_publi[1], "[mensaje]")
					//Hay ficheros de mensaje
					if len(separar_msg) > 1 {
						//Tomamos listados de mensajes, publicidad y los almacenamos
						f_publicidad := strings.Split(separar_msg[0], ";")
						f_mensajes := strings.Split(separar_msg[1], ";")
						//FICHEROS de PUBLICIDAD
						for _, publi := range f_publicidad {
							var cont int
							//Comprobamos si existen los ficheros de publi en la BD interna
							publicidad, errS := db.Query("SELECT * FROM publi WHERE fichero=?", publi)
							if errS != nil {
								Error.Println(errS)
							}
							//Si existe, el contador incrementará
							for publicidad.Next() {
								cont++
							}
							//Contador = 0 --> La BD interna no tiene el fichero publi
							if cont == 0 {
								//Se comprueba si el player_interno tiene el fichero publi.
								_, err := os.Stat(publi_files_location + publi)
								if err != nil {
									//NO lo tiene, se guarda en la BD de player con el estado en N.
									if os.IsNotExist(err) {
										nook, err := db.Prepare("INSERT INTO publi (`fichero`, `existe`, `fecha`) VALUES (?,?,?)")
										if err != nil {
											Error.Println(err)
										}
										db_mu.Lock()
										_, err1 := nook.Exec(publi, "N", string_fecha)
										db_mu.Unlock()
										if err1 != nil {
											Error.Println(err1)
										}
									}
								} else {
									//SI lo tiene, se guarda en la BD de player con el estado en Y.
									ok, err := db.Prepare("INSERT INTO publi (`fichero`, `existe`, `fecha`) VALUES (?,?,?)")
									if err != nil {
										Error.Println(err)
									}
									db_mu.Lock()
									_, err1 := ok.Exec(publi, "Y", string_fecha)
									db_mu.Unlock()
									if err1 != nil {
										Error.Println(err1)
									}
								}
							}
						}
						//FICHEROS de MENSAJES
						for _, msg := range f_mensajes {
							var cont int
							//Separamos entre nombre y playtime de los mensajes
							separar := strings.Split(msg, "<=>")
							msgname := separar[0]
							playtime := separar[1]
							//Comprobamos si existen los mensajes en la BD interna
							mensajes, errS := db.Query("SELECT * FROM mensaje WHERE fichero=?", msgname)
							if errS != nil {
								Error.Println(errS)
							}
							//Si existe, el contador incrementará
							for mensajes.Next() {
								cont++
							}
							//Contador = 0 --> La BD interna no tiene el mensaje
							if cont == 0 {
								//Se comprueba si el player_interno tiene el fichero mensaje.
								_, err := os.Stat(msg_files_location + msgname)
								if err != nil {
									//NO lo tiene, se guarda en la BD de player con el estado en N.
									if os.IsNotExist(err) {
										nook, err := db.Prepare("INSERT INTO mensaje (`fichero`, `playtime`, `existe`, `fecha`) VALUES (?,?,?,?)")
										if err != nil {
											Error.Println(err)
										}
										db_mu.Lock()
										_, err1 := nook.Exec(msgname, playtime, "N", string_fecha)
										db_mu.Unlock()
										if err1 != nil {
											Error.Println(err1)
										}
									}
								} else {
									//SI lo tiene, se guarda en la BD de player con el estado en Y.
									ok, err := db.Prepare("INSERT INTO mensaje (`fichero`, `playtime`, `existe`, `fecha`) VALUES (?,?,?,?)")
									if err != nil {
										Error.Println(err)
									}
									db_mu.Lock()
									_, err1 := ok.Exec(msgname, playtime, "Y", string_fecha)
									db_mu.Unlock()
									if err1 != nil {
										Error.Println(err1)
									}
								}
							}
						}
					}
				} else {
					//No hay ficheros de publicidad, por tanto vamos a comprobar si hay mensajes
					separar_mensaje := strings.Split(respuesta, "[mensaje]")
					if len(separar_mensaje) > 1 {
						//Hay mensajes, vamos a obtenerlos uno a uno
						mensajes := strings.Split(separar_mensaje[1], ";")
						for _, msg := range mensajes {
							var cont int
							//Separamos entre nombre de mensaje y playtime del mensaje
							separar := strings.Split(msg, "<=>")
							msgname := separar[0]
							playtime := separar[1]
							mensajes, errS := db.Query("SELECT * FROM mensaje WHERE fichero=?", msgname)
							if errS != nil {
								Error.Println(errS)
							}
							//Si el mensaje existe, el contador se incrementará
							for mensajes.Next() {
								cont++
							}
							//contador = 0 --> no existe el mensaje en BD, por lo tanto vamos a añadirlo.
							if cont == 0 {
								//Se comprueba si el player_interno tiene el fichero mensaje.
								_, err := os.Stat(msg_files_location + msgname)
								//NO lo tiene, se guarda en la BD de player con el estado en N.
								if err != nil {
									if os.IsNotExist(err) {
										nook, err := db.Prepare("INSERT INTO mensaje (`fichero`, `playtime`, `existe`, `fecha`) VALUES (?,?,?,?)")
										if err != nil {
											Error.Println(err)
										}
										db_mu.Lock()
										_, err1 := nook.Exec(msgname, playtime, "N", string_fecha)
										db_mu.Unlock()
										if err1 != nil {
											Error.Println(err1)
										}
									}
								} else {
									//SI lo tiene, se guarda en la BD de player con el estado en Y.
									ok, err := db.Prepare("INSERT INTO mensaje (`fichero`, `playtime`, `existe`, `fecha`) VALUES (?,?,?,?)")
									if err != nil {
										Error.Println(err)
									}
									db_mu.Lock()
									_, err1 := ok.Exec(msgname, playtime, "Y", string_fecha)
									db_mu.Unlock()
									if err1 != nil {
										Error.Println(err1)
									}
								}
							}
						}
					}
				}
			}
		}
		time.Sleep(5 * time.Minute)
	}
}

func bajadoDeFicheros() {
	for {
		//Sacamos la fecha actual
		y, m, d := time.Now().Date()
		fecha := fmt.Sprintf("%4d%02d%02d", y, int(m), d)
		//Busqueda fichero de publicidad y existencia del mismo por fecha actual
		publiQ, err := db.Query("SELECT fichero, existe FROM publi WHERE fecha=?", fecha)
		if err != nil {
			Error.Println(err)
		}
		for publiQ.Next() {
			var fichero, exist string
			//Tomamos el nombre del fichero de publicidad y su existencia
			err = publiQ.Scan(&fichero, &exist)
			if err != nil {
				Error.Println(err)
			}
			fmt.Sprintf("%s", libs.GenerateFORM(serverint["serverinterno"]+"/downloadPubliFile.cgi", "fichero;"+fichero, "existencia;"+exist))
			fmt.Println(fichero, exist)
		}
		//Busqueda fichero de mensaje y existencia del mismo por fecha actual
		msgQ, err := db.Query("SELECT fichero, existe FROM mensaje WHERE fecha=?", fecha)
		if err != nil {
			Error.Println(err)
		}
		for msgQ.Next() {
			var fichero, exist string
			//Tomamos el nombre del fichero mensaje y su existencia
			err = publiQ.Scan(&fichero, &exist)
			if err != nil {
				Error.Println(err)
			}
			fmt.Sprintf("%s", libs.GenerateFORM(serverint["serverinterno"]+"/downloadMsgFile.cgi", "fichero;"+fichero, "existencia;"+exist))
			fmt.Println(fichero, exist)
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
