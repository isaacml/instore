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
		i, a, cont := 0, 1, 0
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
			publi[cont] = fichero
			cont++
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
			fmt.Println(i, gap)

			//fmt.Printf("%s", line)
			music[i] = strings.TrimRight(line, "\r\n")
			i++
			a++
		}
		cmd.Wait()
		rand.Seed(time.Now().UnixNano())
		shuffle := rand.Perm(len(music))
		rand.Seed(time.Now().UnixNano())
		shuffle2 := rand.Perm(len(publi))
		fmt.Println(shuffle2)

		//Rulamos el Winamp
		win.RunWinamp()
		//Este bucle va a mezclar la musica con la publicidad segun el GAP
		for _, v := range shuffle {
			//var song_duration int
			song = music[v]
			fmt.Println(song)
			// .xxx = musica cifrada; Hay que descifrarla
			if strings.Contains(song, ".xxx") {
				del_ext := strings.Split(song, ".xxx")
				descifrada := del_ext[0] + ".mp3"
				//Proceso de descifrado de la cancion
				cifrado(song, descifrada, []byte{11, 22, 33, 44, 55, 66, 77, 88})
				//Carga de cancion y reproduccion de la cancion
				win.Load("\"" + descifrada + "\"")
				win.Play()
				//Guardamos la duracion total de la cancion
				//song_duration = win.SongLenght(descifrada)
			} else {
				//Carga de cancion y reproduccion de la cancion
				win.Load("\"" + song + "\"")
				win.Play()
				//Guardamos la duracion total de la cancion
				//song_duration = win.SongLenght(song)
			}
			//Cuando el contador de canciones es igual al número de gap, metemos publicidad
			//Un gap = 0, significa que no hay publicidad
			fmt.Println(a, gap)
			/*
				if a == gap {
					for _, val := range shuffle2 {
						arch_publicidad := publi[val]
						total := win.SongEnd()
						inicial := win.SongPlay()
						fmt.Println(arch_publicidad, total, inicial)

							for {
								fmt.Println(total - win.SongEnd())
								if total-win.SongEnd() == 0 {
									fmt.Println("fichero_publi-> " + arch_publicidad)
									win.Load("\"" + arch_publicidad + "\"")
									win.Play()
									break
								}
								time.Sleep(1 * time.Second)
							}

					}

					a = 0
				}
			*/
			//time.Duration(song_duration)
			time.Sleep(10 * time.Second)
			a++
		}
	}
}

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
	ioutil.WriteFile(destino, container, 0666)
	return fail
}
