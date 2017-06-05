package main

import (
	"github.com/isaacml/instore/winamp"
	"time"
	"os/exec"
)

func main() {
	var win winamp.Winamp
	win.RunWinamp()
	
	exec.Command("cmd", "/c", "C:\\instore\\pruebas\\prueba.bat").Run()
	win.Load(`"C:\instore\Music\Rock\System Of A Down - Chop suey.mp3"`)
	time.Sleep(3 * time.Second)
	win.Play()
	/*
	win.Load("\"C:\\instore\\Music\\Rock\\System Of A Down - Chop suey.mp3\"")
	fmt.Println("\"C:\\instore\\Music\\Rock\\System Of A Down - Chop suey.mp3\"")
	win.Play()
	*/
	time.Sleep(15 * time.Second)
}

