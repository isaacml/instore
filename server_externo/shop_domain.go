package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
	"strings"
)

//Funcion que toma la peticion de dominio por parte de la tienda.
func recoger_dominio(w http.ResponseWriter, r *http.Request) {
	var output string
	var domains []string
	fmt.Println(r.FormValue("dominio"))
	doms := strings.Split(r.FormValue("dominio"), ":.:")
	fmt.Println(len(doms))
	doms = doms[:len(doms)-1]
	for _, val := range doms {
		//enviamos el dominio a DomainGenerator() en la libreria de funciones
		dom := libs.DomainGenerator(val)
		for _, v := range dom {
			domains = append(domains, v)
		}
	}
	//Borramos dominios duplicados
	domains = libs.RemoveDuplicates(domains)
	//Nuestra fecha actual personalizada
	fecha := libs.MyCurrentDate()
	//Formamos una cadena con la publicidad y los mensajes para ese dominio y con la fecha que le corresponde
	output += "[publi]"
	for _, val := range domains {
		publicidad, err := db.Query("SELECT fichero, fecha_inicio, gap FROM publi WHERE destino = ? AND fecha_inicio = ?", val, fecha)
		if err != nil {
			Error.Println(err)
		}
		for publicidad.Next() {
			var f_publi, fecha_ini, gap string
			err = publicidad.Scan(&f_publi, &fecha_ini, &gap)
			if err != nil {
				Error.Println(err)
			}
			output += ";" + f_publi + "<=>" + fecha_ini + "<=>" + gap
		}
	}
	output += "[mensaje]"
	for _, val := range domains {
		mensajes, err := db.Query("SELECT fichero, playtime FROM mensaje WHERE destino = ? AND fecha_inicio = ?", val, fecha)
		if err != nil {
			Error.Println(err)
		}
		for mensajes.Next() {
			var f_msg, playtime string
			err = mensajes.Scan(&f_msg, &playtime)
			if err != nil {
				Error.Println(err)
			}
			output += ";" + f_msg + "<=>" + playtime
		}
	}
	//Enviamos la cadena
	fmt.Fprint(w, output) //fmt.Println(output)
}
