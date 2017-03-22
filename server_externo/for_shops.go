package main

import (
	"fmt"
	"net/http"
	"strings"
)

//Variable que va a guardar el dominio de la tienda
var status_dom string

func config_shop(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	var domain string
	accion := r.FormValue("action")

	//Generamos el select de entidades
	if accion == "entidad" {
		user := r.FormValue("username")
		query, err := db.Query("SELECT id FROM usuarios WHERE user = ?", user)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			var entidad, id_ent int
			var list, name string
			err = query.Scan(&entidad)
			if err != nil {
				Error.Println(err)
			}
			query, err := db.Query("SELECT id, nombre FROM entidades WHERE creador_id=?", entidad)
			if err != nil {
				Error.Println(err)
			}
			list = "<div class='panel-heading'>Entidad</div><div class='panel-body'><select name='entidad'>"
			if query.Next() {
				list += "<option value='' selected>Selecciona una entidad</option>"
				query.Scan(&id_ent, &name)
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_ent, name)
				for query.Next() {
					query.Scan(&id_ent, &name)
					if err != nil {
						Error.Println(err)
					}
					list += fmt.Sprintf("<option value='%d'>%s</option>", id_ent, name)
				}
			} else {
				list += "<option value='' selected>No hay entidades</option>"
			}
			list += "</select></div>"
			fmt.Fprint(w, list)
		}
	}
	//Generamos el select de almacenes y guardamos el dominio
	if accion == "almacen" {
		var list, entidad string
		ent := r.FormValue("entidad")

		//Zona donde se genera el select
		if ent != "" {
			var id_alm int
			var almacen string
			query, err := db.Query("SELECT id, almacen FROM almacenes WHERE entidad_id = ?", ent)
			if err != nil {
				Error.Println(err)
			}
			list = "<div class='panel-heading'>Almacen</div><div class='panel-body'><select name='almacen'>"
			if query.Next() {
				list += "<option value='' selected>Selecciona un almacen</option>"
				query.Scan(&id_alm, &almacen)
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_alm, almacen)
				for query.Next() {
					query.Scan(&id_alm, &almacen)
					if err != nil {
						Error.Println(err)
					}
					list += fmt.Sprintf("<option value='%d'>%s</option>", id_alm, almacen)
				}
			} else {
				list += "<option value='' selected>No hay almacenes</option>"
			}
			list += "</select></div>"

			//Zona de Guardado de Dominio
			errdom := db.QueryRow("SELECT nombre FROM entidades WHERE id = ?", ent).Scan(&entidad)
			if errdom != nil {
				Error.Println(errdom)
			}
			status_dom = entidad
			domain = "<span style='color: #B8860B'>" + status_dom + "</span><input type='hidden' name='dominio' value'" + status_dom + "'>"
		} else {
			list = "<div class='panel-heading'>Almacen</div><div class='panel-body'><select name='almacen'><option value='' selected>Requiere una entidad</option></select></div>"
			status_dom = ""
			domain = "<span style='color: #B8860B'>" + status_dom + "</span><input type='hidden' name='dominio' value'" + status_dom + "'>"
		}
		fmt.Fprintf(w, "%s;%s", list, domain)
	}
	//Generamos el select de paises y guardamos el dominio
	if accion == "pais" {
		var list, almacen string
		alm := r.FormValue("almacen")
		partir := strings.Split(status_dom, ".") // partimos el estado del dominio, para poder modificarlo

		//Zona donde se genera el select
		if alm != "" {
			var id_pais int
			var pais string
			query, err := db.Query("SELECT id, pais FROM pais WHERE almacen_id = ?", alm)
			if err != nil {
				Error.Println(err)
			}
			list = "<div class='panel-heading'>País</div><div class='panel-body'><select name='pais'>"
			if query.Next() {
				list += "<option value='' selected>Selecciona un país</option>"
				query.Scan(&id_pais, &pais)
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_pais, pais)
				for query.Next() {
					query.Scan(&id_pais, &pais)
					if err != nil {
						Error.Println(err)
					}
					list += fmt.Sprintf("<option value='%d'>%s</option>", id_pais, pais)
				}
			} else {
				list += "<option value='' selected>No hay paises</option>"
			}
			list += "</select></div>"

			//Zona de Guardado de Dominio
			errdom := db.QueryRow("SELECT almacen FROM almacenes WHERE id = ?", alm).Scan(&almacen)
			if errdom != nil {
				Error.Println(errdom)
			}
			if strings.Contains(status_dom, ".") {
				status_dom = partir[0] + "." + almacen
				domain = "<span style='color: #B8860B'>" + status_dom + "</span><input type='hidden' name='dominio' value'" + status_dom + "'>"
			} else {
				status_dom = status_dom + "." + almacen
				domain = "<span style='color: #B8860B'>" + status_dom + "</span><input type='hidden' name='dominio' value'" + status_dom + "'>"
			}
		} else {
			status_dom = partir[0]
			list = "<div class='panel-heading'>País</div><div class='panel-body'><select id='pais' name='pais'><option value='' selected>Requiere un almacen</option></select></div>"
			domain = "<span style='color: #B8860B'>" + status_dom + "</span><input type='hidden' name='dominio' value'" + status_dom + "'>"
		}
		fmt.Fprintf(w, "%s;%s", list, domain)
	}
	//Generamos el select de regiones y guardamos el dominio
	if accion == "region" {
		var list, country string
		pais := r.FormValue("pais")
		partir := strings.Split(status_dom, ".") // partimos el estado del dominio, para poder modificarlo

		//Zona donde se genera el select
		if pais != "" {
			var id_region int
			var region string
			query, err := db.Query("SELECT id, region FROM region WHERE pais_id = ?", pais)
			if err != nil {
				Error.Println(err)
			}
			list = "<div class='panel-heading'>Región</div><div class='panel-body'><select name='region'>"
			if query.Next() {
				list += "<option value='' selected>Selecciona una región</option>"
				query.Scan(&id_region, &region)
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_region, region)
				for query.Next() {
					query.Scan(&id_region, &region)
					if err != nil {
						Error.Println(err)
					}
					list += fmt.Sprintf("<option value='%d'>%s</option>", id_region, region)
				}
			} else {
				list += "<option value='' selected>No hay regiones</option>"
			}
			list += "</select></div>"

			//Zona de Guardado de Dominio
			errdom := db.QueryRow("SELECT pais FROM pais WHERE id = ?", pais).Scan(&country)
			if errdom != nil {
				Error.Println(errdom)
			}
			status_dom = partir[0] + "." + partir[1] + "." + country
			domain = "<span style='color: #B8860B'>" + status_dom + "</span><input type='hidden' name='dominio' value'" + status_dom + "'>"
		} else {
			list = "<div class='panel-heading'>Región</div><div class='panel-body'><select name='region'><option value='' selected>Requiere un país</option></select></div>"
			status_dom = partir[0] + "." + partir[1]
			domain = "<span style='color: #B8860B'>" + status_dom + "</span><input type='hidden' name='dominio' value'" + status_dom + "'>"
		}
		fmt.Fprintf(w, "%s;%s", list, domain)
	}
	//Generamos el select de provincias y guardamos el dominio
	if accion == "provincia" {
		var list, region string
		reg := r.FormValue("region")
		partir := strings.Split(status_dom, ".") // partimos el estado del dominio, para poder modificarlo

		//Zona donde se genera el select
		if reg != "" {
			var id_prov int
			var prov string
			query, err := db.Query("SELECT id, provincia FROM provincia WHERE region_id = ?", reg)
			if err != nil {
				Error.Println(err)
			}
			list = "<div class='panel-heading'>Provincia</div><div class='panel-body'><select name='provincia'>"
			if query.Next() {
				list += "<option value='' selected>Selecciona una provincia</option>"
				query.Scan(&id_prov, &prov)
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_prov, prov)
				for query.Next() {
					query.Scan(&id_prov, &prov)
					if err != nil {
						Error.Println(err)
					}
					list += fmt.Sprintf("<option value='%d'>%s</option>", id_prov, prov)
				}
			} else {
				list += "<option value='' selected>No hay provincias</option>"
			}
			list += "</select></div>"

			//Zona de Guardado de Dominio
			errdom := db.QueryRow("SELECT region FROM region WHERE id = ?", reg).Scan(&region)
			if errdom != nil {
				Error.Println(errdom)
			}
			status_dom = partir[0] + "." + partir[1] + "." + partir[2] + "." + region
			domain = "<span style='color: #B8860B'>" + status_dom + "</span><input type='hidden' name='dominio' value'" + status_dom + "'>"
		} else {
			list = "<div class='panel-heading'>Provincia</div><div class='panel-body'><select name='provincia'><option value='' selected>Requiere una región</option></select></div>"
			status_dom = partir[0] + "." + partir[1] + "." + partir[2]
			domain = "<span style='color: #B8860B'>" + status_dom + "</span><input type='hidden' name='dominio' value'" + status_dom + "'>"
		}
		fmt.Fprintf(w, "%s;%s", list, domain)
	}
	//Generamos el select de tiendas y guardamos el dominio
	if accion == "tienda" {
		var list, provincia string
		prov := r.FormValue("provincia")
		partir := strings.Split(status_dom, ".") // partimos el estado del dominio, para poder modificarlo

		//Zona donde se genera el select
		if prov != "" {
			var id_tienda int
			var tiendas string
			query, err := db.Query("SELECT id, tienda FROM tiendas WHERE provincia_id = ?", prov)
			if err != nil {
				Error.Println(err)
			}
			list = "<div class='panel-heading'>Tienda</div><div class='panel-body'><select name='tienda'>"
			if query.Next() {
				list += "<option value='' selected>Selecciona una tienda</option>"
				query.Scan(&id_tienda, &tiendas)
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_tienda, tiendas)
				for query.Next() {
					query.Scan(&id_tienda, &tiendas)
					if err != nil {
						Error.Println(err)
					}
					list += fmt.Sprintf("<option value='%d'>%s</option>", id_tienda, tiendas)
				}
			} else {
				list += "<option value='' selected>No hay tiendas</option>"
			}
			list += "</select></div>"

			//Zona de Guardado de Dominio
			errdom := db.QueryRow("SELECT provincia FROM provincia WHERE id = ?", prov).Scan(&provincia)
			if errdom != nil {
				Error.Println(errdom)
			}
			status_dom = partir[0] + "." + partir[1] + "." + partir[2] + "." + partir[3] + "." + provincia
			domain = "<span style='color: #B8860B'>" + status_dom + "</span><input type='hidden' name='dominio' value'" + status_dom + "'>"
		} else {
			list = "<div class='panel-heading'>Tienda</div><div class='panel-body'><select name='tienda'><option value='' selected>Requiere una provincia</option></select></div>"
			status_dom = partir[0] + "." + partir[1] + "." + partir[2] + "." + partir[3]
			domain = "<span style='color: #B8860B'>" + status_dom + "</span><input type='hidden' name='dominio' value'" + status_dom + "'>"
		}
		fmt.Fprintf(w, "%s;%s", list, domain)
	}
	//guardamos el dominio de la tienda
	if accion == "cod_tienda" {
		var tienda string
		shop := r.FormValue("tienda")
		partir := strings.Split(status_dom, ".") // partimos el estado del dominio, para poder modificarlo
		if shop != "" {
			//Zona de Guardado de Dominio
			err := db.QueryRow("SELECT tienda FROM tiendas WHERE id = ?", shop).Scan(&tienda)
			if err != nil {
				Error.Println(err)
			}
			status_dom = partir[0] + "." + partir[1] + "." + partir[2] + "." + partir[3] + "." + partir[4] + "." + tienda
			domain = "<span style='color: #B8860B'>" + status_dom + "</span><input type='hidden' name='dominio' value'" + status_dom + "'>"
		} else {
			status_dom = partir[0] + "." + partir[1] + "." + partir[2] + "." + partir[3] + "." + partir[4]
			domain = "<span style='color: #B8860B'>" + status_dom + "</span><input type='hidden' name='dominio' value'" + status_dom + "'>"
		}
		fmt.Fprintf(w, ";%s", domain)
	}
}
