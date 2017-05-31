package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	//"sync"
	"github.com/isaacml/instore/winamp"
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
		music[i] = strings.TrimRight(line, "\n")
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
	a := 1
	for _, v := range shuffle {
		var win winamp.Winamp
		win.RunWinamp()
		win.Load(music[v])
		win.Play()
		time.Sleep(15 * time.Second)
		//fmt.Println(win.SongLenght())
		//fmt.Println(win.SongPlay())
		if a == gap {
			fmt.Println("Meto publicidad")
			a = 0
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
