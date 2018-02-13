package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func publi_files(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		Error.Println(err)
		return
	}
	defer file.Close()
	//Formato nombre de fichero - yyyymmdd-username-filename -
	nameFileServer := r.FormValue("f_inicio") + "-" + r.FormValue("ownUser") + "-" + r.FormValue("fichero")
	//Creamos el fichero con ese formato, si ya está creado, lo machaca
	f, err := os.OpenFile(publi_files_location+nameFileServer, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		Error.Println(err)
		return
	}
	defer f.Close()
	//Proceso de copia de fichero
	_, copy_err := io.Copy(f, file)
	if copy_err != nil {
		Error.Println(copy_err)
		return
	} else {
		//Si la copia ha ido bien, pasamos a guardar los datos en la BD de servidor
		query, err := db.Query("SELECT id FROM usuarios WHERE user = ?", r.FormValue("ownUser"))
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			//Obtengo el identificador del creador
			var id int
			timestamp := time.Now().Unix()
			err = query.Scan(&id)
			if err != nil {
				Error.Println(err)
			}
			db_mu.Lock()
			_, err1 := db.Exec("INSERT INTO publi (`fichero`, `fecha_inicio`, `fecha_final`, `destino`, `creador_id`, `timestamp`, `gap`) VALUES (?,?,?,?,?,?,?)",
				nameFileServer, r.FormValue("f_inicio"), r.FormValue("f_final"), r.FormValue("destino"), id, timestamp, r.FormValue("gap"))
			db_mu.Unlock()
			if err1 != nil {
				Error.Println(err1)
			}
		}
	}
}
func msg_files(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		Error.Println(err)
		return
	}
	defer file.Close()
	//Formato nombre de fichero - yyyymmdd-username-filename -
	nameFileServer := r.FormValue("f_inicio") + "-" + r.FormValue("ownUser") + "-" + r.FormValue("fichero")
	//Creamos el fichero con ese formato, si ya está creado, lo machaca
	f, err := os.OpenFile(msg_files_location+nameFileServer, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		Error.Println(err)
		return
	}
	defer f.Close()
	//Proceso de copia de fichero
	_, copy_err := io.Copy(f, file)
	if copy_err != nil {
		Error.Println(copy_err)
		return
	} else {
		//Si la copia ha ido bien, pasamos a guardar los datos en la BD de servidor
		query, err := db.Query("SELECT id FROM usuarios WHERE user = ?", r.FormValue("ownUser"))
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			//Obtengo el identificador del creador
			var id int
			timestamp := time.Now().Unix()
			err = query.Scan(&id)
			if err != nil {
				Error.Println(err)
			}
			db_mu.Lock()
			_, err1 := db.Exec("INSERT INTO mensaje (`fichero`, `fecha_inicio`, `fecha_final`, `destino`, `creador_id`, `timestamp`, `playtime`) VALUES (?,?,?,?,?,?,?)",
				nameFileServer, r.FormValue("f_inicio"), r.FormValue("f_final"), r.FormValue("destino"), id, timestamp, r.FormValue("playtime"))
			db_mu.Unlock()
			if err1 != nil {
				Error.Println(err1)
			}
		}
	}
}

//Manda los datos de publicidad/mensajes al administrador
func modo_vista(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var output string
	var id_user string
	//Muestra los datos detallados de publicidad y mensajes
	if r.FormValue("accion") == "mostrar" {
		search := r.FormValue("search")
		err0 := db.QueryRow("SELECT id FROM usuarios WHERE user = ?", r.FormValue("username")).Scan(&id_user)
		if err0 != nil {
			Error.Println(err0)
		}
		if r.FormValue("tabla") == "publi" {
			sql := fmt.Sprintf("SELECT id, fichero, fecha_inicio, fecha_final, destino, timestamp, gap FROM %s WHERE creador_id = %s", r.FormValue("tabla"), id_user)
			query, err := db.Query(sql)
			if err != nil {
				Error.Println(err)
			}
			for query.Next() {
				var fichero, f_ini, f_fin, destinos string
				var id, timestamp, gap int64
				err = query.Scan(&id, &fichero, &f_ini, &f_fin, &destinos, &timestamp, &gap)
				if err != nil {
					Error.Println(err)
				}
				if strings.Contains(destinos, search) {
					//Se obtiene la fecha de creacion de una entidad
					f_creacion := libs.FechaCreacion(timestamp)
					output += fmt.Sprintf("<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Editar Publicidad'>%s</a>", id, fichero)
					output += fmt.Sprintf("<a href='#' onclick='borrar(%d)' title='Borrar Publicidad' style='float:right'><span class='fa fa-trash-o'></a></td>", id)
					output += fmt.Sprintf("<td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%d</td></tr>", f_ini, f_fin, destinos, f_creacion, gap)
				}
			}
		}
		if r.FormValue("tabla") == "mensaje" {
			sql := fmt.Sprintf("SELECT id, fichero, fecha_inicio, fecha_final, destino, timestamp, playtime FROM %s WHERE creador_id = %s", r.FormValue("tabla"), id_user)
			query, err := db.Query(sql)
			if err != nil {
				Error.Println(err)
			}
			for query.Next() {
				var fichero, f_ini, f_fin, destinos, playtime string
				var id, timestamp int64
				err = query.Scan(&id, &fichero, &f_ini, &f_fin, &destinos, &timestamp, &playtime)
				if err != nil {
					Error.Println(err)
				}
				if strings.Contains(destinos, search) {
					//Se obtiene la fecha de creacion de una entidad
					f_creacion := libs.FechaCreacion(timestamp)
					output += fmt.Sprintf("<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Editar Mensaje'>%s</a>", id, fichero)
					output += fmt.Sprintf("<a href='#' onclick='borrar(%d)' title='Borrar Mensaje' style='float:right'><span class='fa fa-trash-o'></a></td>", id)
					output += fmt.Sprintf("<td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", f_ini, f_fin, destinos, f_creacion, playtime)
				}
			}
		}
	}
	if r.FormValue("accion") == "borrar" {
		sql := fmt.Sprintf("DELETE FROM %s WHERE id = %s", r.FormValue("tabla"), r.FormValue("borrar"))
		db_mu.Lock()
		_, err := db.Exec(sql)
		db_mu.Unlock()
		if err != nil {
			Error.Println(err)
		}
		output = "OK"
	}
	if r.FormValue("accion") == "load" {
		var id, gap int
		var f_inicio, f_fin string
		sql := fmt.Sprintf("SELECT id, fecha_inicio, fecha_final, gap FROM %s WHERE id = %s", r.FormValue("tabla"), r.FormValue("edit_id"))
		query, err := db.Query(sql)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			err = query.Scan(&id, &f_inicio, &f_fin, &gap)
			if err != nil {
				Error.Println(err)
			}
			fmt.Fprintf(w, "id=%d&f_inicio=%s&f_fin=%s&gap=%d", id, f_inicio, f_fin, gap)
		}
	}
	fmt.Fprint(w, output) //fmt.Println(output)
}
