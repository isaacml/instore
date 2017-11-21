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
	"strconv"
	"sync"
	"time"
)

var (
	Info                 *log.Logger
	Warning              *log.Logger
	Error                *log.Logger
	db                   *sql.DB
	db_mu                sync.RWMutex
	serverint            map[string]string = make(map[string]string) //Mapa que guarda la direccion del servidor interno
	programmedMusic      map[int]string    = make(map[int]string)    //Guarda el listado de carpetas programadas
	copy_arr             []string                                    //Contenedor que va a guardar los ficheros que van a ser copiados a "C:\instore\\Music\"
	capacidad_arr        int                                         //Guarda la capacidad que tiene el array que guarda la ruta de directorio
	username             string                                      //Variable de usuario y estado global
	directorio_actual    string                                      //Va a contener en todo momento la dirección del explorador WIN(handles_publi.go)
	statusProgammedMusic string                                      //Estado de la programacion: Inicial, Actualizada o Modificar
	block       		 bool	                                     //Estado de bloqueo del reproductor y el gestor de descarga de publicidad/mensajes

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
	go estado_de_entidad()
	go solicitudDeFicheros()
	go reproduccion()
	go saveListInBD()
	go reproduccion_msgs()

	// handlers del servidor HTTP
	http.HandleFunc("/", root)
	http.HandleFunc(login_cgi, login)
	http.HandleFunc(logout_cgi, logout)
	// handler de configuracion de tienda
	http.HandleFunc("/get_orgs.cgi", get_orgs)
	http.HandleFunc("/config_shop.cgi", config_shop)
	//Bitmap Actions
	http.HandleFunc("/acciones.cgi", acciones)
	//Exploradores
	http.HandleFunc("/mensajesInstantaneos.cgi", mensajesInstantaneos)
	http.HandleFunc("/explorerMusic.cgi", explorerMusic)
	http.HandleFunc("/programarMusica.cgi", programarMusica)

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
		//Fecha actual
		fecha := libs.MyCurrentDate()
		//Si el fichero de configuracion existe, enviamos dominio/os de la tienda
		if existe == true {
			if block == false {
				var dominios string
				domainint := make(map[string]string) //Mapa que guarda el dominio de la tienda
				loadSettings(configShop, domainint)
				for _, val := range domainint {
					dominios += val + ":.:"
				}
				respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverint["serverinterno"]+"/acciones.cgi", "action;send_domains", "dominios;"+dominios))
				fmt.Println("La respuesta: ", respuesta)
				//Si la respuesta NO está vacía, comprobamos la respuesta.
				if respuesta != "" {
					//De la respuesta obtenemos el listado de mensajes y publicidad
					separar_publi := strings.Split(respuesta, "[publi];")
					//Hay ficheros de publicidad
					if len(separar_publi) > 1 {
						//Se comprueba si el listado contiene mensajes
						tiene_msg := strings.Contains(separar_publi[1], "[mensaje];")
						if tiene_msg != true {
							//SOLO ARCHIVOS DE PUBLICIDAD
							var publi string
							arch_publi := strings.Split(separar_publi[1], ";")
							for _, publi = range arch_publi {
								var cont int
								if strings.Contains(publi, "[mensaje]") {
									publi = strings.TrimRight(publi, "[mensaje]")
								}
								separar := strings.Split(publi, "<=>")
								f_pub := separar[0]
								fecha_ini := separar[1]
								gap := separar[2]
								//Comprobamos si existen los ficheros de publi en la BD interna
								publicidad, errS := db.Query("SELECT * FROM publi WHERE fichero=?", f_pub)
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
									_, err := os.Stat(publi_files_location + f_pub)
									if err != nil {
										//NO lo tiene, se guarda en la BD de player con el estado en N.
										if os.IsNotExist(err) {
											nook, err := db.Prepare("INSERT INTO publi (`fichero`, `existe`, `fecha_ini`, `gap`) VALUES (?,?,?,?)")
											if err != nil {
												Error.Println(err)
											}
											db_mu.Lock()
											_, err1 := nook.Exec(f_pub, "N", fecha_ini, gap)
											db_mu.Unlock()
											if err1 != nil {
												Error.Println(err1)
											}
										}
									} else {
										//SI lo tiene, se guarda en la BD de player con el estado en Y.
										ok, err := db.Prepare("INSERT INTO publi (`fichero`, `existe`, `fecha_ini`, `gap`) VALUES (?,?,?,?)")
										if err != nil {
											Error.Println(err)
										}
										db_mu.Lock()
										_, err1 := ok.Exec(f_pub, "Y", fecha_ini, gap)
										db_mu.Unlock()
										if err1 != nil {
											Error.Println(err1)
										}
									}
								}
							}
						}
						separar_msg := strings.Split(separar_publi[1], "[mensaje];")
						//Hay ficheros de mensaje
						if len(separar_msg) > 1 {
							//Tomamos listados de mensajes, publicidad y los almacenamos
							f_publicidad := strings.Split(separar_msg[0], ";")
							f_mensajes := strings.Split(separar_msg[1], ";")
							//FICHEROS de PUBLICIDAD
							for _, publi := range f_publicidad {
								var cont int
								separar := strings.Split(publi, "<=>")
								f_pub := separar[0]
								fecha_ini := separar[1]
								gap := separar[2]
								//Comprobamos si existen los ficheros de publi en la BD interna
								publicidad, errS := db.Query("SELECT * FROM publi WHERE fichero=?", f_pub)
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
									_, err := os.Stat(publi_files_location + f_pub)
									if err != nil {
										//NO lo tiene, se guarda en la BD de player con el estado en N.
										if os.IsNotExist(err) {
											nook, err := db.Prepare("INSERT INTO publi (`fichero`, `existe`, `fecha_ini`, `gap`) VALUES (?,?,?,?)")
											if err != nil {
												Error.Println(err)
											}
											db_mu.Lock()
											_, err1 := nook.Exec(f_pub, "N", fecha_ini, gap)
											db_mu.Unlock()
											if err1 != nil {
												Error.Println(err1)
											}
										}
									} else {
										//SI lo tiene, se guarda en la BD de player con el estado en Y.
										ok, err := db.Prepare("INSERT INTO publi (`fichero`, `existe`, `fecha_ini`, `gap`) VALUES (?,?,?,?)")
										if err != nil {
											Error.Println(err)
										}
										db_mu.Lock()
										_, err1 := ok.Exec(f_pub, "Y", fecha_ini, gap)
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
											_, err1 := nook.Exec(msgname, playtime, "N", fecha)
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
										_, err1 := ok.Exec(msgname, playtime, "Y", fecha)
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
						separar_mensaje := strings.Split(respuesta, "[mensaje];")
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
											_, err1 := nook.Exec(msgname, playtime, "N", fecha)
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
										_, err1 := ok.Exec(msgname, playtime, "Y", fecha)
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
		}
		time.Sleep(5 * time.Minute)
	}
}

//Se manda hacia el servidor interno una solicitud de los archivos publi/msg que se tiene que bajar.
func solicitudDeFicheros() {
	for {
		//Sacamos la fecha actual
		y, m, d := time.Now().Date()
		fecha := fmt.Sprintf("%4d%02d%02d", y, int(m), d)
		//Busqueda fichero de publicidad y existencia del mismo por fecha actual
		publiQ, err := db.Query("SELECT fichero, existe, fecha_ini, gap FROM publi WHERE fecha_ini=?", fecha)
		if err != nil {
			Error.Println(err)
		}
		for publiQ.Next() {
			var fichero, exist, fecha_ini, gap string
			//Tomamos el nombre del fichero de publicidad y su existencia
			err = publiQ.Scan(&fichero, &exist, &fecha_ini, &gap)
			if err != nil {
				Error.Println(err)
			}
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverint["serverinterno"]+"/publi_msg.cgi", "action;MsgFiles", "fichero;"+fichero, "existencia;"+exist, "fecha_ini;"+fecha_ini, "gap;"+gap))
			//Si en la respuesta obtenemos el valor "Descarga": el player tiene liste el fichero msg para descargarlo
			if respuesta == "Descarga" {
				b, err := libs.DownloadFile(serverint["serverinterno"]+"/"+fichero+"?accion=publicidad", publi_files_location+fichero, 0, 1000)
				//bytes igual a 0 o error diferente de nulo: la descarga ha ido mal
				if err != nil || b == 0 {
					Error.Println(err)
				}
				//bytes distintos de 0 o error igual a nulo: la descarga se ha realizado correctamente.
				if b != 0 || err == nil {
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
		}
		//Busqueda fichero de mensaje y existencia del mismo por fecha actual
		msgQ, err := db.Query("SELECT fichero, existe FROM mensaje WHERE fecha=?", fecha)
		if err != nil {
			Error.Println(err)
		}
		for msgQ.Next() {
			var fichero, exist string
			//Tomamos el nombre del fichero mensaje y su existencia
			err = msgQ.Scan(&fichero, &exist)
			if err != nil {
				Error.Println(err)
			}
			respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverint["serverinterno"]+"/publi_msg.cgi", "action;PubliFiles", "fichero;"+fichero, "existencia;"+exist))
			//Si en la respuesta obtenemos el valor "Descarga": el player tiene liste el fichero msg para descargarlo
			if respuesta == "Descarga" {
				b, err := libs.DownloadFile(serverint["serverinterno"]+"/"+fichero+"?accion=mensaje", msg_files_location+fichero, 0, 1000)
				if err != nil {
					Error.Println(err)
				}
				//bytes igual a 0: la descarga ha ido mal
				if b == 0 {
					Error.Println("Size Zero: NO se ha descargado el fichero")
				} else { //la descarga se ha realizado correctamente.
					//Cambiamos el estado del fichero de mensaje en BD, a existe.
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
		}
		time.Sleep(1 * time.Minute)
	}
}

//Toma el estado de entidad y se desarrolla el proceso de bloqueo de la tienda
func estado_de_entidad() {
	for {
		//Se obtiene el dominio completo
		dom := libs.MainDomain(configShop)
		//Nos quedamos con el ent[0] que contiende el nombre de la entidad.
		ent := strings.Split(dom, ".")
		res := libs.GenerateFORM(serverint["serverinterno"]+"/acciones.cgi", "action;check_entidad", "ent;"+ent[0])
		db_mu.Lock()
		//Guarda el estado de activo de la entidad: 0 - OFF / 1 - ON
		estado_entidad, _ := strconv.Atoi(res)
		db_mu.Unlock()
		if estado_entidad == 1 { //ON
			var last_connect int64
			var seg_del_mes int64
			year, mes, _ := time.Now().Date()
			timestamp := time.Now().Unix()
			dias_del_mes := libs.DaysIn(mes, year)
			seg_del_mes = dias_del_mes * 4
			//Tomamos la ultima conexion de la tienda
			db.QueryRow("SELECT last_connect FROM tienda WHERE dominio=?", dom).Scan(&last_connect)
			if last_connect-(timestamp-seg_del_mes) < 0 {
				block = true
			}else{
				block = false
			}
		}else{
			block = true
		}
		time.Sleep(5 * time.Second)
	}
}

/*
loadSettings: esta función va a abrir un fichero, leer los datos que contiene y guardarlos en un mapa.
	filename: ruta completa donde se encuentra nuestro fichero(C:\instore\serverext.reg)
*/
func loadSettings(filename string, arr map[string]string) {
	cont := 1
	fr, err := os.Open(filename)
	defer fr.Close()
	if err == nil {
		reader := bufio.NewReader(fr)
		for {
			linea, rerr := reader.ReadString('\n')
			if rerr != nil {
				break
			}
			linea = strings.TrimRight(linea, "\r\n")
			item := strings.Split(linea, " = ")
			if len(item) == 2 {
				if _, ok := arr[item[0]]; ok {
					clave := fmt.Sprintf("%s%d", item[0], cont)
					arr[clave] = item[1]
					cont++
				} else {
					arr[item[0]] = item[1]
				}
			}
		}
	}
}
