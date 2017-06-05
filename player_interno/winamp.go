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
	i := 0
	music := make(map[int]string)
	//Sacamos la fecha actual
	y, m, d := time.Now().Date()
	fecha := fmt.Sprintf("%4d%02d%02d", y, int(m), d)
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
		music[i] = strings.TrimRight(line, "\r\n")
		i++
	}
	cmd.Wait()
	rand.Seed(time.Now().UnixNano())
	shuffle := rand.Perm(len(music))
	fmt.Println(shuffle)

	var gap int
	publicidad, errP := db.Query("SELECT fichero, gap FROM publi WHERE fichero LIKE ?", fecha+"%")
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
	}
	var win winamp.Winamp
	win.RunWinamp()
	a := 1
	for _, v := range shuffle {

		fmt.Println(music[v])
		// .xxx = musica cifrada; Hay que descifrarla
		if strings.Contains(music[v], ".xxx") {
			del_ext := strings.Split(music[v], ".xxx")
			descifrada := del_ext[0] + ".mp3"
				err := cifrado(music[v], descifrada, []byte{11, 22, 33, 44, 55, 66, 77, 88})
				if err != nil {
					Error.Println(err)
				} else {
					win.Load("\"" + descifrada + "\"")
					win.Play()
					fmt.Println(win.SongLenght())
					time.Sleep(50 * time.Second)
				}
		} else {
			win.Load("\"" + music[v] + "\"")
			win.Play()
			fmt.Println(win.SongLenght())
			time.Sleep(50 * time.Second)
			if a == gap {
				fmt.Println("Meto publicidad")
				a = 0
			}
		}
		a++
	}
	/*
			i = 0
			cmd = exec.Command("cmd", "/c", "dir /s /b "+publi_files_location+"*.mp3")
			// comienza la ejecucion del pipe
			stdoutRead, _ = cmd.StdoutPipe()
			reader = bufio.NewReader(stdoutRead)
			cmd.Start()
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					break
				}
				//fmt.Printf("%s", line)
				publi[i] = strings.TrimRight(line, "\n")
				i++
			}
			cmd.Wait()
			rand.Perm(len(publi))


		var gap int
		publicidad, errP := db.Query("SELECT fichero, gap FROM publi WHERE fichero LIKE ?", fecha+"%")
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
			publi[i] = fichero
			i++
		}
		shuffle2 := rand.Perm(len(publi))
		fmt.Println(shuffle2)
		for _, v2 := range shuffle2 {
			for k, v := range shuffle {

				fmt.Println(music[v])
				if gap-1 == k {
					fmt.Println(publi[v2])
				}
			}
		}
	*/
	//creamos la playlist mezclando audio + pub con el gap correspondiente
	//a, p, i := 0, 0, 0
	/*
		getGap, errG := db.Query("SELECT gap FROM publi WHERE fichero LIKE ?", fecha+"%")
		if errG != nil {
			Error.Println(errG)
		}
		if getGap.Next() {
			var gap string
			err = getGap.Scan(&gap)
			if err != nil {
				Error.Println(err)
			}
			gap_int, err = strconv.Atoi(gap)
			if err != nil {
				Error.Println(err)
			}
		}

	*/
	/*
		for _, v := range shuffle {

			fmt.Println(a, audio[v])
			if a == gap_int {
				fmt.Println("meto publicidad")
			}
			//fmt.Println(i % len(pub))
			//fmt.Println(pub[i%len(pub)])
			a++
		}
	*/
	/*
		for k, v := range shuffle {
			fmt.Println(k, audio[v])
			if k == 3 {
				fmt.Println("meto publicidad")
			}
			fmt.Println(i % len(pub))
			//fmt.Println(pub[i%len(pub)])
		}
		for _, _ = range shuffle2 {
			//fmt.Println(pub[v])
			//fmt.Println(pub[i%len(pub)])
		}
	*/

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
