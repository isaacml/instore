package main

import (
	//"bufio"
	"fmt"
	//"math/rand"
	//"os/exec"
	//"strings"
	//"sync"
	//"github.com/isaacml/instore/libs"
	"github.com/isaacml/instore/winamp"
	"os"
	"time"
)

func reproduccion() {
	for {
		/*
			a, p, pl := 0, 0, 1
			var gap int
			var song string
			var win winamp.Winamp
			music := make(map[int]string)
			publi := make(map[int]string)
			//Sacamos la fecha actual
			y, m, d := time.Now().Date()
		*/
		if statusProgammedMusic == "Inicial" {

		}
		time.Sleep(30 * time.Second)
	}
}

//Reproduce los mensajes automáticos de la tienda: bucle infinito que busca cada minuto un mensaje nuevo para reproducir.
func reproduccion_msgs() {
	for {
		//Obtenemos la fecha actual
		y, m, d := time.Now().Date()
		fecha := fmt.Sprintf("%4d%02d%02d", y, int(m), d)
		fecha_sql := fecha + "%"
		//Obtenemos la hora local
		hh, mm, _ := time.Now().Clock()
		clock := fmt.Sprintf("%02d:%02d", hh, mm)
		//Obtenemos todos los mensajes
		mensajes, errM := db.Query("SELECT id, fichero, playtime FROM mensaje WHERE fichero LIKE ?", fecha_sql)
		if errM != nil {
			Error.Println(errM)
		}
		for mensajes.Next() {
			var id int
			var fichero, playtime string
			//Tomamos el id, nombre y playtime de la base de datos mensaje
			err := mensajes.Scan(&id, &fichero, &playtime)
			if err != nil {
				Error.Println(err)
			}
			if playtime == clock {
				var win winamp.Winamp
				st := win.PlayFFplay(msg_files_location + fichero)
				//Si el estado de la reproduccion del mensaje = END (ha acabado), procedemos al borrado.
				if st == "END" {
					fmt.Println("Borro el fichero: "+fichero, id)
					//Borramos el fichero desde el directorio que contiene los mensajes en el player_int
					err = os.Remove(msg_files_location + fichero)
					if err != nil {
						Error.Println(err)
					}
					//Ponemos el estado de mensaje en N (ya que lo hemos borrado y no existe)
					ok, err := db.Prepare("UPDATE mensaje SET existe=? WHERE id = ?")
					if err != nil {
						Error.Println(err)
					}
					db_mu.Lock()
					_, err1 := ok.Exec("N", id)
					if err1 != nil {
						Error.Println(err1)
					}
					db_mu.Unlock()
					break
				}
			}
		}
		time.Sleep(1 * time.Minute)
	}
}
