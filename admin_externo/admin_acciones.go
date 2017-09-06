package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
	"time"
)

//Función que actualiza el tiempo de expiración
func updateExpires(sid string) {
	expires := time.Now().Unix() + int64(session_timeout)
	tiempo[sid] = expires
}
//Acciones independientes
func acciones(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	accion := r.FormValue("accion")
	_, ok := user[sid]
	if ok {
		//Pasamos el nombre de usuario al servidor para obtener los bitmaps de acciones
		if accion == "bitmap_perm" {
			respuesta := libs.GenerateFORM(serverext["serverroot"]+"/acciones.cgi", "accion;bitmap_perm", "user;"+username)
			fmt.Fprint(w, respuesta)
		}
		//Indica si mostrar o no el mantenimiento de las organizaciones
		if accion == "show_org" {
			respuesta := libs.GenerateFORM(serverext["serverroot"]+"/acciones.cgi", "accion;show_org", "user;"+username)
			fmt.Fprint(w, respuesta)
		}
		//Envia un select para mostrar las horas en el panel de mensajes.html
		if accion == "horas_msg" {
			var str string
			for i := 1; i <= 24; i++ {
				str += fmt.Sprintf("<option value='%02d'>%02d</option>", i, i)
			}
			fmt.Fprint(w, str)
		}
		//Envia un select para mostrar los minutos en el panel de mensajes.html
		if accion == "minutos_msg" {
			var str string
			for i := 0; i <= 59; i++ {
				str += fmt.Sprintf("<option value='%02d'>%02d</option>", i, i)
			}
			fmt.Fprint(w, str)
		}
	}
}