package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
	"strings"
	"strconv"
)

//Acciones independientes
func acciones(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	accion := r.FormValue("accion")
	//Toma los valores bitmap de un determinado usuario
	if accion == "bitmap_perm" {
		var output string
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
			msg_normal := libs.BitmapParsing(bitmap_hex, MSG_NORMAL) //res[4]
			//Pasamos los valores al html
			output = fmt.Sprintf("%d;%d;%d;%d;%d", prog_pub, prog_mus, prog_msg, add_mus, msg_normal)
		}
		fmt.Fprint(w, output)
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
		} else {
			output = "Ocultar"
		}
		fmt.Fprintf(w, output)
	}
	if accion == "check_entidad" {
		var st_ent int
		err := db.QueryRow("SELECT status FROM entidades WHERE nombre = ?", r.FormValue("ent")).Scan(&st_ent)
		if err != nil {
			Warning.Println(err)
		}
		fmt.Fprint(w, st_ent)
	}
	//Proporciona la informacion necesaria para generar el explorador de destino (se genera en publi.html)
	//La peticion es por parte del -explorador.go- del admin_externo
	if accion == "destinos" {
		var id_user_admin, id_father_admin int
		var tipo_de_usuario string
		internal_action := r.FormValue("internal_action")
		//Obtengo el identificador y el padre del usuario
		err := db.QueryRow("SELECT id, padre_id FROM usuarios WHERE user = ?", r.FormValue("userAdmin")).Scan(&id_user_admin, &id_father_admin)
		if err != nil {
			Error.Println(err)
		}
		//Si su padre no es un super-admin, solo puede añadir publicidad en su propio dominio
		if id_father_admin > 1 {
			var ent, alm, pais, reg, prov, shop string
			tipo_de_usuario = "NOADMIN"
			SQL_DOMIN := "SELECT entidades.nombre, almacenes.almacen, pais.pais, region.region, provincia.provincia, tiendas.tienda"
			SQL_DOMIN += " FROM usuarios INNER JOIN entidades ON usuarios.entidad_id = entidades.id"
			SQL_DOMIN += " INNER JOIN almacenes ON almacenes.entidad_id = entidades.id INNER JOIN pais ON pais.almacen_id = almacenes.id"
			SQL_DOMIN += " INNER JOIN region ON region.pais_id = pais.id INNER JOIN provincia ON provincia.region_id = region.id"
			SQL_DOMIN += " INNER JOIN tiendas ON tiendas.provincia_id = provincia.id WHERE usuarios.id = ?"
			err := db.QueryRow(SQL_DOMIN, id_user_admin).Scan(&ent, &alm, &pais, &reg, &prov, &shop)
			if err != nil {
				Error.Println(err)
			}
			fmt.Fprintf(w, "%s@@%s.%s.%s.%s.%s.%s", tipo_de_usuario, ent, alm, pais, reg, prov, shop)
		} else {
			var salida string
			if internal_action == "gent_ent" {
				var ents string
				tipo_de_usuario = "ADMIN"
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
					ents += id_string + ";" + ent + "::"
				}
				//Primera vez que se muestran datos: pasamos el tipo de usuario
				salida = tipo_de_usuario + "@@" + ents
			}
			if internal_action == "entidades" {
				var alms string
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
					alms += id_string + ";" + almacen + ";" + entidad + "::"
				}
				salida = alms
			}
			if internal_action == "almacenes" {
				var paises string
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
					paises += id_string + ";" + pais + ";" + almacen + "::"
				}
				salida = paises
			}
			if internal_action == "paises" {
				var regiones string
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
					regiones += id_string + ";" + region + ";" + pais + "::"
				}
				salida = regiones
			}
			if internal_action == "regiones" {
				var provs string
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
					provs += id_string + ";" + provincia + ";" + region + "::"
				}
				salida = provs
			}
			if internal_action == "provincias" {
				var shops string
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
					shops += id_string + ";" + tienda + ";" + provincia + "::"
				}
				salida = shops
			}
			if internal_action == "tiendas" {
				var final string
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
					final += tienda
				}
				salida = final
			}
			//fmt.Println(salida)
			fmt.Fprint(w, salida)
		}
	}
}

//Funcion que toma la peticion de dominio por parte de la tienda.
func recoger_dominio(w http.ResponseWriter, r *http.Request) {
	var output string
	var domains []string
	doms := strings.Split(r.FormValue("dominios"), ":.:")
	doms = doms[:len(doms)-1]
	for _, val := range doms {
		//enviamos el dominio a DomainGenerator() en la libreria de funciones
		dom := libs.DomainGenerator(val)
		for _, v := range dom {
			domains = append(domains, v)
		}
	}
	//Borramos dominios duplicados
	domains = libs.RemoveDuplicates(domains)
	//Nuestra fecha actual personalizada
	fecha := libs.MyCurrentDate()
	//Formamos una cadena con la publicidad y los mensajes para ese dominio y con la fecha que le corresponde
	output += "[publi]"
	for _, val := range domains {
		publicidad, err := db.Query("SELECT fichero, fecha_inicio, fecha_final, gap FROM publi WHERE destino = ?", val)
		if err != nil {
			Error.Println(err)
		}
		for publicidad.Next() {
			var f_publi, fecha_ini, fecha_fin, gap string
			err = publicidad.Scan(&f_publi, &fecha_ini, &fecha_fin, &gap)
			if err != nil {
				Error.Println(err)
			}
			//BETWEEN
			if fecha_ini <= fecha && fecha_fin >= fecha {
				output += ";" + f_publi + "<=>" + fecha_ini + "<=>" + fecha_fin + "<=>" + gap
			}
		}
	}
	output += "[mensaje]"
	for _, val := range domains {
		mensajes, err := db.Query("SELECT fichero, fecha_inicio, fecha_final, playtime FROM mensaje WHERE destino = ?", val)
		if err != nil {
			Error.Println(err)
		}
		for mensajes.Next() {
			var f_msg, fecha_ini, fecha_fin, playtime string
			err = mensajes.Scan(&f_msg, &fecha_ini, &fecha_fin, &playtime)
			if err != nil {
				Error.Println(err)
			}
			//BETWEEN
			if fecha_ini <= fecha && fecha_fin >= fecha {
				output += ";" + f_msg + "<=>" + fecha_ini + "<=>" + fecha_fin + "<=>" + playtime
			}
		}
	}
	//Enviamos la cadena
	fmt.Fprint(w, output) //fmt.Println(output)
}
