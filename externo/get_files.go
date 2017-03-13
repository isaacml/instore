package main

import (
	"fmt"
	//"github.com/isaacml/instore/libs"
	"io"
	"net/http"
	"os"
	"strconv"
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
			_, err1 := db.Exec("INSERT INTO publi (`fichero`, `fecha_inicio`, `fecha_final`, `destino`, `creador_id`, `timestamp`) VALUES (?,?,?,?,?,?)",
				nameFileServer, r.FormValue("f_inicio"), r.FormValue("f_final"), r.FormValue("destino"), id, timestamp)
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
			_, err1 := db.Exec("INSERT INTO publi (`fichero`, `fecha_inicio`, `fecha_final`, `destino`, `creador_id`, `timestamp`, `playtime`) VALUES (?,?,?,?,?,?,?)",
				nameFileServer, r.FormValue("f_inicio"), r.FormValue("f_final"), r.FormValue("destino"), id, timestamp, r.FormValue("playtime"))
			db_mu.Unlock()
			if err1 != nil {
				Error.Println(err1)
			}
		}
	}
}
func destino(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var id_user_admin int
	var salida string
	accion := r.FormValue("action")

	query, err := db.Query("SELECT id FROM usuarios WHERE user = ?", r.FormValue("userAdmin"))
	if err != nil {
		Error.Println(err)
	}
	for query.Next() {
		//Obtengo el identificador del usuario admin externo
		err = query.Scan(&id_user_admin)
		if err != nil {
			Error.Println(err)
		}
	}
	if accion == "destinos" {
		//Seleccionamos las entidades para un user admin concreto
		query2, err := db.Query("SELECT id, nombre FROM entidades WHERE creador_id = ?", id_user_admin)
		if err != nil {
			Error.Println(err)
		}
		for query2.Next() {
			var ent string
			var id_ent int
			err = query2.Scan(&id_ent, &ent)
			if err != nil {
				Error.Println(err)
			}
			id_string := strconv.Itoa(id_ent)
			salida += id_string + ";" + ent + "::"
		}
		fmt.Fprint(w, salida)
	}
	if accion == "entidades" {
		id_entidad := r.FormValue("id_entidad")
		//Seleccionamos el almacen para una entidad concreta
		query3, err := db.Query("SELECT almacenes.id, almacenes.almacen, entidades.nombre FROM almacenes INNER JOIN entidades WHERE entidades.id = almacenes.entidad_id AND almacenes.creador_id = ? AND almacenes.entidad_id = ?", id_user_admin, id_entidad)
		if err != nil {
			Error.Println(err)
		}
		for query3.Next() {
			var almacen, entidad string
			var id_alm int
			err = query3.Scan(&id_alm, &almacen, &entidad)
			if err != nil {
				Error.Println(err)
			}
			id_string := strconv.Itoa(id_alm)
			salida += id_string + ";" + almacen + ";" + entidad + "::"
		}
		fmt.Fprint(w, salida)
	}
	if accion == "almacenes" {
		id_almacen := r.FormValue("id_almacen")
		//Seleccionamos el pais para un almacen concreto
		query4, err := db.Query("SELECT pais.id, pais.pais, almacenes.almacen FROM pais INNER JOIN almacenes WHERE almacenes.id = pais.almacen_id AND pais.creador_id = ? AND pais.almacen_id = ?", id_user_admin, id_almacen)
		if err != nil {
			Error.Println(err)
		}
		for query4.Next() {
			var pais, almacen string
			var id_country int
			err = query4.Scan(&id_country, &pais, &almacen)
			if err != nil {
				Error.Println(err)
			}
			id_string := strconv.Itoa(id_country)
			salida += id_string + ";" + pais + ";" + almacen + "::"
		}
		fmt.Fprint(w, salida)
	}
	if accion == "paises" {
		id_pais := r.FormValue("id_pais")
		//Seleccionamos la region para un país concreto
		query5, err := db.Query("SELECT region.id, region.region, pais.pais FROM region INNER JOIN pais ON region.pais_id = pais.id WHERE region.creador_id = ? AND region.pais_id = ?", id_user_admin, id_pais)
		if err != nil {
			Error.Println(err)
		}
		for query5.Next() {
			var region, pais string
			var id_reg int
			err = query5.Scan(&id_reg, &region, &pais)
			if err != nil {
				Error.Println(err)
			}
			id_string := strconv.Itoa(id_reg)
			salida += id_string + ";" + region + ";" + pais + "::"
		}
		fmt.Fprint(w, salida)
	}
	if accion == "regiones" {
		//Seleccionamos la provincia para una región concreta
		query6, err := db.Query("SELECT provincia.id, provincia.provincia, region.region FROM provincia INNER JOIN region ON provincia.region_id = region.id WHERE provincia.creador_id = ? AND provincia.region_id = ?", id_user_admin, r.FormValue("id_region"))
		if err != nil {
			Error.Println(err)
		}
		for query6.Next() {
			var provincia, region string
			var id_prov int
			err = query6.Scan(&id_prov, &provincia, &region)
			if err != nil {
				Error.Println(err)
			}
			id_string := strconv.Itoa(id_prov)
			salida += id_string + ";" + provincia + ";" + region + "::"
		}
		fmt.Fprint(w, salida)
	}
	if accion == "provincias" {
		//Seleccionamos la tienda para una provincia concreta
		query7, err := db.Query("SELECT tiendas.id, tiendas.tienda, provincia.provincia FROM tiendas INNER JOIN provincia ON tiendas.provincia_id = provincia.id WHERE tiendas.creador_id = ? AND tiendas.provincia_id = ?", id_user_admin, r.FormValue("id_provincia"))
		if err != nil {
			Error.Println(err)
		}
		for query7.Next() {
			var tienda, provincia string
			var id_tiend int
			err = query7.Scan(&id_tiend, &tienda, &provincia)
			if err != nil {
				Error.Println(err)
			}
			id_string := strconv.Itoa(id_tiend)
			salida += id_string + ";" + tienda + ";" + provincia + "::"
		}
		fmt.Fprint(w, salida)
	}
	if accion == "tiendas" {
		//Seleccionamos la tienda para una provincia concreta
		query7, err := db.Query("SELECT tienda FROM tiendas WHERE creador_id = ? AND id = ?", id_user_admin, r.FormValue("id_tienda"))
		if err != nil {
			Error.Println(err)
		}
		for query7.Next() {
			var tienda string
			err = query7.Scan(&tienda)
			if err != nil {
				Error.Println(err)
			}
			salida += tienda
		}
		fmt.Fprint(w, salida)
	}
}
