package winamp

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

var (
	volMax int = 250
)

type Winamp struct {
	// internal status variables
	play   bool
	stop   bool
	pause  bool
	volume int
	run    bool
	mu     sync.Mutex // mutex tu protect the internal variables on multithreads
}

type Status struct {
	Playing  bool
	Stopping bool
	Pausing  bool
}

//Constructor para Winamp
func Winamper() *Winamp {
	win := &Winamp{}
	win.mu.Lock()
	defer win.mu.Unlock()

	// initialize the internal variables values
	win.play = false
	win.stop = false
	win.pause = false
	win.volume = 0
	win.run = false

	return win
}

//Constructor para Status
func (w *Winamp) Status() *Status {
	var st Status

	w.mu.Lock()
	defer w.mu.Unlock()

	st.Playing = w.play
	st.Pausing = w.pause
	st.Stopping = w.stop

	return &st
}

//Función que arranca Winamp, si no está arrancado y establece el volumen a 250
func (w *Winamp) RunWinamp() {
	if w.run == false {
		exec.Command("cmd", "/c", "%winamp%").Start()
		w.mu.Lock()
		w.volume = volMax
		w.run = true
		w.mu.Unlock()
		
	}
}
//Función que establece el volumen del Winamp
func (w *Winamp) Volume() {
	fmt.Println("este es el volumen", w.volume)
	vol := fmt.Sprintf("C:\\instore\\Winamp\\CLEvER.exe volume %d", w.volume)
	exec.Command("cmd", "/c", vol).Run()
}
//Función que cierra Winamp
func (w *Winamp) WinampClose() {
	w.mu.Lock()
	w.run = false
	w.mu.Unlock()
	exec.Command("cmd", "/c", "taskkill /IM winamp.exe").Run()
}

//Función que comprueba si Winamp se está ejecutando o no
func (w *Winamp) WinampIsOpen() bool {
	var gen_bat string
	//Creamos el fichero bat que va a guardar la duracion total(en seg) de la canción
	isOpenFile, err := os.Create("C:\\instore\\Winamp\\isOpenWin.bat")
	if err != nil {
		err = fmt.Errorf("BAT: CANNOT CREATE BAT FILE")
	}
	gen_bat = "@echo off\r\ntasklist /fi \"IMAGENAME eq winamp.exe\" | find /i \"winamp.exe\" > nul\r\nif not errorlevel 1 (echo Existe) else (echo NoExiste)"
	isOpenFile.WriteString(gen_bat)
	bat, err := exec.Command("cmd", "/c", "C:\\instore\\Winamp\\isOpenWin.bat").CombinedOutput()
	if err != nil {
		err = fmt.Errorf("bat: CANNOT OPEN BAT")
	}
	limpio := strings.TrimRight(string(bat), "\r\n")
	//Evaluamos la salida del fichero bat
	if limpio == "Existe" {
		w.run = true
	} else {
		w.run = false
	}
	return w.run
}

//Si Winamp está arrancado, carga una playlist
func (w *Winamp) Load(file string) error {
	var err error
	if w.run == true {
		var gen_fich string
		bat, err := os.Create("song.bat")
		if err != nil {
			err = fmt.Errorf("bat: CANNOT CREATE BAT FILE")
		}
		defer bat.Close()
		gen_fich = "@echo off\r\nC:\\instore\\Winamp\\CLEvER.exe loadnew " + file
		bat.WriteString(gen_fich)
		err = exec.Command("cmd", "/c", "song.bat").Run()
		if err != nil {
			err = fmt.Errorf("load: CANNOT_LOAD_PLAYLIST")
		}
	} else {
		err = fmt.Errorf("winamp: WINAMP_IS_NOT_RUNNING")
	}
	return err
}
func (w *Winamp) Play() {
	w.mu.Lock()
	w.play = true
	w.pause = false
	w.stop = false
	w.mu.Unlock()
	exec.Command("cmd", "/c", "C:\\instore\\Winamp\\CLEvER.exe play").Run()
}
func (w *Winamp) Stop() {
	w.mu.Lock()
	w.play = false
	w.pause = false
	w.stop = true
	w.mu.Unlock()
	exec.Command("cmd", "/c", "C:\\instore\\Winamp\\CLEvER.exe stop").Run()
}
func (w *Winamp) Pause() {
	w.mu.Lock()
	w.play = false
	w.pause = true
	w.stop = false
	w.mu.Unlock()
	exec.Command("cmd", "/c", "C:\\instore\\Winamp\\CLEvER.exe  pause").Run()
}

