package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	//"sync"
	"github.com/isaacml/instore/winamp"
	"io/ioutil"
	"os"
	"time"
)

func reproduccion() {
	for {
		a, p, pl := 0, 0, 1
		var gap int
		var song string
		var win winamp.Winamp
		music := make(map[int]string)
		publi := make(map[int]string)
		//Sacamos la fecha actual
		y, m, d := time.Now().Date()
		fecha := fmt.Sprintf("%4d%02d%02d", y, int(m), d)
		//Obtenemos el GAP
		publicidad, errP := db.Query("SELECT fichero, gap FROM publi  WHERE fecha_ini = ? OR fecha_fin >= ?", fecha, fecha)
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
		fmt.Println(publi)
		//Se crean lo comandos
		cmd := exec.Command("cmd", "/c", "dir /s /b "+music_files+"*.mp3 & dir /s /b "+music_files+"*.xxx")
		// comienza la ejecucion del pipe
		stdoutRead, _ := cmd.StdoutPipe()
		reader := bufio.NewReader(stdoutRead)
		cmd.Start()
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			//fmt.Printf("%s", line)
			music[a] = strings.TrimRight(line, "\r\n")
			a++
		}
		cmd.Wait()
		rand.Seed(time.Now().UnixNano())
		shuffle := rand.Perm(len(music))
		//Rulamos el Winamp
		win.RunWinamp()
		//ZONA DE CREACION DE PLAYLIST
		//Este bucle va a mezclar la musica con la publicidad segun el GAP
		for _, v := range shuffle {
			var song_duration int
			song = music[v] //Tomamos las canciones, teniendo en cuenta que hay musica cif/NO cif
			fmt.Println(song)
			// .xxx = musica cif; Hay que descifrarla
			if strings.Contains(song, ".xxx") {
				del_ext := strings.Split(song, ".xxx")
				descifrada := del_ext[0] + ".mp3"
				//Proceso de descifrado de la cancion: ver función mas abajo.
				cifrado(song, descifrada, []byte{11, 22, 33, 44, 55, 66, 77, 88})
				//Carga y reproduccion de cancion
				win.Load("\"" + descifrada + "\"")
				win.Play()
				//Guardamos la duracion total de la cancion
				song_duration = win.SongLenght(descifrada)
			} else {
				//Carga y reproduccion de cancion
				win.Load("\"" + song + "\"")
				win.Play()
				//Guardamos la duracion total de la cancion
				song_duration = win.SongLenght(song)
			}
			//Controlamos el GAP: Cuando el contador de canciones es igual al número de gap, metemos publicidad.
			//Un gap = 0 --> No hay publicidad, las canciones corren una detrás de otra.
			if pl == gap {
				//Movemos aleatoriamente todos los ficheros publi guardados en nuestro arr.
				rand.Seed(time.Now().UnixNano())
				shuffle2 := rand.Perm(len(publi))
				//Una vez mezclado, cogemos el primer fichero de publicidad y lo reproducimos.
				for _, val := range shuffle2 {
					publi_file := publi[val]
					fmt.Println(publi_file)
					win.Load("\"" + publi_files_location + publi_file + "\"")
					win.Play()
					song_duration = win.SongLenght(publi_files_location + publi_file)
					break
				}
				//Volvemos a poner el contador de playlist 0
				pl = 0
			}
			//time.Duration(song_duration)
			time.Sleep(time.Duration(song_duration) * time.Second)
			pl++
		}
	}
}

func reproduccion_msgs() {
	for {
		hh, mm, _ := time.Now().Clock()
		clock := fmt.Sprintf("%2d:%2d", hh, mm)
		fmt.Println(clock)
		mensajes, errM := db.Query("SELECT fichero FROM mensaje WHERE playtime=?", clock)
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
			var win winamp.Winamp
			fmt.Printf("%s\n", fichero)
			win.PlayFFplay(fichero)
		}
		time.Sleep(1 * time.Minute)
	}
}

//Este proceso lo que hace es cifrar o descifrar un fichero existente
func cifrado(origen, destino string, key []byte) error {
	var fail error
	p := make([]byte, 8) //Va a contener el archivo origen en bloques de 8 bytes
	var container []byte //Va almacenar los datos del fichero de destino
	file, err := os.OpenFile(origen, os.O_RDONLY, 0666)
	if err != nil {
		fail = fmt.Errorf("Error en la apertura")
	}
	lector := bufio.NewReader(file)
	for {
		num, err := lector.Read(p)
		if err != nil {
			fail = fmt.Errorf("Fin de lectura")
			break
		}
		if num <= 0 {
			break
		} else {
			for i := 0; i < num; i++ {
				container = append(container, p[i]^key[i])
			}
		}
	}
	//Escribimos el fichero
	ioutil.WriteFile(destino, container, 0666)
	return fail
}
