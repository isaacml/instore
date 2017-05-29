package main

import (
	"fmt"
	//"github.com/isaacml/instore/winamp"
	"github.com/isaacml/instore/playlist"
	//"strconv"
	//"time"
)

func reproduccion() {
	var list playlist.PlayList
	//var err error
	//var gap_int int
	//Sacamos la fecha actual
	//y, m, d := time.Now().Date()
	//fecha := fmt.Sprintf("%4d%02d%02d", y, int(m), d)
	fmt.Println(list.MusicMap())
	fmt.Println(list.PubliMap())
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
		publi, errP := db.Query("SELECT fichero FROM publi WHERE fichero LIKE ?", fecha+"%")
		if errP != nil {
			Error.Println(errP)
		}
		for k, v := range list.PubliMap() {
			fmt.Println(k, audio[v])
			if gap_int-1 == k {
				fmt.Println("meto publicidad")
			}
		}
		for publi.Next() {
			var fichero string
			//Tomamos el nombre del fichero mensaje
			err := publi.Scan(&fichero)
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
