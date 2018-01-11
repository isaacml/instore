package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"github.com/isaacml/instore/winamp"
	"math/rand"
	"os"
	"strings"
	"time"
)

//Comparamos la hora guardada con la hora del sistema
func horario_reproduccion() {
	for {
		var actual int                      //Minutos actuales totales
		var hora_inicial, hora_final string //Variables de base de datos
		mins_of_day := 1440                 //Minutos que tiene un dia
		//Obtenemos la hora local
		clock := libs.MyCurrentClock()
		//Obtenemos hora inicial y final de la SQL
		db.QueryRow("SELECT hora_inicial, hora_final FROM horario").Scan(&hora_inicial, &hora_final)
		//Comprobamos que los datos obtenidos en base de datos no son vacíos
		if hora_inicial != "" && hora_final != "" {
			//Segmentamos para obtener horas y mins actuales
			arr_clock := strings.Split(clock, ":")
			//Segmentamos para obtener la hora inicial y final
			arr_hinicial := strings.Split(hora_inicial, ":")
			arr_hfinal := strings.Split(hora_final, ":")
			//Se comprueba que el array no está vacío
			if len(arr_clock) > 1 || len(arr_hinicial) > 1 || len(arr_hfinal) > 1 {
				//Pasamos las horas y minutos a minutos totales
				actual = libs.Hour2min(libs.ToInt(arr_clock[0]), libs.ToInt(arr_clock[1]))
				inicial := libs.Hour2min(libs.ToInt(arr_hinicial[0]), libs.ToInt(arr_hinicial[1]))
				final := libs.Hour2min(libs.ToInt(arr_hfinal[0]), libs.ToInt(arr_hfinal[1]))
				//Si la hora inicial es mayor que la final, pasa un dia completo
				if inicial > final {
					final = mins_of_day + final
					actual = mins_of_day + actual
				}
				fmt.Println(inicial, actual, final)
				//Miramos que la hora actual de reproduccion esté dentro del rango
				if actual >= inicial && final >= actual {
					db_mu.Lock()
					schedule = true
					db_mu.Unlock()
				}
				//Fuera de horario
				if actual > final {
					db_mu.Lock()
					schedule = false
					db_mu.Unlock()
				}
			}
		}
		time.Sleep(1 * time.Minute)
	}
}

