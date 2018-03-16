package winamp

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	//"time"
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

//Función que arranca Winamp, si no está arrancado
func (w *Winamp) RunWinamp() {
	if w.run == false {
		exec.Command("cmd", "/c", "start /MIN winamp\\winamp.exe").Start()
		w.mu.Lock()
		w.run = true
		w.mu.Unlock()

	}
}

//Función que establece el volumen del Winamp
func (w *Winamp) Volume() {
	err := exec.Command("cmd", "/c", "apps\\CLEvER.exe volume 250").Run()
	if err != nil {
		err = fmt.Errorf("VOL: FAIL TO SEND VOLUME")
	}
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
	isOpenFile, err := os.Create("winamp\\isOpenWin.bat")
	if err != nil {
		err = fmt.Errorf("BAT: CANNOT CREATE BAT FILE")
	}
	defer isOpenFile.Close()
	gen_bat = "@echo off\r\ntasklist /fi \"IMAGENAME eq winamp.exe\" | find /i \"winamp.exe\" > nul\r\nif not errorlevel 1 (echo Existe) else (echo NoExiste)"
	isOpenFile.WriteString(gen_bat)
	bat, err := exec.Command("cmd", "/c", "winamp\\isOpenWin.bat").CombinedOutput()
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
	var bat *os.File
	if w.run == true {
		var gen_fich string
		bat, err = os.Create("song.bat")
		if err != nil {
			err = fmt.Errorf("bat: CANNOT CREATE BAT FILE")
		}
		defer bat.Close()
		gen_fich = "@echo off\r\napps\\CLEvER.exe loadnew " + file
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
	exec.Command("cmd", "/c", "apps\\CLEvER.exe play").Run()
}
func (w *Winamp) Stop() {
	w.mu.Lock()
	w.play = false
	w.pause = false
	w.stop = true
	w.mu.Unlock()
	exec.Command("cmd", "/c", "apps\\CLEvER.exe stop").Run()
}
func (w *Winamp) Pause() {
	w.mu.Lock()
	w.play = false
	w.pause = true
	w.stop = false
	w.mu.Unlock()
	exec.Command("cmd", "/c", "apps\\CLEvER.exe  pause").Run()
}

//Muestra el tiempo de reproducción (en seg) de la canción
func (w *Winamp) SongPlay() int {
	var min, sec int
	lector, _ := exec.Command("cmd", "/c", "apps\\CLEvER.exe position").CombinedOutput()
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
	lector, _ := exec.Command("cmd", "/c", "apps\\CLEvER.exe timeleft").CombinedOutput()
	//formato de timend -> 05:02
	timend := strings.Split(fmt.Sprintf("%s", string(lector)), ":")
	//hago un split para sacar los minutos y los segundos
	min, _ = strconv.Atoi(timend[0])
	sec, _ = strconv.Atoi(timend[1])
	totalsec := (min * 60) + sec

	return totalsec
}

//Limpia la playlist
func (w *Winamp) Clear() {
	exec.Command("cmd", "/c", "apps\\CLEvER.exe clear").Run()
}

//Tiempo total de un fichero de musica en segundos
func (w *Winamp) SongLenght(file string) int {
	var gen_bat string
	var total_sec, final_segs int
	//Creamos el fichero bat que va a guardar la duracion total(en seg) de la canción
	song_lenght_bat, err := os.Create("song_lenght.bat")
	if err != nil {
		err = fmt.Errorf("lenght_bat: CANNOT CREATE BAT FILE")
	}
	defer song_lenght_bat.Close()
	gen_bat = "@echo off\r\napps\\ffprobe.exe -v quiet -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 \"" + file + "\""
	song_lenght_bat.WriteString(gen_bat)
	//Una vez creado el fichero, lo ejecutamos y tomamos su salida
	seg, _ := exec.Command("cmd", "/c", "song_lenght.bat").CombinedOutput()
	if strings.TrimSpace(string(seg)) == "NA" {
		st, err := os.Stat(file)
		if err != nil {
			err = fmt.Errorf("Statlenght: ERROR_STAT_FILE")
		}
		tamanio := st.Size()
		total_sec = int(tamanio / 16000)
	} else {
		//formato de SongLenght -> 201.91234 seg
		song := strings.Split(fmt.Sprintf("%s", string(seg)), ".")
		total_sec, err = strconv.Atoi(song[0])
		if err != nil {
			err = fmt.Errorf("conv: CANNOT_CONVERSION")
		}
	}
	if total_sec > 30 && total_sec < 300 {
		final_segs = total_sec
	} else if total_sec < 30 {
		final_segs = 30
	} else if total_sec > 300 {
		final_segs = 300
	}
	return final_segs
}

// Metodo que introduce la publicidad por ffplay
func (w *Winamp) PlayFFplay(publi string) string {
	var gen_bat string
	//Paramos la cancion
	exec.Command("cmd", "/c", "apps\\CLEvER.exe volume 0").Run()
	//Creamos el fichero bat que va a guardar la duracion total(en seg) de la canción
	msg_file, err := os.Create("msg_file.bat")
	if err != nil {
		err = fmt.Errorf("msg_file: CANNOT CREATE MSG FILE")
	}
	defer msg_file.Close()
	gen_bat = "@echo off\r\napps\\ffplay.exe -nodisp \"" + publi + "\" -autoexit"
	msg_file.WriteString(gen_bat)
	//Una vez creado el fichero, lo ejecutamos (se reproduce el mensaje)
	exec.Command("cmd", "/c", "msg_file.bat").Run()
	//Vuelve a sonar la cancion
	exec.Command("cmd", "/c", "apps\\CLEvER.exe volume 250").Run()
	return "END"
}
