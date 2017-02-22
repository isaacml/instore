package main

import (
	"fmt"
	"net/http"
)

//Función que va a establecer un permiso o padre (ROOT o normal)
func user_permiso(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	user := r.FormValue("username")
	query, err := db.Query("SELECT id, padre_id FROM usuarios WHERE user = ?", user)
	if err != nil {
		Error.Println(err)
	}
	for query.Next() {
		var id int
		var permiso int
		err = query.Scan(&id, &permiso)
		if err != nil {
			Error.Println(err)
		}
		//permiso = 0 : es un usuario ROOT, le permitimos estas opciones
		if permiso == 0{
			fmt.Fprintf(w, "<div class='panel-heading'>¿Quien es el padre?</div><div class='panel-body'><select id='padre' name='padre'><option value='0'>ROOT</option><option value='%d' selected>%s</option></select></div>", id, user)
		}
	}
}
//Función que va a establecer una entidad (ROOT o normal)
func user_entidad(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	user := r.FormValue("username")
	query, err := db.Query("SELECT id, entidad_id FROM usuarios WHERE user = ?", user)
	if err != nil {
		Error.Println(err)
	}
	for query.Next() {
		var id int
		var entidad int
		err = query.Scan(&id, &entidad)
		if err != nil {
			Error.Println(err)
		}
		//permiso = 0 : es un usuario ROOT, le permitimos estas opciones
		if entidad == 0 {
			var list string
			query, err := db.Query("SELECT id, nombre FROM entidades")
			if err != nil {
				Error.Println(err)
			}
			list = "<div class='panel-heading'>Entidad</div><div class='panel-body'><select id='entidad' name='entidad'><option value='0' selected>ROOT</option>"
			for query.Next() {
				var id_ent int
				var name string
				err = query.Scan(&id_ent, &name)
				if err != nil {
					Error.Println(err)
				}
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_ent, name)
			}
			list += "</select></div>"
			fmt.Fprint(w, list)
		}
	}
}
//Función que va a mostrar un select de entidades según el usuario (ROOT o normal)
func almacen_entidad(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	user := r.FormValue("username")
	query, err := db.Query("SELECT id, entidad_id, padre_id FROM usuarios WHERE user = ?", user)
	if err != nil {
		Error.Println(err)
	}
	for query.Next() {
		var id, entidad, padre int
		err = query.Scan(&id, &entidad, &padre)
		if err != nil {
			Error.Println(err)
		}
		//entidad = 0 & padre = 0: es un usuario ROOT, puede ver cualquier entidad
		if entidad == 0{
			var list string
			query, err := db.Query("SELECT id, nombre FROM entidades")
			if err != nil {
				Error.Println(err)
			}
			list = "<div class='panel-heading'>Entidad</div><div class='panel-body'><select id='entidad' name='entidad'>"
			if query.Next(){
				var id_ent int
				var name string
				err = query.Scan(&id_ent, &name)
				if err != nil {
					Error.Println(err)
				}
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_ent, name)
				for query.Next() {
					err = query.Scan(&id_ent, &name)
					if err != nil {
						Error.Println(err)
					}
					list += fmt.Sprintf("<option value='%d'>%s</option>", id_ent, name)
				}
				list += "</select></div>"
				fmt.Fprint(w, list)
			}else{
				list += "<option value=''>No hay entidades</option></select></div>"
				fmt.Fprint(w, list)
			}
		}else{
			//usuario NORMAL: solo puede ver las entidades que él creó
			var list string
			query, err := db.Query("SELECT id, nombre FROM entidades WHERE creador_id = ?", id)
			if err != nil {
				Error.Println(err)
			}
			list = "<div class='panel-heading'>Entidad</div><div class='panel-body'><select id='entidad' name='entidad'>"
			if query.Next(){
				var id_ent int
				var name string
				err = query.Scan(&id_ent, &name)
				if err != nil {
					Error.Println(err)
				}
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_ent, name)
				for query.Next() {
					err = query.Scan(&id_ent, &name)
					if err != nil {
						Error.Println(err)
					}
					list += fmt.Sprintf("<option value='%d'>%s</option>", id_ent, name)
				}
				list += "</select></div>"
				fmt.Fprint(w, list)
			}else{
				list += "<option value=''>No hay entidades</option></select></div>"
				fmt.Fprint(w, list)
			}
		}
	}
}

