package main

import (
	"net/http"
	"fmt"
)

//Función que va a cargar en el formulario a un usuario concreto
func load_user(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	edit_id := r.FormValue("edit_id")
	
	var id, ent_id, dad_id int
	var user, all_name, pass string
	query, err := db.Query("SELECT id, user, nombre_completo, pass, entidad_id, padre_id FROM usuarios WHERE id = ?", edit_id)
	if err != nil {
		Error.Println(err)
	}
	for query.Next() {
		err = query.Scan(&id, &user, &all_name, &pass, &ent_id, &dad_id)
		if err != nil {
			Error.Println(err)
		}
		fmt.Fprintf(w, "id=%d&user=%s&name_user=%s&pass=%s&entidad=%d&padre=%d", id, user, all_name, pass, ent_id, dad_id)
	}
}
//Función que va a cargar en el formulario una entidad concreta
func load_entidad(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	edit_id := r.FormValue("edit_id")
	
	var id int
	var nombre string
	query, err := db.Query("SELECT id, nombre FROM entidades WHERE id = ?", edit_id)
	if err != nil {
		Error.Println(err)
	}
	for query.Next() {
		err = query.Scan(&id, &nombre)
		if err != nil {
			Error.Println(err)
		}
		fmt.Fprintf(w, "id=%d&entidad=%s", id, nombre)
	}
}
//Función que va a cargar en el formulario un almacen concreto
func load_almacen(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	edit_id := r.FormValue("edit_id")
	
	var id, ent_id int
	var almacen string
	query, err := db.Query("SELECT id, almacen, entidad_id FROM almacenes WHERE id = ?", edit_id)
	if err != nil {
		Error.Println(err)
	}
	for query.Next() {
		err = query.Scan(&id, &almacen, &ent_id)
		if err != nil {
			Error.Println(err)
		}
		fmt.Fprintf(w, "id=%d&almacen=%s&entidad=%d", id, almacen, ent_id)
	}
}
//Función que va a cargar en el formulario un país concreto
func load_pais(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	edit_id := r.FormValue("edit_id")
	
	var id, almacen_id int
	var pais string
	query, err := db.Query("SELECT id, pais, almacen_id FROM pais WHERE id = ?", edit_id)
	if err != nil {
		Error.Println(err)
	}
	for query.Next() {
		err = query.Scan(&id, &pais, &almacen_id)
		if err != nil {
			Error.Println(err)
		}
		fmt.Fprintf(w, "id=%d&pais=%s&almacen=%d", id, pais, almacen_id)
	}
}
//Función que va a cargar en el formulario una región concreta
func load_region(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	edit_id := r.FormValue("edit_id")
	
	var id, pais_id int
	var region string
	query, err := db.Query("SELECT id, region, pais_id FROM region WHERE id = ?", edit_id)
	if err != nil {
		Error.Println(err)
	}
	for query.Next() {
		err = query.Scan(&id, &region, &pais_id)
		if err != nil {
			Error.Println(err)
		}
		fmt.Fprintf(w, "id=%d&region=%s&pais=%d", id, region, pais_id)
	}
}
//Función que va a cargar en el formulario una provincia concreta
func load_provincia(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	edit_id := r.FormValue("edit_id")
	
	var id, region_id int
	var provincia string
	query, err := db.Query("SELECT id, provincia, region_id FROM provincia WHERE id = ?", edit_id)
	if err != nil {
		Error.Println(err)
	}
	for query.Next() {
		err = query.Scan(&id, &provincia, &region_id)
		if err != nil {
			Error.Println(err)
		}
		fmt.Fprintf(w, "id=%d&provincia=%s&region=%d", id, provincia, region_id)
	}
}
//Función que va a cargar en el formulario una tienda concreta
func load_tienda(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	edit_id := r.FormValue("edit_id")
	
	var id, prov_id int
	var tienda, address, phone, extra string
	query, err := db.Query("SELECT id, tienda, provincia_id, address, phone, extra FROM tiendas WHERE id = ?", edit_id)
	if err != nil {
		Error.Println(err)
	}
	for query.Next() {
		err = query.Scan(&id, &tienda, &prov_id, &address, &phone, &extra)
		if err != nil {
			Error.Println(err)
		}
		fmt.Fprintf(w, "id=%d&tienda=%s&provincia=%d&address=%s&phone=%s&extra=%s", id, tienda, prov_id, address, phone, extra)
	}
}