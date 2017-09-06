package main

import (
	"io"
	"net/http"
	"os"
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
