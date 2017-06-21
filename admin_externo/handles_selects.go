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

//Función que muestra el usuario en activo
func user_admin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		updateExpires(sid)
		fmt.Fprint(w, username)
	}
}

//Envia un select del tipo de entidad(entidad_id) a la que pertenece el usuario (ROOT=0 o Normal=other)
func user_entidad(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/user_entidad.cgi", "username;"+username))
		fmt.Fprint(w, respuesta)
	}
}

//Envia un select para mostrar las entidades por usuario
func almacen_entidad(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/almacen_entidad.cgi", "username;"+username))
		fmt.Fprint(w, respuesta)
	}
}

//Envia un select para mostrar los almacenes por usuario
func pais_almacen(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/pais_almacen.cgi", "username;"+username))
		fmt.Fprint(w, respuesta)
	}
}

//Envia un select para mostrar los paises por usuario
func region_pais(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/region_pais.cgi", "username;"+username))
		fmt.Fprint(w, respuesta)
	}
}

//Envia un select para mostrar las regiones por usuario
func provincia_region(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/provincia_region.cgi", "username;"+username))
		fmt.Fprint(w, respuesta)
	}
}

//Envia un select para mostrar las provincias por usuario
func tienda_provincia(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		loadSettings(serverRoot)
		updateExpires(sid)
		respuesta := fmt.Sprintf("%s", libs.GenerateFORM(serverext["serverroot"]+"/tienda_provincia.cgi", "username;"+username))
		fmt.Fprint(w, respuesta)
	}
}

//Envia un select para mostrar las horas en el panel de mensajes.html
func horas_msg(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		var str string
		for i := 1; i <= 24; i++ {
			str += fmt.Sprintf("<option value='%02d'>%02d</option>", i, i)
		}
		fmt.Fprint(w, str)
	}
}

//Envia un select para mostrar los minutos en el panel de mensajes.html
func minutos_msg(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		var str string
		for i := 1; i <= 59; i++ {
			str += fmt.Sprintf("<option value='%02d'>%02d</option>", i, i)
		}
		fmt.Fprint(w, str)
	}

}
