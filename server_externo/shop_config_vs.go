package main

import (
	"fmt"
	"net/http"
)

func config_shop_vs(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() 					//recupera campos del form tanto GET como POST
	accion := r.FormValue("action") //variable que tiene la accion de cada org
	var list string 				//variable para imprimir la salida
	
	//Generamos el select de entidades
	if accion == "entidad" {
		ent := r.FormValue("nom_ent")
		db_mu.Lock()
		err := db.QueryRow("SELECT id FROM entidades WHERE nombre = ?", ent).Scan(&list)
		db_mu.Unlock()
		if err != nil {
			Error.Println(err)
			return
		}
	}
	//Generamos el select de almacenes y guardamos el dominio
	if accion == "almacen" {
		ent := r.FormValue("entidad")
		//Zona donde se genera el select
		if ent != "" {
			var id_alm int
			var almacen string
			db_mu.Lock()
			query, err := db.Query("SELECT id, almacen FROM almacenes WHERE entidad_id = ?", ent)
			db_mu.Unlock()
			if err != nil {
				Error.Println(err)
				return
			}
			if query.Next() {
				query.Scan(&id_alm, &almacen)
				list = fmt.Sprintf("%d<=>%s;", id_alm, almacen)
				for query.Next() {
					query.Scan(&id_alm, &almacen)
					if err != nil {
						Error.Println(err)
						continue
					}
					list += fmt.Sprintf("%d<=>%s;", id_alm, almacen)
				}
			}
		}
	}
	//Generamos el select de paises y guardamos el dominio
	if accion == "pais" {
		alm := r.FormValue("almacen")
		//Zona donde se genera el select
		if alm != "" {
			var id_pais int
			var pais string
			db_mu.Lock()
			query, err := db.Query("SELECT id, pais FROM pais WHERE almacen_id = ?", alm)
			db_mu.Unlock()
			if err != nil {
				Error.Println(err)
				return
			}
			if query.Next() {
				query.Scan(&id_pais, &pais)
				list = fmt.Sprintf("%d<=>%s;", id_pais, pais)
				for query.Next() {
					query.Scan(&id_pais, &pais)
					if err != nil {
						Error.Println(err)
						continue
					}
					list += fmt.Sprintf("%d<=>%s;", id_pais, pais)
				}
			}
		}
	}
	//Generamos el select de regiones y guardamos el dominio
	if accion == "region" {
		pais := r.FormValue("pais")
		//Zona donde se genera el select
		if pais != "" {
			var id_region int
			var region string
			db_mu.Lock()
			query, err := db.Query("SELECT id, region FROM region WHERE pais_id = ?", pais)
			db_mu.Unlock()
			if err != nil {
				Error.Println(err)
				return
			}
			if query.Next() {
				query.Scan(&id_region, &region)
				list = fmt.Sprintf("%d<=>%s;", id_region, region)
				for query.Next() {
					query.Scan(&id_region, &region)
					if err != nil {
						Error.Println(err)
						continue
					}
					list += fmt.Sprintf("%d<=>%s;", id_region, region)
				}
			}
		}
	}
	//Generamos el select de provincias y guardamos el dominio
	if accion == "provincia" {
		reg := r.FormValue("region")
		//Zona donde se genera el select
		if reg != "" {
			var id_prov int
			var prov string
			db_mu.Lock()
			query, err := db.Query("SELECT id, provincia FROM provincia WHERE region_id = ?", reg)
			db_mu.Unlock()
			if err != nil {
				Error.Println(err)
				return
			}
			if query.Next() {
				query.Scan(&id_prov, &prov)
				list = fmt.Sprintf("%d<=>%s;", id_prov, prov)
				for query.Next() {
					query.Scan(&id_prov, &prov)
					if err != nil {
						Error.Println(err)
						continue
					}
					list += fmt.Sprintf("%d<=>%s;", id_prov, prov)
				}
			}
		}
	}
	//Generamos el select de tiendas y guardamos el dominio
	if accion == "tienda" {
		prov := r.FormValue("provincia")
		//Zona donde se genera el select
		if prov != "" {
			var id_tienda int
			var tiendas string
			db_mu.Lock()
			query, err := db.Query("SELECT id, tienda FROM tiendas WHERE provincia_id = ?", prov)
			db_mu.Unlock()
			if err != nil {
				Error.Println(err)
				return
			}
			if query.Next() {
				query.Scan(&id_tienda, &tiendas)
				list = fmt.Sprintf("%d<=>%s;", id_tienda, tiendas)
				for query.Next() {
					query.Scan(&id_tienda, &tiendas)
					if err != nil {
						Error.Println(err)
						continue
					}
					list += fmt.Sprintf("%d<=>%s;", id_tienda, tiendas)
				}
			}
		}
	}
	fmt.Fprint(w, list)
}