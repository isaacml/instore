package main

import (
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/isaacml/instore/libs"
)

//Comparamos la hora guardada con la hora del sistema
func horario_reproduccion() {
	for {
		var hora_inicial, hora_final int
		var sol bool
		//Obtenemos la hora local
		clock := libs.MyCurrentClock()
		//Segmentamos para obtener horas y mins actuales
		arr_clock := strings.Split(clock, ":")
		//Pasamos las horas y minutos --> Minutos actuales totales
		actual := libs.Hour2min(libs.ToInt(arr_clock[0]), libs.ToInt(arr_clock[1]))
		//Obtenemos hora inicial y final de la SQL
		query, err := db.Query("SELECT hora_inicial, hora_final FROM aux")
		if err != nil {
			Warning.Println(err)
		}
		for query.Next() {
			err = query.Scan(&hora_inicial, &hora_final)
			if err != nil {
				Error.Println(err)
			}
			//Miramos que la hora actual de reproduccion esté dentro del rango
			if actual > hora_inicial && actual < hora_final {
				sol = sol || true
			}
		}
		db_mu.Lock()
		schedule = sol
		db_mu.Unlock()
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
			musica := make(map[int]string)
			pl := 1
			st_prog, err := libs.St_Prog_Music(db, db_mu)
			if err != nil {
				Error.Println(err)
			}
			if st_prog == "" {
				arr_music := libs.MusicToPlay(music_files, st_music)
				for k, v := range arr_music {
					musica[k] = v
				}
				rand.Seed(time.Now().UnixNano())
				shuffle := rand.Perm(len(musica))
				for _, v := range shuffle {
					st_prog_int, err := libs.St_Prog_Music(db, db_mu)
					if err != nil {
						Error.Println(err)
					}
					if st_prog_int == "PrimerCambio" || block == true || schedule == false {
						break
					}
					//Comprobamos si winamp está abierto
					isOpen := winplayer.WinampIsOpen()
					if isOpen == false {
						//Rulamos el Winamp
						winplayer.RunWinamp()
						time.Sleep(1 * time.Second)
						winplayer.Volume()
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
							winplayer.Load("\"" + song_to_play + "\"")
							winplayer.Play()
							//Esperamos el tiempo de duracion de la canción
							time.Sleep(time.Duration(winplayer.SongLenght(song_to_play)) * time.Second)
							//Una vez finalizada la reproduccion del fichero encriptado: Limpiamos la playlist
							winplayer.Clear()
							//Borramos el descifrado(.mp3)
							os.Remove(song_to_play)
						}
					} else {
						//Carga y reproduccion de cancion
						winplayer.Load("\"" + musica[v] + "\"")
						winplayer.Play()
						time.Sleep(time.Duration(winplayer.SongLenght(musica[v])) * time.Second)
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
							winplayer.Load("\"" + all_publi_file + "\"")
							winplayer.Play()
							//Esperamos lo que dure el archivo de publicidad
							time.Sleep(time.Duration(winplayer.SongLenght(all_publi_file)) * time.Second)
							break
						}
						//Volvemos a poner el contador de playlist 0
						pl = 0
					}
					pl++

				}
			}
			if st_prog == "PrimerCambio" {
				var carpeta string
				query, err := db.Query("SELECT carpeta FROM musica")
				if err != nil {
					Error.Println(err)
				}
				cont := 0
				for query.Next() {
					err = query.Scan(&carpeta)
					if err != nil {
						Error.Println(err)
					}
					//generamos la ruta completa a esas carpetas
					full_route := music_files + carpeta + "\\"
					arr_music := libs.MusicToPlay(full_route, st_music)
					for _, v := range arr_music {
						musica[cont] = v
						cont++
					}
				}
				rand.Seed(time.Now().UnixNano())
				shuffle := rand.Perm(len(musica))
				for _, v := range shuffle {
					st_prog_int, err := libs.St_Prog_Music(db, db_mu)
					if err != nil {
						Error.Println(err)
					}
					if st_prog_int == "SegundoCambio" || block == true || schedule == false {
						break
					}
					//Comprobamos si winamp está abierto
					isOpen := winplayer.WinampIsOpen()
					if isOpen == false {
						//Rulamos el Winamp
						winplayer.RunWinamp()
						time.Sleep(1 * time.Second)
						winplayer.Volume()
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
							winplayer.Load("\"" + song_to_play + "\"")
							winplayer.Play()
							//Esperamos el tiempo de duracion de la canción
							time.Sleep(time.Duration(winplayer.SongLenght(song_to_play)) * time.Second)
							//Una vez finalizada la reproduccion del fichero encriptado: Limpiamos la playlist
							winplayer.Clear()
							//Borramos el descifrado(.mp3)
							os.Remove(song_to_play)
						}
					} else {
						//Carga y reproduccion de cancion
						winplayer.Load("\"" + musica[v] + "\"")
						winplayer.Play()
						time.Sleep(time.Duration(winplayer.SongLenght(musica[v])) * time.Second)
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
							winplayer.Load("\"" + all_publi_file + "\"")
							winplayer.Play()
							time.Sleep(time.Duration(winplayer.SongLenght(all_publi_file)) * time.Second)
							break
						}
						//Volvemos a poner el contador de playlist 0
						pl = 0
					}
					pl++
				}
			}
			if st_prog == "SegundoCambio" {
				var carpeta string
				query, err := db.Query("SELECT carpeta FROM musica")
				if err != nil {
					Error.Println(err)
				}
				cont := 0
				for query.Next() {
					err = query.Scan(&carpeta)
					if err != nil {
						Error.Println(err)
					}
					//generamos la ruta completa a esas carpetas
					full_route := music_files + carpeta + "\\"
					arr_music := libs.MusicToPlay(full_route, st_music)
					for _, v := range arr_music {
						musica[cont] = v
						cont++
					}
				}
				rand.Seed(time.Now().UnixNano())
				shuffle := rand.Perm(len(musica))
				for _, v := range shuffle {
					st_prog_int, err := libs.St_Prog_Music(db, db_mu)
					if err != nil {
						Error.Println(err)
					}
					if st_prog_int == "PrimerCambio" || block == true || schedule == false {
						break
					}
					//Comprobamos si winamp está abierto
					isOpen := winplayer.WinampIsOpen()
					if isOpen == false {
						//Rulamos el Winamp
						winplayer.RunWinamp()
						time.Sleep(1 * time.Second)
						winplayer.Volume()
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
							winplayer.Load("\"" + song_to_play + "\"")
							winplayer.Play()
							//Esperamos el tiempo de duracion de la canción
							time.Sleep(time.Duration(winplayer.SongLenght(song_to_play)) * time.Second)
							//Una vez finalizada la reproduccion del fichero encriptado: Limpiamos la playlist
							winplayer.Clear()
							//Borramos el descifrado(.mp3)
							os.Remove(song_to_play)
						}
					} else {
						//Carga y reproduccion de cancion
						winplayer.Load("\"" + musica[v] + "\"")
						winplayer.Play()
						time.Sleep(time.Duration(winplayer.SongLenght(musica[v])) * time.Second)
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
							winplayer.Load("\"" + all_publi_file + "\"")
							winplayer.Play()
							time.Sleep(time.Duration(winplayer.SongLenght(all_publi_file)) * time.Second)
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
		if block != true {
			if schedule != false {
				//Obtenemos la fecha actual
				fecha := libs.MyCurrentDate()
				//Obtenemos la hora local
				clock := libs.MyCurrentClock()
				//Obtenemos todos los mensajes
				mensajes, errM := db.Query("SELECT id, fichero, fecha_ini, fecha_fin, playtime FROM mensaje")
				if errM != nil {
					Error.Println(errM)
				}
				for mensajes.Next() {
					var id int
					var fichero, fecha_ini, fecha_fin, playtime string
					//Tomamos el id, nombre y playtime de la base de datos mensaje
					err := mensajes.Scan(&id, &fichero, &fecha_ini, &fecha_fin, &playtime)
					if err != nil {
						Error.Println(err)
					}
					//BETWEEN
					if fecha_ini <= fecha && fecha_fin >= fecha {
						if playtime == clock {
							go winplayer.PlayFFplay(msg_files_location + fichero)
						}
					}
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
	var fichero, fecha_ini, fecha_fin string
	publicidad, errP := db.Query("SELECT fichero, fecha_ini, fecha_fin, gap FROM publi")
	if errP != nil {
		Error.Println(errP)
		gap = 0
	}
	for publicidad.Next() {
		//Tomamos el nombre del fichero publi
		err := publicidad.Scan(&fichero, &fecha_ini, &fecha_fin, &gap)
		if err != nil {
			Error.Println(err)
		}
		//BETWEEN
		if fecha_ini <= fecha && fecha_fin >= fecha {
			publi[p] = fichero
			p++
		}
	}
	return publi, gap
}
