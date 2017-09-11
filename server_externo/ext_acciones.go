package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
	"strconv"
)

//Acciones independientes
func acciones(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	accion := r.FormValue("accion")
	//Toma los valores bitmap de un determinado usuario
	if accion == "bitmap_perm" {
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
	if accion == "show_org" {
		var padre int
		var output string
		err := db.QueryRow("SELECT padre_id FROM usuarios WHERE user = ?", r.FormValue("user")).Scan(&padre)
		if err != nil {
			Error.Println(err)
		}
		if padre == 0 || padre == 1 {
			output = "Mostrar"
		}else{
			output = "Ocultar"
		}
		fmt.Fprintf(w, output)
	}
	//Proporciona la informacion necesaria para generar el explorador de destino (se genera en publi.html)
	//La peticion es por parte del -explorador.go- del admin_externo
	if accion == "destinos" {
		var id_user_admin int
		var salida string
		internal_action := r.FormValue("internal_action")
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
		if internal_action == "gent_ent" {
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
		}
		if internal_action == "entidades" {
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
		}
		if internal_action == "almacenes" {
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
		}
		if internal_action == "paises" {
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
		}
		if internal_action == "regiones" {
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
		}
		if internal_action == "provincias" {
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
		}
		if internal_action == "tiendas" {
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
		}
		//fmt.Println(salida)
		fmt.Fprint(w, salida)
	}
}