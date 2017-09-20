package main

import (
	"fmt"
	"math/rand"
	"strings"
	//"sync"
	"github.com/isaacml/instore/libs"
	"github.com/isaacml/instore/winamp"
	"os"
	"time"
)
					
func reproduccion() {
	for {
		//a, p, pl := 0, 0, 1
		p := 0
		var gap int
		var song string
		var win winamp.Winamp
		publi := make(map[int]string)
		musica := make(map[int]string)
		//Sacamos la fecha actual
		y, m, d := time.Now().Date()
		fecha := fmt.Sprintf("%4d%02d%02d", y, int(m), d)
		//INICIAL
		if statusProgammedMusic == "Inicial" {
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
			for _, val := range programmedMusic {
				//generamos la ruta completa a esas carpetas
				full_route := music_files + val + "\\"
				musica = libs.MusicToPlay(full_route)
			}
			rand.Seed(time.Now().UnixNano())
			shuffle := rand.Perm(len(musica))
			//Rulamos el Winamp
			win.RunWinamp()
			for _, v := range shuffle {
				var song_duration, song_end int
				song = musica[v]
				if strings.Contains(song, ".xxx") {
					del_ext := strings.Split(song, ".xxx")
					descifrada := del_ext[0] + ".mp3"
					//Proceso de descifrado de la cancion: ver en libreria de funciones.
					libs.Cifrado(song, descifrada, []byte{11, 22, 33, 44, 55, 66, 77, 88})
					//Carga y reproduccion de cancion
					win.Load("\"" + descifrada + "\"")
					win.Play()
					//Guardamos la duracion total de la cancion
					song_duration = win.SongLenght(descifrada)
					for {
						song_end = win.SongEnd()
						song_play := win.SongPlay()
						if song_play == song_duration {
							err := os.Remove(descifrada)
							fmt.Println(err)
							continue
						}
						fmt.Println(song_duration, song_end, song_play)
						time.Sleep(1 * time.Second)
					}
				}
			}
		}
		time.Sleep(30 * time.Second)
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
