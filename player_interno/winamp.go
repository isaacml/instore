package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"github.com/isaacml/instore/winamp"
	"math/rand"
	"os"
	"time"
)

//Zona de reproduccion del player de la tienda
func reproduccion() {
	for {
		if estado_entidad != 0 {
			if estado_entidad == 0 {
				continue
			}
			var win winamp.Winamp
			publi := make(map[int]string)
			musica := make(map[int]string)
			p, pl := 0, 1
			var gap int
			var song string
			//Sacamos la fecha actual
			fecha := libs.MyCurrentDate()
			//INICIAL
			fmt.Println("-" + statusProgammedMusic + "-")
			//Obtenemos el GAP
			publicidad, errP := db.Query("SELECT fichero, gap FROM publi  WHERE fecha_ini = ?", fecha)
			if errP != nil {
				Error.Println(errP)
				gap = 0
			}
			for publicidad.Next() {
				var fichero string
				//Tomamos el nombre del fichero mensaje
				err := publicidad.Scan(&fichero, &gap)
				if err != nil {
					Error.Println(err)
				}
				//fmt.Printf("%s", fichero)
				publi[p] = fichero
				p++
			}
			//Comprobamos si winamp está abierto
			isOpen := win.WinampIsOpen()
			if isOpen == false {
				//Rulamos el Winamp
				win.RunWinamp()
				time.Sleep(1 * time.Second)
				win.Volume()
			}
			if statusProgammedMusic == "Inicial" {
				for _, val := range programmedMusic {
					//generamos la ruta completa a esas carpetas
					full_route := music_files + val + "\\"
					musica = libs.MusicToPlay(full_route)
				}
				rand.Seed(time.Now().UnixNano())
				shuffle := rand.Perm(len(musica))
				for _, v := range shuffle {
					if statusProgammedMusic == "Actualizada" || estado_entidad == 0 {
						break
					}
					song = musica[v]
					libs.PlaySong(song, win)
					//Controlamos el GAP: Cuando el contador de canciones es igual al número de gap, metemos publicidad.
					//Un gap = 0 --> No hay publicidad, las canciones corren una detrás de otra.
					if pl == gap {
						//Movemos aleatoriamente todos los ficheros publi guardados en nuestro arr.
						rand.Seed(time.Now().UnixNano())
						shuffle2 := rand.Perm(len(publi))
						//Una vez mezclado, cogemos el primer fichero de publicidad y lo reproducimos.
						for _, val := range shuffle2 {
							publi_file := publi[val]
							all_publi_file := publi_files_location + publi_file
							libs.PlayPubli(all_publi_file, win)
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
					musica = libs.MusicToPlay(full_route)
				}
				rand.Seed(time.Now().UnixNano())
				shuffle := rand.Perm(len(musica))
				for _, v := range shuffle {
					if statusProgammedMusic == "Modificar" || estado_entidad == 0 {
						break
					}
					song = musica[v]
					libs.PlaySong(song, win)
					//Controlamos el GAP: Cuando el contador de canciones es igual al número de gap, metemos publicidad.
					//Un gap = 0 --> No hay publicidad, las canciones corren una detrás de otra.
					if pl == gap {
						//Movemos aleatoriamente todos los ficheros publi guardados en nuestro arr.
						rand.Seed(time.Now().UnixNano())
						shuffle2 := rand.Perm(len(publi))
						//Una vez mezclado, cogemos el primer fichero de publicidad y lo reproducimos.
						for _, val := range shuffle2 {
							publi_file := publi[val]
							all_publi_file := publi_files_location + publi_file
							libs.PlayPubli(all_publi_file, win)
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
					musica = libs.MusicToPlay(full_route)
				}
				rand.Seed(time.Now().UnixNano())
				shuffle := rand.Perm(len(musica))
				for _, v := range shuffle {
					if statusProgammedMusic == "Actualizada" || estado_entidad == 0 {
						break
					}
					song = musica[v]
					libs.PlaySong(song, win)
					//Controlamos el GAP: Cuando el contador de canciones es igual al número de gap, metemos publicidad.
					//Un gap = 0 --> No hay publicidad, las canciones corren una detrás de otra.
					if pl == gap {
						//Movemos aleatoriamente todos los ficheros publi guardados en nuestro arr.
						rand.Seed(time.Now().UnixNano())
						shuffle2 := rand.Perm(len(publi))
						//Una vez mezclado, cogemos el primer fichero de publicidad y lo reproducimos.
						for _, val := range shuffle2 {
							publi_file := publi[val]
							all_publi_file := publi_files_location + publi_file
							libs.PlayPubli(all_publi_file, win)
							break
						}
						//Volvemos a poner el contador de playlist 0
						pl = 0
					}
					pl++
				}
			} else {
				var song string
				musica = libs.MusicToPlay(music_files)
				rand.Seed(time.Now().UnixNano())
				shuffle := rand.Perm(len(musica))
				for _, v := range shuffle {
					if statusProgammedMusic == "Inicial" || estado_entidad == 0 {
						break
					}
					song = musica[v] //Tomamos las canciones, teniendo en cuenta que hay musica cif/NO cif
					libs.PlaySong(song, win)
					//Controlamos el GAP: Cuando el contador de canciones es igual al número de gap, metemos publicidad.
					//Un gap = 0 --> No hay publicidad, las canciones corren una detrás de otra.
					if pl == gap {
						//Movemos aleatoriamente todos los ficheros publi guardados en nuestro arr.
						rand.Seed(time.Now().UnixNano())
						shuffle2 := rand.Perm(len(publi))
						//Una vez mezclado, cogemos el primer fichero de publicidad y lo reproducimos.
						for _, val := range shuffle2 {
							publi_file := publi[val]
							all_publi_file := publi_files_location + publi_file
							libs.PlayPubli(all_publi_file, win)
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
		y, m, d := time.Now().Date()
		fecha := fmt.Sprintf("%4d%02d%02d", y, int(m), d)
		fecha_sql := fecha + "%"
		//Obtenemos la hora local
		hh, mm, _ := time.Now().Clock()
		clock := fmt.Sprintf("%02d:%02d", hh, mm)
		//Obtenemos todos los mensajes
		mensajes, errM := db.Query("SELECT id, fichero, playtime FROM mensaje WHERE fichero LIKE ?", fecha_sql)
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
