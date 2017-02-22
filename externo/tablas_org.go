package main

import (
	"net/http"
	"fmt"
	"time"
)

//Función que va a mostrar los usuarios segun su padre en una tabla
func get_users(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	username := r.FormValue("username")
	
	query, err := db.Query("SELECT id FROM usuarios WHERE user = ?", username)
	if err != nil {
		Error.Println(err)
	}
	for query.Next() {
		var id, dad_id int
		var user, all_name, pass string
		err = query.Scan(&dad_id)
		if err != nil {
			Error.Println(err)
		}
		query, err := db.Query("SELECT id, user, nombre_completo, pass FROM usuarios WHERE padre_id = ?", dad_id)
		if err != nil {
			Warning.Println(err)
		}
		for query.Next() {
			err = query.Scan(&id, &user, &all_name, &pass)
			if err != nil {
				Error.Println(err)
			}
			fmt.Fprintf(w, "<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Pulsa para editar el usuario'>%s</a></td><td>%s</td><td>%s</td></tr>", 
						id, user, all_name, pass)
		}
	}
}
//Función que va a mostrar las entidades en una tabla segun su usuario creador
func get_entidad(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	username := r.FormValue("username")
	
	query, err := db.Query("SELECT id FROM usuarios WHERE user = ?", username)
	if err != nil {
		Error.Println(err)
	}
	for query.Next() {
		var id, creador_id int
		var tiempo int64
		var nombre string
		err = query.Scan(&creador_id)
		if err != nil {
			Error.Println(err)
		}
		query, err := db.Query("SELECT id, nombre, timestamp FROM entidades WHERE creador_id = ?", creador_id)
		if err != nil {
			Warning.Println(err)
		}
		for query.Next() {
			err = query.Scan(&id, &nombre, &tiempo)
			if err != nil {
				Error.Println(err)
			}
			creacion := time.Unix(tiempo, 0)
			fmt.Fprintf(w, "<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Pulsa para editar entidad'>%s</a></td><td>%s</td></tr>", 
						id, nombre, creacion)
		}
	}
}
//Función que va a mostrar los usuarios segun su padre en una tabla
func get_almacen(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	username := r.FormValue("username")
	
	query, err := db.Query("SELECT id FROM usuarios WHERE user = ?", username)
	if err != nil {
		Error.Println(err)
	}
	for query.Next() {
		var id, creador_id int
		var tiempo int64
		var almacen, entidad string
		err = query.Scan(&creador_id)
		if err != nil {
			Error.Println(err)
		}
		query, err := db.Query("SELECT almacenes.id, almacenes.almacen, almacenes.timestamp, entidades.nombre FROM entidades INNER JOIN almacenes ON almacenes.entidad_id = entidades.id WHERE almacenes.creador_id = ?", creador_id)
		if err != nil {
			Warning.Println(err)
		}
		for query.Next() {
			err = query.Scan(&id, &almacen, &tiempo, &entidad)
			if err != nil {
				Error.Println(err)
			}
			creacion := time.Unix(tiempo, 0)
			fmt.Fprintf(w, "<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Pulsa para editar almacen'>%s</a></td><td>%s</td><td>%s</td></tr>", 
						id, almacen, creacion, entidad)
		}
	}
}
//Función que va a mostrar los paises segun su padre en una tabla
func get_pais(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	username := r.FormValue("username")
	
	query, err := db.Query("SELECT id FROM usuarios WHERE user = ?", username)
	if err != nil {
		Error.Println(err)
	}
	for query.Next() {
		var id, creador_id int
		var tiempo int64
		var pais, almacen string
		err = query.Scan(&creador_id)
		if err != nil {
			Error.Println(err)
		}
		query, err := db.Query("SELECT pais.id, pais.pais, pais.timestamp, almacenes.almacen FROM pais INNER JOIN almacenes ON pais.almacen_id = almacenes.id WHERE pais.creador_id = ?", creador_id)
		if err != nil {
			Warning.Println(err)
		}
		for query.Next() {
			err = query.Scan(&id, &pais, &tiempo, &almacen)
			if err != nil {
				Error.Println(err)
			}
			creacion := time.Unix(tiempo, 0)
			fmt.Fprintf(w, "<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Pulsa para editar país'>%s</a></td><td>%s</td><td>%s</td></tr>", 
						id, pais, creacion, almacen)
		}
	}
}
//Función que va a mostrar las regiones segun su padre en una tabla
func get_region(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	username := r.FormValue("username")
	
	query, err := db.Query("SELECT id FROM usuarios WHERE user = ?", username)
	if err != nil {
		Error.Println(err)
	}
	for query.Next() {
		var id, creador_id int
		var tiempo int64
		var pais, region string
		err = query.Scan(&creador_id)
		if err != nil {
			Error.Println(err)
		}
		query, err := db.Query("SELECT region.id, region.region, region.timestamp, pais.pais FROM region INNER JOIN pais ON region.pais_id = pais.id WHERE region.creador_id = ?", creador_id)
		if err != nil {
			Warning.Println(err)
		}
		for query.Next() {
			err = query.Scan(&id, &region, &tiempo, &pais)
			if err != nil {
				Error.Println(err)
			}
			creacion := time.Unix(tiempo, 0)
			fmt.Fprintf(w, "<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Pulsa para editar región'>%s</a></td><td>%s</td><td>%s</td></tr>", 
						id, region, creacion, pais)
		}
	}
}
//Función que va a mostrar las provincias segun su padre en una tabla
func get_provincia(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	username := r.FormValue("username")
	
	query, err := db.Query("SELECT id FROM usuarios WHERE user = ?", username)
	if err != nil {
		Error.Println(err)
	}
	for query.Next() {
		var id, creador_id int
		var tiempo int64
		var provincia, region string
		err = query.Scan(&creador_id)
		if err != nil {
			Error.Println(err)
		}
		query, err := db.Query("SELECT provincia.id, provincia.provincia, provincia.timestamp, region.region FROM provincia INNER JOIN region ON provincia.region_id = region.id WHERE provincia.creador_id = ?", creador_id)
		if err != nil {
			Warning.Println(err)
		}
		for query.Next() {
			err = query.Scan(&id, &provincia, &tiempo, &region)
			if err != nil {
				Error.Println(err)
			}
			creacion := time.Unix(tiempo, 0)
			fmt.Fprintf(w, "<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Pulsa para editar provincia'>%s</a></td><td>%s</td><td>%s</td></tr>", 
						id, provincia, creacion, region)
		}
	}
}
//Función que va a mostrar las tiendas segun su padre en una tabla
func get_tienda(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	username := r.FormValue("username")
	
	query, err := db.Query("SELECT id FROM usuarios WHERE user = ?", username)
	if err != nil {
		Error.Println(err)
	}
	for query.Next() {
		var id, creador_id int
		var tiempo int64
		var provincia, tienda, address, phone, extra string
		err = query.Scan(&creador_id)
		if err != nil {
			Error.Println(err)
		}
		query, err := db.Query("SELECT tiendas.id, tiendas.tienda, tiendas.timestamp, provincia.provincia, tiendas.address, tiendas.phone, tiendas.extra FROM tiendas INNER JOIN provincia ON tiendas.provincia_id = provincia.id WHERE tiendas.creador_id = ?", creador_id)
		if err != nil {
			Warning.Println(err)
		}
		for query.Next() {
			err = query.Scan(&id, &tienda, &tiempo, &provincia, &address, &phone, &extra)
			if err != nil {
				Error.Println(err)
			}
			creacion := time.Unix(tiempo, 0)
			fmt.Fprintf(w, "<tr class='odd gradeX'><td><a href='#' onclick='load(%d)' title='Pulsa para editar tienda'>%s</a></td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", 
						id, tienda, creacion, provincia, address, phone, extra)
		}
	}
}