//Muestra el tiempo de reproducción (en seg) de la canción
func (w *Winamp) SongPlay() int {
	var min, sec int
	lector, _ := exec.Command("cmd", "/c", "C:\\instore\\Winamp\\CLEvER.exe position").CombinedOutput()
	//formato de timeplay -> 02:12
	timeplay := strings.Split(fmt.Sprintf("%s", string(lector)), ":")
	//hago un split para sacar los minutos y los segundos
	min, _ = strconv.Atoi(timeplay[0])
	sec, _ = strconv.Atoi(timeplay[1])
	totalsec := (min * 60) + sec

	return totalsec
}

//Muestra el tiempo(en seg) que queda para acabar la cancion
func (w *Winamp) SongEnd() int {
	var min, sec int
	lector, _ := exec.Command("cmd", "/c", "C:\\instore\\Winamp\\CLEvER.exe timeleft").CombinedOutput()
	//formato de timend -> 05:02
	timend := strings.Split(fmt.Sprintf("%s", string(lector)), ":")
	//hago un split para sacar los minutos y los segundos
	min, _ = strconv.Atoi(timend[0])
	sec, _ = strconv.Atoi(timend[1])
	totalsec := (min * 60) + sec

	return totalsec
}
func (w *Winamp) VolumeUp() {
	var cont int
	for i := 1; i <= volMax; i += 25 {
		exec.Command("cmd", "/c", "C:\\instore\\Winamp\\CLEvER.exe volup").Run()
		cont++
	}
	if w.volume >= volMax {
		w.mu.Lock()
		w.volume = volMax
		w.mu.Unlock()
		return
	} else {
		w.mu.Lock()
		w.volume = (w.volume + (cont * 4))
		w.mu.Unlock()
	}
	fmt.Println(w.volume)
}
func (w *Winamp) VolumeDown() {
	var cont int
	for i := 1; i <= volMax; i += 25 {
		exec.Command("cmd", "/c", "C:\\instore\\Winamp\\CLEvER.exe voldn").Run()
		cont++
	}
	if w.volume < 10 {
		w.mu.Lock()
		w.volume = 0
		w.mu.Unlock()
		return
	} else {
		w.mu.Lock()
		w.volume = (w.volume - (cont * 4))
		w.mu.Unlock()
	}
	fmt.Println(w.volume)
}
//Limpia la playlist
func (w *Winamp) Clear(){
	exec.Command("cmd", "/c", "C:\\instore\\Winamp\\CLEvER.exe clear").Run()
}
//Tiempo total de un fichero de musica en segundos
func (w *Winamp) SongLenght(file string) int {
	var gen_bat string
	//Creamos el fichero bat que va a guardar la duracion total(en seg) de la canción
	song_lenght_bat, err := os.Create("song_lenght.bat")
	if err != nil {
		err = fmt.Errorf("lenght_bat: CANNOT CREATE BAT FILE")
	}
	defer song_lenght_bat.Close()
	gen_bat = "@echo off\r\nC:\\instore\\ffprobe.exe -v quiet -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 \"" + file + "\""
	song_lenght_bat.WriteString(gen_bat)
	//Una vez creado el fichero, lo ejecutamos y tomamos su salida
	seg, _ := exec.Command("cmd", "/c", "song_lenght.bat").CombinedOutput()
	//formato de SongLenght -> 201.91234 seg
	song := strings.Split(fmt.Sprintf("%s", string(seg)), ".")
	total_sec, err := strconv.Atoi(song[0])
	if err != nil {
		err = fmt.Errorf("conv: CANNOT_CONVERSION")
	}
	return total_sec
}

// Metodo que introduce la publicidad por ffplay
func (w *Winamp) PlayFFplay(publi string) string {
	var st string
	//Bajo el volumen del reproductor Winamp a 0
	exec.Command("cmd", "/c", "C:\\instore\\Winamp\\CLEvER.exe volume 0").Run()
	//Reproduzco la publicidad del ffplay
	play := fmt.Sprintf("C:\\instore\\ffplay.exe -nodisp %s -autoexit", publi)
	exec.Command("cmd", "/c", play).Run()
	//Vuelvo a subir el volumen a como estaba
	inc := fmt.Sprintf("C:\\instore\\Winamp\\CLEvER.exe volume %d", volMax)
	exec.Command("cmd", "/c", inc).Run()
	st = "END"
	return st
}
