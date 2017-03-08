package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
)

func bitmap_actions(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	query, err := db.Query("SELECT bitmap_acciones FROM usuarios WHERE user = ?", r.FormValue("user"))
	if err != nil {
		Error.Println(err)
	}
	for query.Next() {
		var bitmap_hex string
		err = query.Scan(&bitmap_hex)
		if err != nil {
			Error.Println(err)
		}
		prog_pub := libs.BitmapParsing(bitmap_hex, PROG_PUB)     //res[0]
		prog_mus := libs.BitmapParsing(bitmap_hex, PROG_MUS)     //res[1]
		prog_msg := libs.BitmapParsing(bitmap_hex, PROG_MSG)     //res[2]
		add_mus := libs.BitmapParsing(bitmap_hex, ADD_MUS)       //res[3]
		msg_auto := libs.BitmapParsing(bitmap_hex, MSG_AUTO)     //res[4]
		msg_normal := libs.BitmapParsing(bitmap_hex, MSG_NORMAL) //res[5]
		//Pasamos los valores al html
		fmt.Fprintf(w, "%d;%d;%d;%d;%d;%d", prog_pub, prog_mus, prog_msg, add_mus, msg_auto, msg_normal)
	}
}

//Al hacer click para editar un usuario, esta función va a determinar si a un usuario le pertenece o no una acción
func bitmap_checked(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	edit_id := r.FormValue("edit_id")

	var output string
	query, err := db.Query("SELECT bitmap_acciones FROM usuarios WHERE id = ?", edit_id)
	if err != nil {
		Error.Println(err)
	}
	for query.Next() {
		var bitmap string
		err = query.Scan(&bitmap)
		if err != nil {
			Error.Println(err)
		}
		//Checkeado o No, segun el resultado al pasarle la mascara
		prog_pub := libs.BitmapParsing(bitmap, PROG_PUB) //res[0]
		if prog_pub != 0 {
			output += "<input type='checkbox' name='prog_pub_edit' value='1' checked/> Programar Publicidad<br>"
		} else {
			output += "<input type='checkbox' name='prog_pub_edit' value='1'/> Programar Publicidad<br>"
		}
		prog_mus := libs.BitmapParsing(bitmap, PROG_MUS) //res[1]
		if prog_mus != 0 {
			output += "<input type='checkbox' name='prog_mus_edit' value='2' checked/> Programar Música<br>"
		} else {
			output += "<input type='checkbox' name='prog_mus_edit' value='2'/> Programar Música<br>"
		}
		prog_msg := libs.BitmapParsing(bitmap, PROG_MSG) //res[2]
		if prog_msg != 0 {
			output += "<input type='checkbox' name='prog_msg_edit' value='4' checked/> Programar Mensajes Nuevos<br>"
		} else {
			output += "<input type='checkbox' name='prog_msg_edit' value='4'/> Programar Mensajes Nuevos<br>"
		}
		add_mus := libs.BitmapParsing(bitmap, ADD_MUS) //res[3]
		if add_mus != 0 {
			output += "<input type='checkbox' name='add_mus_edit' value='8' checked/> Añadir Música No Cifrada<br>"
		} else {
			output += "<input type='checkbox' name='add_mus_edit' value='8'/> Añadir Música No Cifrada<br>"
		}
		msg_auto := libs.BitmapParsing(bitmap, MSG_AUTO) //res[4]
		if msg_auto != 0 {
			output += "<input type='checkbox' name='msg_auto_edit' value='16' checked/> Reproducir Mensajes Automatizados<br>"
		} else {
			output += "<input type='checkbox' name='msg_auto_edit' value='16'/> Reproducir Mensajes Automatizados<br>"
		}
		msg_normal := libs.BitmapParsing(bitmap, MSG_NORMAL) //res[5]
		if msg_normal != 0 {
			output += "<input type='checkbox' name='msg_normal_edit' value='32' checked/> Reproducir Mensajes Normales<br>"
		} else {
			output += "<input type='checkbox' name='msg_normal_edit' value='32'/> Reproducir Mensajes Normales<br>"
		}
		fmt.Fprint(w, output)
	}
}
