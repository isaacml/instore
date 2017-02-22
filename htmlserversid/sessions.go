package main

import (
	"fmt"
	"github.com/todostreaming/realip"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// mapas de control de sessions
var user map[string]string = make(map[string]string)
var ip map[string]string = make(map[string]string)
var tiempo map[string]int64 = make(map[string]int64)
var level map[string]int = make(map[string]int)
var agent map[string]string = make(map[string]string)

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

	usuario := r.FormValue(name_username)
	pass := r.FormValue(name_password)

	if authentication(usuario, pass) {
		var agente string
		// Generamos la Cookie a escibir en el navegador del usuario
		aleat := rand.New(rand.NewSource(time.Now().UnixNano()))
		sid := sessionid(aleat, session_value_len)
		expiration := time.Now().Unix() + int64(session_timeout)
		ipReal := realip.RealIP(r) // Direccion IP real
		for k, v := range r.Header {
			if k == "User-Agent" {
				agente = v[0]
			}
		}
		//Cuando se repite autenticacion de usuario
		for key, _ := range user {
			//Si el usuario e IP existen pero el navegador es distinto, se abre una nueva sesion para el
			if usuario == user[key] && ipReal == ip[key] && agente != agent[key] {
				user[sid] = usuario
				ip[sid] = realip.RealIP(r)
				tiempo[sid] = expiration
				agent[sid] = agente
				http.Redirect(w, r, "/"+enter_page+"?"+sid, http.StatusFound)
				return
			}
			//Si el usuario, IP y navegador existen, se actualizan los mapas
			if usuario == user[key] && ipReal == ip[key] && agente == agent[key] {
				user[key] = usuario
				ip[key] = realip.RealIP(r)
				tiempo[key] = expiration
				agent[key] = agente
				http.Redirect(w, r, "/"+enter_page+"?"+key, http.StatusFound)
				return
			}
		}
		// Guardamos constancia de la session en nuestros mapas internos, si es la primera vez que se autentica
		user[sid] = usuario
		ip[sid] = ipReal
		tiempo[sid] = expiration
		level[sid] = 5
		agent[sid] = agente

		// Enviamos a la pagina de entrada tras el login correcto
		http.Redirect(w, r, "/"+enter_page+"?"+sid, http.StatusFound)
	} else {
		// Te devolvemos a la pagina inicial de login
		fmt.Println("Login incorrecto")
		http.Redirect(w, r, "/"+first_page+".html?err", http.StatusFound)
	}
}

// función que tramita el logout de la session
func logout(w http.ResponseWriter, r *http.Request) {
	if checkCGI(r) == true {
		for k, _ := range user {
			delete(user, k)
			delete(ip, k)
			delete(tiempo, k)
			delete(level, k)
			delete(agent, k)
		}
		http.Redirect(w, r, "/"+first_page+".html", http.StatusFound)
	} else {
		http.Redirect(w, r, "/"+first_page+".html", http.StatusFound)
	}
}

// convierte un string numérico en un entero int
func toInt(cant string) (res int) {
	res, _ = strconv.Atoi(cant)
	return
}
