package main

import (
	"fmt"
	//"github.com/isaacml/instore/winamp"
	"bufio"
	"math/rand"
	"os/exec"
	"strings"
	"time"
)

func reproduccion() {
	audio := make(map[int]string)
	pub := make(map[int]string)
	i := 0
	/*
		var win winamp.Winamp
		win.RunWinamp()
		win.Load("c:/instore/musica/ACDC.mp3")
		win.Play()
		time.Sleep(15 * time.Second)
		fmt.Println(win.SongLenght())
		fmt.Println(win.SongPlay())
		win.PlayFFplay("c:/instore/intros/Mysterious.mp3")
		fmt.Println(win.Status().FFplaying)
	*/
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
		audio[i] = strings.TrimRight(line, "\n")
		i++
	}
	cmd.Wait()
	rand.Seed(time.Now().UnixNano())
	shuffle := rand.Perm(len(audio))
	fmt.Println(shuffle)
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
		pub[i] = strings.TrimRight(line, "\n")
		i++
	}
	cmd.Wait()
	shuffle2 := rand.Perm(len(pub))
	fmt.Println(shuffle2)
	//creamos la playlist mezclando audio + pub con el gap correspondiente
	//a, p, i := 0, 0, 0

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

}
