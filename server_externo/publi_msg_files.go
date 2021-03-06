package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/isaacml/instore/libs"
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
		db_mu.Lock()
		query, err := db.Query("SELECT id FROM usuarios WHERE user = ?", r.FormValue("ownUser"))
		db_mu.Unlock()
		if err != nil {
			Error.Println(err)
			return
		}
		for query.Next() {
			//Obtengo el identificador del creador
			var id int
			timestamp := time.Now().Unix()
			err = query.Scan(&id)
			if err != nil {
				Error.Println(err)
				continue
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
		db_mu.Lock()
		query, err := db.Query("SELECT id FROM usuarios WHERE user = ?", r.FormValue("ownUser"))
		db_mu.Unlock()
		if err != nil {
			Error.Println(err)
			return
		}
		for query.Next() {
			//Obtengo el identificador del creador
			var id int
			timestamp := time.Now().Unix()
			err = query.Scan(&id)
			if err != nil {
				Error.Println(err)
				continue
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
	//Primera vez que se muestran los ficheros publi /msg
	if r.FormValue("accion") == "first_show" {
		db_mu.Lock()
		err0 := db.QueryRow("SELECT id FROM usuarios WHERE user = ?", r.FormValue("username")).Scan(&id_user)
		db_mu.Unlock()
		if err0 != nil {
			Error.Println(err0)
			return
		}
		if r.FormValue("tabla") == "publi" {
			db_mu.Lock()
			sql := fmt.Sprintf("SELECT id, fichero, fecha_inicio, fecha_final, destino, timestamp, gap FROM %s WHERE creador_id = %s", r.FormValue("tabla"), id_user)
			db_mu.Unlock()
			query, err := db.Query(sql)
			if err != nil {
				Error.Println(err)
				return
			}
			for query.Next() {
				var fichero, f_ini, f_fin, destinos string
				var id, timestamp, gap int64
				err = query.Scan(&id, &fichero, &f_ini, &f_fin, &destinos, &timestamp, &gap)
				if err != nil {
					Error.Println(err)
					continue
				}
				//Se obtiene la fecha de creacion de una entidad
				f_creacion := libs.FechaCreacion(timestamp)
				//Convertimos a fecha normal
				f_ini_conv := libs.FechaSQLtoNormal(f_ini)
				f_fin_conv := libs.FechaSQLtoNormal(f_fin)
				//Generamos la tabla
				output += fmt.Sprintf("<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Editar Publicidad'>%s</a>", id, fichero)
				output += fmt.Sprintf("<a href='#' onclick='borrar(%d)' title='Borrar Publicidad' style='float:right'><span class='fa fa-trash-o'></a></td>", id)
				output += fmt.Sprintf("<td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%d</td></tr>", f_ini_conv, f_fin_conv, destinos, f_creacion, gap)
			}
		}
		if r.FormValue("tabla") == "mensaje" {
			db_mu.Lock()
			sql := fmt.Sprintf("SELECT id, fichero, fecha_inicio, fecha_final, destino, timestamp, playtime FROM %s WHERE creador_id = %s", r.FormValue("tabla"), id_user)
			db_mu.Unlock()
			query, err := db.Query(sql)
			if err != nil {
				Error.Println(err)
				return
			}
			for query.Next() {
				var fichero, f_ini, f_fin, destinos, playtime string
				var id, timestamp int64
				err = query.Scan(&id, &fichero, &f_ini, &f_fin, &destinos, &timestamp, &playtime)
				if err != nil {
					Error.Println(err)
					continue
				}
				//Se obtiene la fecha de creacion de una entidad
				f_creacion := libs.FechaCreacion(timestamp)
				//Convertimos a fecha normal
				f_ini_conv := libs.FechaSQLtoNormal(f_ini)
				f_fin_conv := libs.FechaSQLtoNormal(f_fin)
				//Generamos la tabla
				output += fmt.Sprintf("<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Editar Mensaje'>%s</a>", id, fichero)
				output += fmt.Sprintf("<a href='#' onclick='borrar(%d)' title='Borrar Mensaje' style='float:right'><span class='fa fa-trash-o'></a></td>", id)
				output += fmt.Sprintf("<td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", f_ini_conv, f_fin_conv, destinos, f_creacion, playtime)
			}
		}
	}
	//Muestra los datos detallados de publicidad y mensajes recibiendo un patrón de busqueda
	if r.FormValue("accion") == "mostrar" {
		search := r.FormValue("search") //Patrón de Busqueda
		db_mu.Lock()
		err0 := db.QueryRow("SELECT id FROM usuarios WHERE user = ?", r.FormValue("username")).Scan(&id_user)
		db_mu.Unlock()
		if err0 != nil {
			Error.Println(err0)
			return
		}
		if r.FormValue("tabla") == "publi" {
			db_mu.Lock()
			sql := fmt.Sprintf("SELECT id, fichero, fecha_inicio, fecha_final, destino, timestamp, gap FROM %s WHERE creador_id = %s", r.FormValue("tabla"), id_user)
			db_mu.Unlock()
			query, err := db.Query(sql)
			if err != nil {
				Error.Println(err)
				return
			}
			for query.Next() {
				var fichero, f_ini, f_fin, destinos string
				var id, timestamp, gap int64
				err = query.Scan(&id, &fichero, &f_ini, &f_fin, &destinos, &timestamp, &gap)
				if err != nil {
					Error.Println(err)
					continue
				}
				if strings.Contains(destinos, search) {
					//Se obtiene la fecha de creacion de una entidad
					f_creacion := libs.FechaCreacion(timestamp)
					//Convertimos a fecha normal
					f_ini_conv := libs.FechaSQLtoNormal(f_ini)
					f_fin_conv := libs.FechaSQLtoNormal(f_fin)
					//Generamos la tabla
					output += fmt.Sprintf("<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Editar Publicidad'>%s</a>", id, fichero)
					output += fmt.Sprintf("<a href='#' onclick='borrar(%d)' title='Borrar Publicidad' style='float:right'><span class='fa fa-trash-o'></a></td>", id)
					output += fmt.Sprintf("<td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%d</td></tr>", f_ini_conv, f_fin_conv, destinos, f_creacion, gap)
				}
			}
		}
		if r.FormValue("tabla") == "mensaje" {
			db_mu.Lock()
			sql := fmt.Sprintf("SELECT id, fichero, fecha_inicio, fecha_final, destino, timestamp, playtime FROM %s WHERE creador_id = %s", r.FormValue("tabla"), id_user)
			db_mu.Unlock()
			query, err := db.Query(sql)
			if err != nil {
				Error.Println(err)
				return
			}
			for query.Next() {
				var fichero, f_ini, f_fin, destinos, playtime string
				var id, timestamp int64
				err = query.Scan(&id, &fichero, &f_ini, &f_fin, &destinos, &timestamp, &playtime)
				if err != nil {
					Error.Println(err)
					continue
				}
				if strings.Contains(destinos, search) {
					//Se obtiene la fecha de creacion de una entidad
					f_creacion := libs.FechaCreacion(timestamp)
					//Convertimos a fecha normal
					f_ini_conv := libs.FechaSQLtoNormal(f_ini)
					f_fin_conv := libs.FechaSQLtoNormal(f_fin)
					output += fmt.Sprintf("<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Editar Mensaje'>%s</a>", id, fichero)
					output += fmt.Sprintf("<a href='#' onclick='borrar(%d)' title='Borrar Mensaje' style='float:right'><span class='fa fa-trash-o'></a></td>", id)
					output += fmt.Sprintf("<td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", f_ini_conv, f_fin_conv, destinos, f_creacion, playtime)
				}
			}
		}
	}
	if r.FormValue("accion") == "borrar" {
		var filename string
		//Obtenemos el nombre del fichero y lo borramos de manera local
		sql := fmt.Sprintf("SELECT fichero FROM %s WHERE id = %s", r.FormValue("tabla"), r.FormValue("borrar"))
		err := db.QueryRow(sql).Scan(&filename)
		if filename != "" {
			fmt.Println(filename)
			if r.FormValue("tabla") == "mensaje" {
				//Borramos el fichero de mensaje
				err = os.Remove(msg_files_location + filename)
				if err != nil {
					Error.Println(err)
					return
				}
			}
			if r.FormValue("tabla") == "publi" {
				//Borramos el fichero de publicidad
				err = os.Remove(publi_files_location + filename)
				if err != nil {
					Error.Println(err)
					return
				}
			}
		}
		//Despues lo borramos de base de datos
		db_mu.Lock()
		sql = fmt.Sprintf("DELETE FROM %s WHERE id = %s", r.FormValue("tabla"), r.FormValue("borrar"))
		db_mu.Unlock()
		_, err = db.Exec(sql)
		if err != nil {
			Error.Println(err)
			return
		}
		output = "OK"
	}
	if r.FormValue("accion") == "load" {
		var id, gap int
		var f_inicio, f_fin, destino, horario string
		if r.FormValue("tabla") == "publi" {
			db_mu.Lock()
			sql := fmt.Sprintf("SELECT id, fecha_inicio, fecha_final, destino, gap FROM %s WHERE id = %s", r.FormValue("tabla"), r.FormValue("edit_id"))
			db_mu.Unlock()
			query, err := db.Query(sql)
			if err != nil {
				Error.Println(err)
				return
			}
			for query.Next() {
				err := query.Scan(&id, &f_inicio, &f_fin, &destino, &gap)
				if err != nil {
					Error.Println(err)
					continue
				}
				//Convertimos a fecha normal
				f_ini_conv := libs.FechaSQLtoNormal(f_inicio)
				f_fin_conv := libs.FechaSQLtoNormal(f_fin)
				panel_destino := destino + "&nbsp;&nbsp;<a href='#' onclick='edit_dom()' title='Pulsa para editar dominio'><span class='fa fa-edit'></a>"
				fmt.Fprintf(w, "id=%d&f_inicio=%s&f_fin=%s&origen=%s&gap=%d:.:%s", id, f_ini_conv, f_fin_conv, destino, gap, panel_destino)
			}
		}
		if r.FormValue("tabla") == "mensaje" {
			db_mu.Lock()
			sql := fmt.Sprintf("SELECT id, fecha_inicio, fecha_final, destino, playtime FROM %s WHERE id = %s", r.FormValue("tabla"), r.FormValue("edit_id"))
			db_mu.Unlock()
			query, err := db.Query(sql)
			if err != nil {
				Error.Println(err)
				return
			}
			for query.Next() {
				err := query.Scan(&id, &f_inicio, &f_fin, &destino, &horario)
				if err != nil {
					Error.Println(err)
					continue
				}
				//Convertimos a fecha normal
				f_ini_conv := libs.FechaSQLtoNormal(f_inicio)
				f_fin_conv := libs.FechaSQLtoNormal(f_fin)
				//Obtenemos hora y minutos por separado
				hora := libs.MostrarHoras(horario)
				minuto := libs.MostrarMinutos(horario)
				panel_destino := destino + "&nbsp;&nbsp;<a href='#' onclick='edit_dom()' title='Pulsa para editar dominio'><span class='fa fa-edit'></a>"
				fmt.Fprintf(w, "id=%d&f_inicio=%s&f_fin=%s&origen=%s:.:%s:.:%s:.:%s", id, f_ini_conv, f_fin_conv, destino, panel_destino, hora, minuto)
			}
		}
	}
	if r.FormValue("accion") == "modificar" {
		var err1 error
		timestamp := time.Now().Unix()
		//Generamos las fechas en formato SQL (Ej. 20170212)
		f_ini := libs.FechaNormaltoSQL(r.FormValue("f_ini"))
		f_fin := libs.FechaNormaltoSQL(r.FormValue("f_fin"))
		if r.FormValue("tabla") == "publi" {
			//Si el destino está vacio, no lo modificamos.
			if r.FormValue("destino") == "" {
				db_mu.Lock()
				_, err1 = db.Exec("UPDATE publi SET fecha_inicio=?, fecha_final=?, gap=?, timestamp=? WHERE id = ?", f_ini, f_fin, r.FormValue("gap"), timestamp, r.FormValue("id"))
				db_mu.Unlock()
			} else {
				db_mu.Lock()
				_, err1 = db.Exec("UPDATE publi SET fecha_inicio=?, fecha_final=?, destino=?, gap=?, timestamp=? WHERE id = ?", f_ini, f_fin, r.FormValue("destino"), r.FormValue("gap"), timestamp, r.FormValue("id"))
				db_mu.Unlock()
			}
			if err1 != nil {
				Error.Println(err1)
				return
			}
		}
		if r.FormValue("tabla") == "mensaje" {
			//Formamos el horario
			horario := r.FormValue("hora") + ":" + r.FormValue("minutos")
			//Si el destino está vacio, no lo modificamos.
			if r.FormValue("destino") == "" {
				db_mu.Lock()
				_, err1 = db.Exec("UPDATE mensaje SET fecha_inicio=?, fecha_final=?, playtime=?, timestamp=? WHERE id = ?", f_ini, f_fin, horario, timestamp, r.FormValue("id"))
				db_mu.Unlock()
			} else {
				db_mu.Lock()
				_, err1 = db.Exec("UPDATE mensaje SET fecha_inicio=?, fecha_final=?, destino=?, playtime=?, timestamp=? WHERE id = ?", f_ini, f_fin, r.FormValue("destino"), horario, timestamp, r.FormValue("id"))
				db_mu.Unlock()
			}
			if err1 != nil {
				Error.Println(err1)
				return
			}
		}
	}
	fmt.Fprint(w, output) //fmt.Println(output)
}
