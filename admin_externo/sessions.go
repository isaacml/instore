package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"github.com/todostreaming/realip"
	"math/rand"
	"net/http"
	"time"
)

//Mapas de control de sesiones
var (
	user   map[string]string = make(map[string]string)
	ip     map[string]string = make(map[string]string)
	tiempo map[string]int64  = make(map[string]int64)
	level  map[string]int    = make(map[string]int)
	agent  map[string]string = make(map[string]string)
)

func controlinternalsessions() {
	for {
		for k, v := range tiempo {
			//fmt.Println(time.Now().Unix())
			if time.Now().Unix() > v {
				delete(user, k)
				delete(ip, k)
				delete(tiempo, k)
				delete(level, k)
				delete(agent, k)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

// genera una session id o Value del Cookie aleatoria y de la longitud que se quiera
func sessionid(r *rand.Rand, n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[r.Intn(len(letterRunes))]
	}
	return string(b)
}

// funcion q tramita el login correcto o erroneo
func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	libs.LoadSettingsLin(serverRoot, settings)
	var agente string
	username = r.FormValue(name_username)
	password := r.FormValue(name_password)
	aleat := rand.New(rand.NewSource(time.Now().UnixNano()))
	sid := sessionid(aleat, session_value_len)
	expiration := time.Now().Unix() + int64(session_timeout)
	ipReal := realip.RealIP(r) // Direccion IP real
	for k, v := range r.Header {
		if k == "User-Agent" {
			agente = v[0]
		}
	}
	//SE PASAN LAS VARIABLES POST AL SERVIDOR EXTERNO PARA LA AUTENTICACION
	respuesta := fmt.Sprintf("%s", libs.GenerateFORM(settings["serverroot"]+"/login.cgi", "user;"+username, "pass;"+password))
	//RECOGEMOS LA RESPUESTA
	if respuesta == "OK" {
		//Cuando se repite autenticacion de usuario
		for key, _ := range user {
			//Si el usuario e IP existen pero el navegador es distinto, se abre una nueva sesion para el
			if username == user[key] && ipReal == ip[key] && agente != agent[key] {
				user[sid] = username
				ip[sid] = realip.RealIP(r)
				tiempo[sid] = expiration
				agent[sid] = agente
				http.Redirect(w, r, "/"+enter_page+"?"+sid, http.StatusFound)
				return
			}
			//Si el usuario, IP y navegador existen, se actualizan los mapas
			if username == user[key] && ipReal == ip[key] && agente == agent[key] {
				user[key] = username
				ip[key] = realip.RealIP(r)
				tiempo[key] = expiration
				agent[key] = agente
				http.Redirect(w, r, "/"+enter_page+"?"+key, http.StatusFound)
				return
			}
		}
		// Guardamos constancia de la session en nuestros mapas internos, si es la primera vez que se autentica
		user[sid] = username
		ip[sid] = ipReal
		tiempo[sid] = expiration
		level[sid] = 0
		agent[sid] = agente
		// Enviamos a la pagina de entrada tras el login correcto
		http.Redirect(w, r, "/"+enter_page+"?"+sid, http.StatusFound)
	} else {
		// Te devolvemos a la pagina inicial de login
		http.Redirect(w, r, "/"+first_page+".html?err", http.StatusFound)
	}
}

// función que tramita el logout de la session
func logout(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	_, ok := user[sid]
	if ok {
		for k, _ := range user {
			delete(user, k)
			delete(ip, k)
			delete(tiempo, k)
			delete(level, k)
			delete(agent, k)
		}
		http.Redirect(w, r, "/"+first_page+".html", http.StatusFound)
	}
}