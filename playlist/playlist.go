package playlist

import (
	"bufio"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var (
	music_files string = "C:\\instore\\Music\\"
	publi_files string = "C:\\instore\\PubliShop\\"
	i           int    = 0
)

type PlayList struct {
	// internal status variables
	music map[int]string
	publi map[int]string
	mu    sync.Mutex // mutex tu protect the internal variables on multithreads
}

//Constructor para Winamp
func Listado() *PlayList {
	list := &PlayList{}
	list.mu.Lock()
	defer list.mu.Unlock()

	// initialize the internal variables values
	list.music = make(map[int]string)
	list.publi = make(map[int]string)
	return list
}

//Función que cierra Winamp
func (p *PlayList) MusicMap() []int {
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
		fmt.Printf("%s", line)
		//p.music[i] = strings.TrimRight(line, "\n")
		i++
	}
	cmd.Wait()
	rand.Seed(time.Now().UnixNano())
	return rand.Perm(len(p.music))
}

//Función que cierra Winamp
func (p *PlayList) PubliMap() []int {
	cmd := exec.Command("cmd", "/c", "dir /s /b "+publi_files+"*.mp3")
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
		p.publi[i] = strings.TrimRight(line, "\n")
		i++
	}
	cmd.Wait()
	return rand.Perm(len(p.publi))
}
