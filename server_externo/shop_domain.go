package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
	"strings"
	"time"
)

//Funcion que toma la peticion de dominio por parte de la tienda.
func recoger_dominio(w http.ResponseWriter, r *http.Request) {
	var output string
	var domains []string
	doms := strings.Split(r.FormValue("dominio"), ":.:")
	doms = doms[:len(doms)-1]
	for _, val := range doms {
		//enviamos el dominio a DomainGenerator() en la libreria de funciones
		dom := libs.DomainGenerator(val)
		for _, v := range dom {
			fmt.Println("val1", v)
			for _, valor := range domains {
				fmt.Println("val2", v)
				if v != valor {
					domains = append(domains, v)
				}
			}
		}
	}
	fmt.Println(domains)
	fecha_actual := time.Now()
	//Formato de la fecha actual --> 20070405
	string_fecha := fmt.Sprintf("%4d%02d%02d", fecha_actual.Year(), int(fecha_actual.Month()), fecha_actual.Day())
	//Formamos una cadena con la publicidad y los mensajes para ese dominio y con la fecha que le corresponde
	output += "[publi]"
	for _, val := range domains {
		publicidad, err := db.Query("SELECT fichero, fecha_inicio, gap FROM publi WHERE destino = ? AND fecha_inicio = ?", val, string_fecha)
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
		mensajes, err := db.Query("SELECT fichero, playtime FROM mensaje WHERE destino = ? AND fecha_inicio = ?", val, string_fecha)
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
	fmt.Fprint(w, output)
}
