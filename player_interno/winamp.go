package main

import (
	"fmt"
	"github.com/isaacml/instore/winamp"
	"math/rand"
	"net/http"
	"time"
)

func reproduccion(w http.ResponseWriter, r *http.Request) {
	audio := make(map[int]string)

	var win winamp.Winamp
	win.RunWinamp()
	win.Load("c:/instore/musica/ACDC.mp3")
	win.Play()
	time.Sleep(15 * time.Second)
	fmt.Println(win.SongLenght())
	fmt.Println(win.SongPlay())
	win.PlayFFplay("c:/instore/intros/Mysterious.mp3")
	fmt.Println(win.Status().FFplaying)
	shuffle := rand.Perm(len(audio))
	fmt.Println(shuffle)
}