//Zona de reproduccion del player de la tienda
func reproduccion() {
	for {
		if block == false && schedule == true {
			if block == true && schedule == false {
				continue
			}
			var win winamp.Winamp
			musica := make(map[int]string)
			pl := 1
			if statusProgammedMusic == "Inicial" {
				for _, val := range programmedMusic {
					//generamos la ruta completa a esas carpetas
					full_route := music_files + val + "\\"
					libs.MusicToPlay(full_route, st_music, musica)
				}
				rand.Seed(time.Now().UnixNano())
				shuffle := rand.Perm(len(musica))
				for _, v := range shuffle {
					if statusProgammedMusic == "Actualizada" || block == true || schedule == false {
						break
					}
					//Comprobamos si winamp está abierto
					isOpen := win.WinampIsOpen()
					if isOpen == false {
						//Rulamos el Winamp
						win.RunWinamp()
						time.Sleep(1 * time.Second)
						win.Volume()
					}
					//Obtenemos la publicidad
					publi, gap := publi_q_toca()
					//Evaluamos cada una de las canciones: cif o nocif
					if strings.Contains(musica[v], ".xxx") {
						segment := strings.Split(musica[v], ".xxx")
						song_to_play := segment[0] + ".mp3"
						//Proceso de descifrado de la cancion: ver en libreria de funciones.
						_, st_cif := libs.Cifrado(musica[v], song_to_play, []byte{11, 22, 33, 44, 55, 66, 77, 88})
						if st_cif == "GOOD" {
							//Carga y reproduccion de cancion
							win.Load("\"" + song_to_play + "\"")
							win.Play()
							//Esperamos el tiempo de duracion de la canción
							time.Sleep(time.Duration(win.SongLenght(song_to_play)) * time.Second)
							//Una vez finalizada la reproduccion del fichero encriptado: Limpiamos la playlist
							win.Clear()
							//Borramos el descifrado(.mp3)
							os.Remove(song_to_play)
						}
					} else {
						//Carga y reproduccion de cancion
						win.Load("\"" + musica[v] + "\"")
						win.Play()
						time.Sleep(time.Duration(win.SongLenght(musica[v])) * time.Second)
					}
					//Controlamos el GAP: Cuando el contador de canciones es igual al número de gap, metemos publicidad.
					//Un gap = 0 --> No hay publicidad, las canciones corren una detrás de otra.
					if pl == gap {
						//Movemos aleatoriamente todos los ficheros publi guardados en nuestro arr.
						rand.Seed(time.Now().UnixNano())
						shuffle2 := rand.Perm(len(publi))
						//Una vez mezclado, cogemos el primer fichero de publicidad y lo reproducimos.
						for _, val := range shuffle2 {
							//Directorio publi + Fichero publi
							all_publi_file := publi_files_location + publi[val]
							win.Load("\"" + all_publi_file + "\"")
							win.Play()
							time.Sleep(time.Duration(win.SongLenght(all_publi_file)) * time.Second)
							break
						}
						//Volvemos a poner el contador de playlist 0
						pl = 0
					}
					pl++
				}
			} else if statusProgammedMusic == "Actualizada" {
				for _, val := range programmedMusic {
					//generamos la ruta completa a esas carpetas
					full_route := music_files + val + "\\"
					libs.MusicToPlay(full_route, st_music, musica)
				}
				rand.Seed(time.Now().UnixNano())
				shuffle := rand.Perm(len(musica))
				for _, v := range shuffle {
					if statusProgammedMusic == "Modificar" || block == true || schedule == false {
						break
					}
					//Comprobamos si winamp está abierto
					isOpen := win.WinampIsOpen()
					if isOpen == false {
						//Rulamos el Winamp
						win.RunWinamp()
						time.Sleep(1 * time.Second)
						win.Volume()
					}
					//Obtenemos la publicidad
					publi, gap := publi_q_toca()
					//Evaluamos cada una de las canciones: cif o nocif
					if strings.Contains(musica[v], ".xxx") {
						segment := strings.Split(musica[v], ".xxx")
						song_to_play := segment[0] + ".mp3"
						//Proceso de descifrado de la cancion: ver en libreria de funciones.
						_, st_cif := libs.Cifrado(musica[v], song_to_play, []byte{11, 22, 33, 44, 55, 66, 77, 88})
						if st_cif == "GOOD" {
							//Carga y reproduccion de cancion
							win.Load("\"" + song_to_play + "\"")
							win.Play()
							//Esperamos el tiempo de duracion de la canción
							time.Sleep(time.Duration(win.SongLenght(song_to_play)) * time.Second)
							//Una vez finalizada la reproduccion del fichero encriptado: Limpiamos la playlist
							win.Clear()
							//Borramos el descifrado(.mp3)
							os.Remove(song_to_play)
						}
					} else {
						//Carga y reproduccion de cancion
						win.Load("\"" + musica[v] + "\"")
						win.Play()
						time.Sleep(time.Duration(win.SongLenght(musica[v])) * time.Second)
					}
					//Controlamos el GAP: Cuando el contador de canciones es igual al número de gap, metemos publicidad.
					//Un gap = 0 --> No hay publicidad, las canciones corren una detrás de otra.
					if pl == gap {
						//Movemos aleatoriamente todos los ficheros publi guardados en nuestro arr.
						rand.Seed(time.Now().UnixNano())
						shuffle2 := rand.Perm(len(publi))
						//Una vez mezclado, cogemos el primer fichero de publicidad y lo reproducimos.
						for _, val := range shuffle2 {
							//Directorio publi + Fichero publi
							all_publi_file := publi_files_location + publi[val]
							win.Load("\"" + all_publi_file + "\"")
							win.Play()
							time.Sleep(time.Duration(win.SongLenght(all_publi_file)) * time.Second)
							break
						}
						//Volvemos a poner el contador de playlist 0
						pl = 0
					}
					pl++
				}
			} else if statusProgammedMusic == "Modificar" {
				for _, val := range programmedMusic {
					//generamos la ruta completa a esas carpetas
					full_route := music_files + val + "\\"
					libs.MusicToPlay(full_route, st_music, musica)
				}
				rand.Seed(time.Now().UnixNano())
				shuffle := rand.Perm(len(musica))
				for _, v := range shuffle {
					if statusProgammedMusic == "Actualizada" || block == true || schedule == false {
						break
					}
					//Comprobamos si winamp está abierto
					isOpen := win.WinampIsOpen()
					if isOpen == false {
						//Rulamos el Winamp
						win.RunWinamp()
						time.Sleep(1 * time.Second)
						win.Volume()
					}
					//Obtenemos la publicidad
					publi, gap := publi_q_toca()
					//Evaluamos cada una de las canciones: cif o nocif
					if strings.Contains(musica[v], ".xxx") {
						segment := strings.Split(musica[v], ".xxx")
						song_to_play := segment[0] + ".mp3"
						//Proceso de descifrado de la cancion: ver en libreria de funciones.
						_, st_cif := libs.Cifrado(musica[v], song_to_play, []byte{11, 22, 33, 44, 55, 66, 77, 88})
						if st_cif == "GOOD" {
							//Carga y reproduccion de cancion
							win.Load("\"" + song_to_play + "\"")
							win.Play()
							//Esperamos el tiempo de duracion de la canción
							time.Sleep(time.Duration(win.SongLenght(song_to_play)) * time.Second)
							//Una vez finalizada la reproduccion del fichero encriptado: Limpiamos la playlist
							win.Clear()
							//Borramos el descifrado(.mp3)
							os.Remove(song_to_play)
						}
					} else {
						//Carga y reproduccion de cancion
						win.Load("\"" + musica[v] + "\"")
						win.Play()
						time.Sleep(time.Duration(win.SongLenght(musica[v])) * time.Second)
					}
					//Controlamos el GAP: Cuando el contador de canciones es igual al número de gap, metemos publicidad.
					//Un gap = 0 --> No hay publicidad, las canciones corren una detrás de otra.
					if pl == gap {
						//Movemos aleatoriamente todos los ficheros publi guardados en nuestro arr.
						rand.Seed(time.Now().UnixNano())
						shuffle2 := rand.Perm(len(publi))
						//Una vez mezclado, cogemos el primer fichero de publicidad y lo reproducimos.
						for _, val := range shuffle2 {
							//Directorio publi + Fichero publi
							all_publi_file := publi_files_location + publi[val]
							win.Load("\"" + all_publi_file + "\"")
							win.Play()
							time.Sleep(time.Duration(win.SongLenght(all_publi_file)) * time.Second)
							break
						}
						//Volvemos a poner el contador de playlist 0
						pl = 0
					}
					pl++
				}
			} else {
				libs.MusicToPlay(music_files, st_music, musica)
				rand.Seed(time.Now().UnixNano())
				shuffle := rand.Perm(len(musica))
				for _, v := range shuffle {
					if statusProgammedMusic == "Inicial" || block == true || schedule == false {
						break
					}
					//Comprobamos si winamp está abierto
					isOpen := win.WinampIsOpen()
					if isOpen == false {
						//Rulamos el Winamp
						win.RunWinamp()
						time.Sleep(1 * time.Second)
						win.Volume()
					}
					//Obtenemos la publicidad
					publi, gap := publi_q_toca()
					//Evaluamos cada una de las canciones: cif o nocif
					if strings.Contains(musica[v], ".xxx") {
						segment := strings.Split(musica[v], ".xxx")
						song_to_play := segment[0] + ".mp3"
						//Proceso de descifrado de la cancion: ver en libreria de funciones.
						_, st_cif := libs.Cifrado(musica[v], song_to_play, []byte{11, 22, 33, 44, 55, 66, 77, 88})
						if st_cif == "GOOD" {
							//Carga y reproduccion de cancion
							win.Load("\"" + song_to_play + "\"")
							win.Play()
							//Esperamos el tiempo de duracion de la canción
							time.Sleep(time.Duration(win.SongLenght(song_to_play)) * time.Second)
							//Una vez finalizada la reproduccion del fichero encriptado: Limpiamos la playlist
							win.Clear()
							//Borramos el descifrado(.mp3)
							os.Remove(song_to_play)
						}
					} else {
						//Carga y reproduccion de cancion
						win.Load("\"" + musica[v] + "\"")
						win.Play()
						time.Sleep(time.Duration(win.SongLenght(musica[v])) * time.Second)
					}
					//Controlamos el GAP: Cuando el contador de canciones es igual al número de gap, metemos publicidad.
					//Un gap = 0 --> No hay publicidad, las canciones corren una detrás de otra.
					if pl == gap {
						//Movemos aleatoriamente todos los ficheros publi guardados en nuestro arr.
						rand.Seed(time.Now().UnixNano())
						shuffle2 := rand.Perm(len(publi))
						//Una vez mezclado, cogemos el primer fichero de publicidad y lo reproducimos.
						for _, val := range shuffle2 {
							//Directorio publi + Fichero publi
							all_publi_file := publi_files_location + publi[val]
							win.Load("\"" + all_publi_file + "\"")
							win.Play()
							//Esperamos lo que dure el archivo de publicidad
							time.Sleep(time.Duration(win.SongLenght(all_publi_file)) * time.Second)
							break
						}
						//Volvemos a poner el contador de playlist 0
						pl = 0
					}
					pl++
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}

//Reproduce los mensajes automáticos de la tienda: bucle infinito que busca cada minuto un mensaje nuevo para reproducir.
func reproduccion_msgs() {
	for {
		//Obtenemos la fecha actual
		fecha := libs.MyCurrentDate()
		fecha_sql := fecha + "%"
		//Obtenemos la hora local
		clock := libs.MyCurrentClock()
		//Obtenemos todos los mensajes
		mensajes, errM := db.Query("SELECT id, fichero, playtime FROM mensaje WHERE fecha LIKE ?", fecha_sql)
		if errM != nil {
			Error.Println(errM)
		}
		for mensajes.Next() {
			var id int
			var fichero, playtime string
			//Tomamos el id, nombre y playtime de la base de datos mensaje
			err := mensajes.Scan(&id, &fichero, &playtime)
			if err != nil {
				Error.Println(err)
			}
			if playtime == clock {
				var win winamp.Winamp
				st := win.PlayFFplay(msg_files_location + fichero)
				//Si el estado de la reproduccion del mensaje = END (ha acabado), procedemos al borrado.
				if st == "END" {
					fmt.Println("Borro el fichero: "+fichero, id)
					//Borramos el fichero desde el directorio que contiene los mensajes en el player_int
					err = os.Remove(msg_files_location + fichero)
					if err != nil {
						Error.Println(err)
					}
					//Ponemos el estado de mensaje en N (ya que lo hemos borrado y no existe)
					ok, err := db.Prepare("UPDATE mensaje SET existe=? WHERE id = ?")
					if err != nil {
						Error.Println(err)
					}
					db_mu.Lock()
					_, err1 := ok.Exec("N", id)
					if err1 != nil {
						Error.Println(err1)
					}
					db_mu.Unlock()
					break
				}
			}
		}
		time.Sleep(1 * time.Minute)
	}
}

//Hace un select de la publicidad diaria y la guarda en un mapa junto con el GAP
func publi_q_toca() (map[int]string, int) {
	p := 0
	publi := make(map[int]string)
	//Sacamos la fecha actual
	fecha := libs.MyCurrentDate()
	//Obtenemos el GAP
	var gap int
	publicidad, errP := db.Query("SELECT fichero, gap FROM publi  WHERE fecha_ini = ?", fecha)
	if errP != nil {
		Error.Println(errP)
		gap = 0
	}
	for publicidad.Next() {
		var fichero string
		//Tomamos el nombre del fichero publi
		err := publicidad.Scan(&fichero, &gap)
		if err != nil {
			Error.Println(err)
		}
		publi[p] = fichero
		p++
	}
	return publi, gap
}
