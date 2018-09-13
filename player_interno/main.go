package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/isaacml/instore/libs"
	"github.com/isaacml/instore/winamp"
	_ "github.com/mattn/go-sqlite3"
)

var (
	Info          *log.Logger
	Warning       *log.Logger
	Error         *log.Logger
	db            *sql.DB
	db_mu         sync.Mutex
	settings      map[string]string = make(map[string]string) //Guarda los settings de la tienda
	capacidad_arr int                                         //Guarda la capacidad que tiene el array que guarda la ruta de directorio
	block         bool                                        //Estado de bloqueo del reproductor y el gestor de descarga de publicidad/mensajes
	schedule      bool                                        //Guarda el estado que genera el horario de reproducción (true: reproduce | false: no reproduce)
	winplayer     *winamp.Winamp
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
	db, err_db = sql.Open("sqlite3", bd_name)
	if err_db != nil {
		Error.Println(err_db)
		log.Fatalln("Fallo al abrir el archivo de error:", err_db)
	}
	db.Exec("PRAGMA journal_mode=WAL;")
	libs.LoadSettingsLin(serverRoot, settings) // Se carga los valores del fichero SettingsShop.reg
	winplayer = winamp.Winamper()
}

// Funcion principal del programa
func main() {
	fmt.Printf("Golang HTTP Server starting at Port %s ...\n", settings["port"])
	go controlinternalsessions() // Controla la caducidad de la sesion
	go estado_de_entidad()
	go horario_reproduccion()
	go solicitudDeFicheros()
	go saveListInBD()
	go reproduccion()
	go reproduccion_msgs()
	go borrarPublicidad()

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
	http.HandleFunc("/instantaneos.cgi", instantaneos)
	http.HandleFunc("/mostrar_boton.cgi", mostrar_boton)
	http.HandleFunc("/playInstantaneos.cgi", playInstantaneos)
	http.HandleFunc("/programarMusica.cgi", programarMusica)
	http.HandleFunc("/volumen_global.cgi", volumen_global)

	s := &http.Server{
		Addr:           ":" + settings["port"],
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
		existe := libs.Existencia(configShop)
		//Si el fichero de configuracion existe, enviamos dominio/os de la tienda
		if existe == true {
			if block == false {
				var dominios string
				domainint := make(map[string]string) //Mapa que guarda el dominio de la tienda
				libs.LoadDomains(configShop, domainint)
				for _, val := range domainint {
					dominios += val + ":.:"
				}
				respuesta := fmt.Sprintf("%s", libs.GenerateFORM(settings["serverinterno"]+"/acciones.cgi", "action;send_domains", "dominios;"+dominios))
				//fmt.Println("La respuesta: ", respuesta)
				//Si la respuesta NO está vacía, comprobamos la respuesta.
				if respuesta != "" {
					//De la respuesta obtenemos el listado de mensajes y publicidad
					separar_publi := strings.Split(respuesta, "[publi];")
					if len(separar_publi) > 1 { //Hay ficheros de publicidad
						tiene_msg := strings.Contains(separar_publi[1], "[mensaje];")
						if tiene_msg != true { //Se comprueba si el listado contiene mensajes
							//SOLO ARCHIVOS DE PUBLICIDAD
							var publi string
							borrar_todos_mensajes()
							borrar_publi_int(separar_publi[1])
							arch_publi := strings.Split(separar_publi[1], ";")
							for _, publi = range arch_publi {
								var cont int
								var bd_gap, bd_f_ini, bd_f_fin string
								if strings.Contains(publi, "[mensaje]") {
									publi = strings.TrimRight(publi, "[mensaje]")
								}
								separar := strings.Split(publi, "<=>")
								f_pub := separar[0]
								fecha_ini := separar[1]
								fecha_fin := separar[2]
								gap := separar[3]
								//Comprobamos si existen los ficheros de publi en la BD interna
								db.QueryRow("SELECT count(id), fecha_ini, fecha_fin, gap FROM publi WHERE fichero=?", f_pub).Scan(&cont, &bd_f_ini, &bd_f_fin, &bd_gap)
								//Contador = 0 --> La BD interna no tiene el fichero publi
								if cont == 0 {
									//Se comprueba si el player_interno tiene el fichero publi.
									_, err := os.Stat(publi_files_location + f_pub)
									if err != nil {
										//NO lo tiene, se guarda en la BD de player con el estado en N.
										if os.IsNotExist(err) {
											insert_publi(f_pub, "N", fecha_ini, fecha_fin, gap)
										}
									} else {
										//SI lo tiene, se guarda en la BD de player con el estado en Y.
										insert_publi(f_pub, "Y", fecha_ini, fecha_fin, gap)
									}
								} else {
									if bd_gap != gap || bd_f_ini != fecha_ini || bd_f_fin != fecha_fin {
										update_publi(f_pub, fecha_ini, fecha_fin, gap)
									}
								}
							}
						}
						separar_msg := strings.Split(separar_publi[1], "[mensaje];")
						//Hay ficheros de mensaje
						if len(separar_msg) > 1 {
							borrar_publi_int(separar_msg[0])
							borrar_mensajes_int(separar_msg[1])
							//Tomamos listados de mensajes, publicidad y los almacenamos
							f_publicidad := strings.Split(separar_msg[0], ";")
							f_mensajes := strings.Split(separar_msg[1], ";")
							//FICHEROS de PUBLICIDAD
							for _, publi := range f_publicidad {
								var cont int
								var bd_gap, bd_f_ini, bd_f_fin string
								separar := strings.Split(publi, "<=>")
								f_pub := separar[0]
								fecha_ini := separar[1]
								fecha_fin := separar[2]
								gap := separar[3]
								//Comprobamos si existen los ficheros de publi en la BD interna
								db.QueryRow("SELECT count(id), fecha_ini, fecha_fin, gap FROM publi WHERE fichero=?", f_pub).Scan(&cont, &bd_f_ini, &bd_f_fin, &bd_gap)
								//Contador = 0 --> La BD interna no tiene el fichero publi
								if cont == 0 {
									//Se comprueba si el player_interno tiene el fichero publi.
									_, err := os.Stat(publi_files_location + f_pub)
									if err != nil {
										//NO lo tiene, se guarda en la BD de player con el estado en N.
										if os.IsNotExist(err) {
											insert_publi(f_pub, "N", fecha_ini, fecha_fin, gap)
										}
									} else {
										//SI lo tiene, se guarda en la BD de player con el estado en Y.
										insert_publi(f_pub, "Y", fecha_ini, fecha_fin, gap)
									}
								} else {
									if bd_gap != gap || bd_f_ini != fecha_ini || bd_f_fin != fecha_fin {
										update_publi(f_pub, fecha_ini, fecha_fin, gap)
									}
								}	
							}
							//FICHEROS de MENSAJES
							for _, msg := range f_mensajes {
								var cont int
								var bd_playtime, bd_f_ini, bd_f_fin string
								//Separamos entre nombre y playtime de los mensajes
								separar := strings.Split(msg, "<=>")
								msgname := separar[0]
								fecha_ini := separar[1]
								fecha_fin := separar[2]
								playtime := separar[3]
								//Comprobamos si existen los mensajes en la BD interna
								db.QueryRow("SELECT count(id), fecha_ini, fecha_fin, playtime FROM mensaje WHERE fichero=?", msgname).Scan(&cont, &bd_f_ini, &bd_f_fin, &bd_playtime)
								//Contador = 0 --> La BD interna no tiene el mensaje
								if cont == 0 {
									//Se comprueba si el player_interno tiene el fichero mensaje.
									_, err := os.Stat(msg_files_location + msgname)
									if err != nil {
										//NO lo tiene, se guarda en la BD de player con el estado en N.
										if os.IsNotExist(err) {
											insert_msg(msgname, "N", fecha_ini, fecha_fin, playtime)
										}
									} else {
										//SI lo tiene, se guarda en la BD de player con el estado en Y.
										insert_msg(msgname, "Y", fecha_ini, fecha_fin, playtime)
									}
								} else {
									if bd_playtime != playtime || bd_f_ini != fecha_ini || bd_f_fin != fecha_fin {
										update_msg(msgname, fecha_ini, fecha_fin, playtime)
									}
								}
							}
						}
					} else {
						//No hay ficheros de publicidad, por tanto vamos a comprobar si hay mensajes
						separar_mensaje := strings.Split(respuesta, "[mensaje];")
						if len(separar_mensaje) > 1 {
							borrar_toda_publi()
							borrar_mensajes_int(separar_mensaje[1])
							//Hay mensajes, vamos a obtenerlos uno a uno
							mensajes := strings.Split(separar_mensaje[1], ";")
							for _, msg := range mensajes {
								var cont int
								var bd_playtime, bd_f_ini, bd_f_fin string
								//Separamos entre nombre de mensaje y playtime del mensaje
								separar := strings.Split(msg, "<=>")
								msgname := separar[0]
								fecha_ini := separar[1]
								fecha_fin := separar[2]
								playtime := separar[3]
								db.QueryRow("SELECT count(id), fecha_ini, fecha_fin, playtime FROM mensaje WHERE fichero=?", msgname).Scan(&cont, &bd_f_ini, &bd_f_fin, &bd_playtime)
								//contador = 0 --> no existe el mensaje en BD, por lo tanto vamos a añadirlo.
								if cont == 0 {
									//Se comprueba si el player_interno tiene el fichero mensaje.
									_, err := os.Stat(msg_files_location + msgname)
									//NO lo tiene, se guarda en la BD de player con el estado en N.
									if err != nil {
										if os.IsNotExist(err) {
											insert_msg(msgname, "N", fecha_ini, fecha_fin, playtime)
										}
									} else {
										//SI lo tiene, se guarda en la BD de player con el estado en Y.
										insert_msg(msgname, "Y", fecha_ini, fecha_fin, playtime)
									}
								} else {
									if bd_playtime != playtime || bd_f_ini != fecha_ini || bd_f_fin != fecha_fin {
										update_msg(msgname, fecha_ini, fecha_fin, playtime)
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
		//Solo bajamos ficheros si la tienda está desbloqueada (block = false)
		if block == false {
			//Sacamos la fecha actual
			fecha := libs.MyCurrentDate()
			//Busqueda fichero de publicidad y existencia del mismo por fecha actual
			publiQ, err := db.Query("SELECT fichero, existe, fecha_ini, fecha_fin, gap FROM publi")
			if err != nil {
				Error.Println(err)
			}
			for publiQ.Next() {
				var fichero, exist, fecha_ini, fecha_fin, gap string
				//Tomamos el nombre del fichero de publicidad y su existencia
				err = publiQ.Scan(&fichero, &exist, &fecha_ini, &fecha_fin, &gap)
				if err != nil {
					Error.Println(err)
				}
				//BETWEEN
				if fecha_ini <= fecha && fecha_fin >= fecha {
					respuesta := fmt.Sprintf("%s", libs.GenerateFORM(settings["serverinterno"]+"/publi_msg.cgi", "action;PubliFiles", "fichero;"+fichero, "existencia;"+exist, "fecha_ini;"+fecha_ini, "fecha_fin;"+fecha_fin, "gap;"+gap))
					//Si en la respuesta obtenemos el valor "Descarga": el player tiene liste el fichero msg para descargarlo
					if respuesta == "Descarga" {
						b, err := libs.DownloadFile(settings["serverinterno"]+"/"+fichero+"?accion=publicidad", publi_files_location+fichero, 0, 1000)
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
			}
			//Busqueda fichero de mensaje y existencia del mismo por fecha actual
			msgQ, err := db.Query("SELECT fichero, existe, fecha_ini, fecha_fin, playtime FROM mensaje")
			if err != nil {
				Error.Println(err)
			}
			for msgQ.Next() {
				var fichero, exist, fecha_ini, fecha_fin, playtime string
				//Tomamos el nombre del fichero mensaje y su existencia
				err = msgQ.Scan(&fichero, &exist, &fecha_ini, &fecha_fin, &playtime)
				if err != nil {
					Error.Println(err)
				}
				//BETWEEN
				if fecha_ini <= fecha && fecha_fin >= fecha {
					respuesta := fmt.Sprintf("%s", libs.GenerateFORM(settings["serverinterno"]+"/publi_msg.cgi", "action;MsgFiles", "fichero;"+fichero, "existencia;"+exist, "fecha_ini;"+fecha_ini, "fecha_fin;"+fecha_fin, "playtime;"+playtime))
					//Si en la respuesta obtenemos el valor "Descarga": el player tiene liste el fichero msg para descargarlo
					if respuesta == "Descarga" {
						b, err := libs.DownloadFile(settings["serverinterno"]+"/"+fichero+"?accion=mensaje", msg_files_location+fichero, 0, 1000)
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
		res := libs.GenerateFORM(settings["serverinterno"]+"/acciones.cgi", "action;check_entidad", "ent;"+ent[0])
		if res == "" {
			var last_connect int64
			var seg_del_mes int64
			year, mes, _ := time.Now().Date()
			timestamp := time.Now().Unix()
			dias_del_mes := libs.DaysIn(mes, year)
			//total de dias * segundos que tiene un dia
			seg_del_mes = dias_del_mes * segs_of_day
			db.QueryRow("SELECT last_connect FROM tienda WHERE dominio=?", dom).Scan(&last_connect)
			db_mu.Lock()
			if last_connect-(timestamp-seg_del_mes) < 0 {
				block = true
			} else {
				block = false
			}
			db_mu.Unlock()
		} else {
			//Guarda el estado de activo de la entidad: 0 - OFF / 1 - ON
			estado_entidad, _ := strconv.Atoi(res)
			//Dominio principal de la tienda guardado en el fichero de configuracion
			domain := libs.MainDomain(configShop)
			if estado_entidad == 1 {
				timestamp := time.Now().Unix()
				db_mu.Lock()
				_, err1 := db.Exec("UPDATE tienda SET last_connect = ? WHERE dominio = ?", timestamp, domain)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
				}
				block = false
			}
			if estado_entidad == 0 {
				db_mu.Lock()
				_, err1 := db.Exec("UPDATE tienda SET last_connect = ? WHERE dominio = ?", 1000, domain)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
				}
				block = true
			}
		}
		time.Sleep(5 * time.Minute)
	}
}

//Borramos los ficheros de publicidad con dos años de antigüedad
func borrarPublicidad() {
	for {
		//año actual
		anio_actual := time.Now().Year()
		//PUBLICIDAD
		publi, errP := db.Query("SELECT id, fichero, fecha_fin FROM publi")
		if errP != nil {
			Error.Println(errP)
		}
		for publi.Next() {
			var id int
			var fichero, fecha_fin, anio, mes, dia string
			//Tomamos el id, nombre y playtime de la base de datos mensaje
			err := publi.Scan(&id, &fichero, &fecha_fin)
			if err != nil {
				Error.Println(err)
			}
			fmt.Sscanf(fecha_fin, "%4s%2s%2s", &anio, &mes, &dia)
			if libs.ToInt(anio) <= anio_actual-2 {
				//Borramos el fichero desde la ruta interna
				err = os.Remove(publi_files_location + fichero)
				if err != nil {
					Error.Println(err)
				}
				//Borramos de la base de datos los ficheros de publicidad
				db_mu.Lock()
				db.Exec("DELETE FROM publi WHERE fichero = ?", fichero)
				db_mu.Unlock()
			}
		}
		time.Sleep(20 * time.Hour) //Cada 20h revisa los ficheros para borrar
	}
}

//Inserta la publicidad en la base de datos de la tienda
func insert_publi(f_pub, existe, fecha_ini, fecha_fin, gap string) {
	stm, err := db.Prepare("INSERT INTO publi (`fichero`, `existe`, `fecha_ini`, `fecha_fin`, `gap`) VALUES (?,?,?,?,?)")
	if err != nil {
		Error.Println(err)
	}
	db_mu.Lock()
	_, err1 := stm.Exec(f_pub, existe, fecha_ini, fecha_fin, gap)
	db_mu.Unlock()
	if err1 != nil {
		Error.Println(err1)
	}
}

//Inserta los mensajes en la base de datos de la tienda
func insert_msg(msgname, existe, fecha_ini, fecha_fin, playtime string) {
	stm, err := db.Prepare("INSERT INTO mensaje (`fichero`, `existe`, `fecha_ini`, `fecha_fin`, `playtime`) VALUES (?,?,?,?,?)")
	if err != nil {
		Error.Println(err)
	}
	db_mu.Lock()
	_, err1 := stm.Exec(msgname, existe, fecha_ini, fecha_fin, playtime)
	db_mu.Unlock()
	if err1 != nil {
		Error.Println(err1)
	}
}

//Modifica los mensajes en la base de datos de la tienda
func update_msg(msgname, fecha_ini, fecha_fin, playtime string) {
	//Cambiamos el estado del fichero de publicidad en BD, a existe.
	ok, err := db.Prepare("UPDATE mensaje SET fecha_ini=?, fecha_fin=?, playtime=? WHERE fichero=?")
	if err != nil {
		Error.Println(err)
	}
	db_mu.Lock()
	_, err1 := ok.Exec(fecha_ini, fecha_fin, playtime, msgname)
	db_mu.Unlock()
	if err1 != nil {
		Error.Println(err1)
	}
}

//Modifica la publicidad en la base de datos de la tienda
func update_publi(f_pub, fecha_ini, fecha_fin, gap string) {
	//Cambiamos el estado del fichero de publicidad en BD, a existe.
	ok, err := db.Prepare("UPDATE publi SET fecha_ini=?, fecha_fin=?, gap=? WHERE fichero=?")
	if err != nil {
		Error.Println(err)
	}
	db_mu.Lock()
	_, err1 := ok.Exec(fecha_ini, fecha_fin, gap, f_pub)
	db_mu.Unlock()
	if err1 != nil {
		Error.Println(err1)
	}
}

func borrar_publi_int(listado string){
	q1, err := db.Query("SELECT fichero FROM publi")
	if err != nil {
		Error.Println(err)
	}
	for q1.Next() {
		var fichero string
		//Miramos y guardamos en un array los ficheros que tenemos en BD
		err = q1.Scan(&fichero)
		if err != nil {
			Error.Println(err)
		}
		//Borramos los mensajes que no coincidan con el listado
		if !strings.Contains(listado, fichero){
			//Borramos el fichero desde la ruta interna
			err = os.Remove(publi_files_location + fichero)
			if err != nil {
				Error.Println(err)
			}
			//Borramos de la base de datos los ficheros de publicidad
			db_mu.Lock()
			db.Exec("DELETE FROM publi WHERE fichero = ?", fichero)
			db_mu.Unlock()
		}
	}
}

func borrar_mensajes_int(listado string){
	q1, err := db.Query("SELECT fichero FROM mensaje")
	if err != nil {
		Error.Println(err)
	}
	for q1.Next() {
		var fichero string
		//Miramos y guardamos en un array los ficheros que tenemos en BD
		err = q1.Scan(&fichero)
		if err != nil {
			Error.Println(err)
		}
		//Borramos los mensajes que no coincidan con el listado
		if !strings.Contains(listado, fichero){
			//Borramos el fichero desde la ruta interna
			err = os.Remove(msg_files_location + fichero)
			if err != nil {
				Error.Println(err)
			}
			//Borramos de la base de datos los ficheros de mensaje
			db_mu.Lock()
			db.Exec("DELETE FROM mensaje WHERE fichero = ?", fichero)
			db_mu.Unlock()
		}
	}
}

func borrar_toda_publi(){
	q1, err := db.Query("SELECT fichero FROM publi")
	if err != nil {
		Error.Println(err)
	}
	for q1.Next() {
		var fichero string
		//Miramos y guardamos en un array los ficheros que tenemos en BD
		err = q1.Scan(&fichero)
		if err != nil {
			Error.Println(err)
		}
		//Borramos el fichero desde la ruta interna
		err = os.Remove(publi_files_location + fichero)
		if err != nil {
			Error.Println(err)
		}
		//Borramos de la base de datos los ficheros de mensaje
		db_mu.Lock()
		db.Exec("DELETE FROM publi WHERE fichero = ?", fichero)
		db_mu.Unlock()
	}
}
func borrar_todos_mensajes(){
	q1, err := db.Query("SELECT fichero FROM mensaje")
	if err != nil {
		Error.Println(err)
	}
	for q1.Next() {
		var fichero string
		//Miramos y guardamos en un array los ficheros que tenemos en BD
		err = q1.Scan(&fichero)
		if err != nil {
			Error.Println(err)
		}
		//Borramos el fichero desde la ruta interna
		err = os.Remove(msg_files_location + fichero)
		if err != nil {
			Error.Println(err)
		}
		//Borramos de la base de datos los ficheros de mensaje
		db_mu.Lock()
		db.Exec("DELETE FROM mensaje WHERE fichero = ?", fichero)
		db_mu.Unlock()
	}
}