//Función que va a mostrar un select de almacenes según el usuario (ROOT o normal)
func pais_almacen(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	user := r.FormValue("username")
	query, err := db.Query("SELECT id, entidad_id, padre_id FROM usuarios WHERE user = ?", user)
	if err != nil {
		Error.Println(err)
	}
	for query.Next() {
		var id, entidad, padre int
		err = query.Scan(&id, &entidad, &padre)
		if err != nil {
			Error.Println(err)
		}
		//entidad = 0 & padre = 0: es un usuario ROOT, puede ver cualquier almacen
		if entidad == 0{
			var list string
			query, err := db.Query("SELECT id, almacen FROM almacenes")
			if err != nil {
				Error.Println(err)
			}
			list = "<div class='panel-heading'>Almacen</div><div class='panel-body'><select id='almacen' name='almacen'>"
			if query.Next(){
				var id_alm int
				var name string
				err = query.Scan(&id_alm, &name)
				if err != nil {
					Error.Println(err)
				}
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_alm, name)
				for query.Next() {
					err = query.Scan(&id_alm, &name)
					if err != nil {
						Error.Println(err)
					}
					list += fmt.Sprintf("<option value='%d'>%s</option>", id_alm, name)
				}
				list += "</select></div>"
				fmt.Fprint(w, list)
			}else{
				list += "<option value=''>No hay almacenes</option></select></div>"
				fmt.Fprint(w, list)
			}
		}else{
			//usuario NORMAL: solo puede ver los paises que él creó
			var list string
			query, err := db.Query("SELECT id, almacen FROM almacenes WHERE creador_id = ?", id)
			if err != nil {
				Error.Println(err)
			}
			list = "<div class='panel-heading'>Almacen</div><div class='panel-body'><select id='almacen' name='almacen'>"
			if query.Next(){
				var id_alm int
				var name string
				err = query.Scan(&id_alm, &name)
				if err != nil {
					Error.Println(err)
				}
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_alm, name)
				for query.Next() {
					err = query.Scan(&id_alm, &name)
					if err != nil {
						Error.Println(err)
					}
					list += fmt.Sprintf("<option value='%d'>%s</option>", id_alm, name)
				}
				list += "</select></div>"
				fmt.Fprint(w, list)
			}else{
				list += "<option value=''>No hay almacenes</option></select></div>"
				fmt.Fprint(w, list)
			}
		}
	}
}
//Función que va a mostrar un select de paises según el usuario (ROOT o normal)
func region_pais(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	user := r.FormValue("username")
	query, err := db.Query("SELECT id, entidad_id, padre_id FROM usuarios WHERE user = ?", user)
	if err != nil {
		Error.Println(err)
	}
	for query.Next() {
		var id, entidad, padre int
		err = query.Scan(&id, &entidad, &padre)
		if err != nil {
			Error.Println(err)
		}
		//entidad = 0 & padre = 0: es un usuario ROOT, puede ver cualquier paíss
		if entidad == 0{
			var list string
			query, err := db.Query("SELECT id, pais FROM pais")
			if err != nil {
				Error.Println(err)
			}
			list = "<div class='panel-heading'>País</div><div class='panel-body'><select id='pais' name='pais'>"
			if query.Next(){
				var id_pais int
				var name string
				err = query.Scan(&id_pais, &name)
				if err != nil {
					Error.Println(err)
				}
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_pais, name)
				for query.Next() {
					err = query.Scan(&id_pais, &name)
					if err != nil {
						Error.Println(err)
					}
					list += fmt.Sprintf("<option value='%d'>%s</option>", id_pais, name)
				}
				list += "</select></div>"
				fmt.Fprint(w, list)
			}else{
				list += "<option value=''>No hay paises</option></select></div>"
				fmt.Fprint(w, list)
			}
		}else{
			//usuario NORMAL: solo puede ver los paises que él creó
			var list string
			query, err := db.Query("SELECT id, pais FROM pais WHERE creador_id = ?", id)
			if err != nil {
				Error.Println(err)
			}
			list = "<div class='panel-heading'>País</div><div class='panel-body'><select id='pais' name='pais'>"
			if query.Next(){
				var id_pais int
				var name string
				err = query.Scan(&id_pais, &name)
				if err != nil {
					Error.Println(err)
				}
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_pais, name)
				for query.Next() {
					err = query.Scan(&id_pais, &name)
					if err != nil {
						Error.Println(err)
					}
					list += fmt.Sprintf("<option value='%d'>%s</option>", id_pais, name)
				}
				list += "</select></div>"
				fmt.Fprint(w, list)
			}else{
				list += "<option value=''>No hay paises</option></select></div>"
				fmt.Fprint(w, list)
			}
		}
	}
}
//Función que va a mostrar un select de regiones según el usuario (ROOT o normal)
func provincia_region(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	user := r.FormValue("username")
	query, err := db.Query("SELECT id, entidad_id, padre_id FROM usuarios WHERE user = ?", user)
	if err != nil {
		Error.Println(err)
	}
	for query.Next() {
		var id, entidad, padre int
		err = query.Scan(&id, &entidad, &padre)
		if err != nil {
			Error.Println(err)
		}
		//entidad = 0 & padre = 0: es un usuario ROOT, puede ver cualquier paíss
		if entidad == 0{
			var list string
			query, err := db.Query("SELECT id, region FROM region")
			if err != nil {
				Error.Println(err)
			}
			list = "<div class='panel-heading'>Región</div><div class='panel-body'><select id='region' name='region'>"
			if query.Next(){
				var id_region int
				var name string
				err = query.Scan(&id_region, &name)
				if err != nil {
					Error.Println(err)
				}
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_region, name)
				for query.Next() {
					err = query.Scan(&id_region, &name)
					if err != nil {
						Error.Println(err)
					}
					list += fmt.Sprintf("<option value='%d'>%s</option>", id_region, name)
				}
				list += "</select></div>"
				fmt.Fprint(w, list)
			}else{
				list += "<option value=''>No hay regiones</option></select></div>"
				fmt.Fprint(w, list)
			}
		}else{
			//usuario NORMAL: solo puede ver las regiones que él creó
			var list string
			query, err := db.Query("SELECT id, region FROM region WHERE creador_id = ?", id)
			if err != nil {
				Error.Println(err)
			}
			list = "<div class='panel-heading'>Región</div><div class='panel-body'><select id='region' name='region'>"
			if query.Next(){
				var id_region int
				var name string
				err = query.Scan(&id_region, &name)
				if err != nil {
					Error.Println(err)
				}
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_region, name)
				for query.Next() {
					err = query.Scan(&id_region, &name)
					if err != nil {
						Error.Println(err)
					}
					list += fmt.Sprintf("<option value='%d'>%s</option>", id_region, name)
				}
				list += "</select></div>"
				fmt.Fprint(w, list)
			}else{
				list += "<option value=''>No hay regiones</option></select></div>"
				fmt.Fprint(w, list)
			}
		}
	}
}
//Función que va a mostrar un select de provincias según el usuario (ROOT o normal)
func tienda_provincia(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	user := r.FormValue("username")
	query, err := db.Query("SELECT id, entidad_id, padre_id FROM usuarios WHERE user = ?", user)
	if err != nil {
		Error.Println(err)
	}
	for query.Next() {
		var id, entidad, padre int
		err = query.Scan(&id, &entidad, &padre)
		if err != nil {
			Error.Println(err)
		}
		//entidad = 0 & padre = 0: es un usuario ROOT, puede ver cualquier paíss
		if entidad == 0{
			var list string
			query, err := db.Query("SELECT id, provincia FROM provincia")
			if err != nil {
				Error.Println(err)
			}
			list = "<div class='panel-heading'>Provincia</div><div class='panel-body'><select id='provincia' name='provincia'>"
			if query.Next(){
				var id_prov int
				var name string
				err = query.Scan(&id_prov, &name)
				if err != nil {
					Error.Println(err)
				}
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_prov, name)
				for query.Next() {
					err = query.Scan(&id_prov, &name)
					if err != nil {
						Error.Println(err)
					}
					list += fmt.Sprintf("<option value='%d'>%s</option>", id_prov, name)
				}
				list += "</select></div>"
				fmt.Fprint(w, list)
			}else{
				list += "<option value=''>No hay provincias</option></select></div>"
				fmt.Fprint(w, list)
			}
		}else{
			//usuario NORMAL: solo puede ver las regiones que él creó
			var list string
			query, err := db.Query("SELECT id, region FROM region WHERE creador_id = ?", id)
			if err != nil {
				Error.Println(err)
			}
			list = "<div class='panel-heading'>Provincia</div><div class='panel-body'><select id='provincia' name='provincia'>"
			if query.Next(){
				var id_prov int
				var name string
				err = query.Scan(&id_prov, &name)
				if err != nil {
					Error.Println(err)
				}
				list += fmt.Sprintf("<option value='%d'>%s</option>", id_prov, name)
				for query.Next() {
					err = query.Scan(&id_prov, &name)
					if err != nil {
						Error.Println(err)
					}
					list += fmt.Sprintf("<option value='%d'>%s</option>", id_prov, name)
				}
				list += "</select></div>"
				fmt.Fprint(w, list)
			}else{
				list += "<option value=''>No hay provincias</option></select></div>"
				fmt.Fprint(w, list)
			}
		}
	}
}