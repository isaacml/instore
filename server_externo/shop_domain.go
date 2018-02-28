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
	doms := strings.Split(r.FormValue("dominios"), ":.:")
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
		publicidad, err := db.Query("SELECT fichero, fecha_inicio, fecha_final, gap FROM publi WHERE destino = ?", val)
		if err != nil {
			Error.Println(err)
		}
		for publicidad.Next() {
			var f_publi, fecha_ini, fecha_fin, gap string
			err = publicidad.Scan(&f_publi, &fecha_ini, &fecha_fin, &gap)
			if err != nil {
				Error.Println(err)
			}
			fmt.Println("FECHAS: ", fecha_ini, fecha, fecha_fin)
			//BETWEEN
			if fecha_ini <= fecha && fecha_fin >= fecha {
				output += ";" + f_publi + "<=>" + fecha_ini + "<=>" + fecha_fin + "<=>" + gap
			}
		}
	}
	output += "[mensaje]"
	for _, val := range domains {
		mensajes, err := db.Query("SELECT fichero, fecha_inicio, fecha_final, playtime FROM mensaje WHERE destino = ?", val)
		if err != nil {
			Error.Println(err)
		}
		for mensajes.Next() {
			var f_msg, fecha_ini, fecha_fin, playtime string
			err = mensajes.Scan(&f_msg, &fecha_ini, &fecha_fin, &playtime)
			if err != nil {
				Error.Println(err)
			}
			//BETWEEN
			if fecha_ini <= fecha && fecha_fin >= fecha {
				output += ";" + f_msg + "<=>" + fecha_ini + "<=>" + fecha_fin + "<=>" + playtime
			}
		}
	}
	//Enviamos la cadena
	fmt.Fprint(w, output) //fmt.Println(output)
}
