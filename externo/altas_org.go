package main

import (
	"fmt"
	//"github.com/isaacml/instore/libs"
	"net/http"
	"strconv"
	"time"
)

//MASCARAS
const (
	PROG_PUB = 1 << iota
	PROG_MUS
	PROG_MSG
	ADD_MUS
	MSG_AUTO
	MSG_NORMAL
)

//Variables de estado global
var bad, empty string

//Función para dar de alta nuevos usuarios
func alta_users(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	user := r.FormValue("user")
	name_user := r.FormValue("name_user")
	pass := r.FormValue("pass")
	padre := r.FormValue("padre")
	input_entidad := r.FormValue("input_entidad")

	if user == "" || name_user == "" || pass == "" {
		empty := "Los campos no pueden estar vacios"
		fmt.Fprintf(w, "<div class='form-group text-warning'>%s</div>", empty)
	} else {
		//Generar el bitmap de acciones en hexadecimal
		var bitmap int
		if r.FormValue("prog_pub") != "" {
			bitmap = PROG_PUB
		}
		if r.FormValue("prog_mus") != "" {
			bitmap = bitmap + PROG_MUS
		}
		if r.FormValue("prog_msg") != "" {
			bitmap = bitmap + PROG_MSG
		}
		if r.FormValue("add_mus") != "" {
			bitmap = bitmap + ADD_MUS
		}
		if r.FormValue("msg_auto") != "" {
			bitmap = bitmap + MSG_AUTO
		}
		if r.FormValue("msg_normal") != "" {
			bitmap = bitmap + MSG_NORMAL
		}

		//Aquí se guarda el valor del bitmap en hexadecimal
		bitmap_hex := fmt.Sprintf("%x", bitmap)
		//Seleccionamos el usuario y la entidad de un padre concreto
		query, err := db.Query("SELECT id, user, entidad_id FROM usuarios WHERE user = ?", padre)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			var id, ent_bd int
			var user_bd string
			err = query.Scan(&id, &user_bd, &ent_bd)
			if err != nil {
				Error.Println(err)
			}
			if ent_bd != 0 {
				bad := "Solo un ROOT puede añadir nuevos usuarios"
				fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
				return
			} else if user_bd == user { //Comprobamos que no hay dos Usuarios con el mismo nombre
				bad := "Fallo al añadir: el usuario ya existe"
				fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
			} else {
				var entidad int
				if input_entidad != "" {
					//tomamos el id_entidad, proporcionado por el select de formulario
					entidad, err = strconv.Atoi(input_entidad)
					if err != nil {
						Error.Println(err)
					}
				}
				db_mu.Lock()
				_, err1 := db.Exec("INSERT INTO usuarios (`user`, `old_user`, `pass`, `nombre_completo`, `entidad_id`, `padre_id`, `bitmap_acciones`) VALUES (?,?,?,?,?,?,?)",
					user, user, pass, name_user, entidad, id, bitmap_hex)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
					bad := "Fallo al añadir usuario"
					fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
				} else {
					good := "Usuario añadido correctamente"
					fmt.Fprintf(w, "<div class='form-group text-success'>%s</div>", good)
				}
			}
		}
	}
}

//Función para dar de alta una nueva entidad
func entidad(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	username := r.FormValue("username")
	entidad := r.FormValue("entidad")

	if entidad == "" {
		empty = "El campo no puede estar vacio"
		fmt.Fprintf(w, "<div class='form-group text-warning'>%s</div>", empty)
	} else {
		query, err := db.Query("SELECT id, entidad_id FROM usuarios WHERE user = ?", username)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			var id, entidad_id int
			err = query.Scan(&id, &entidad_id)
			if err != nil {
				Error.Println(err)
			}
			if entidad_id == 0 {
				timestamp := time.Now().Unix()
				db_mu.Lock()
				_, err1 := db.Exec("INSERT INTO entidades (`nombre`, `creador_id`, `timestamp`, `last_access`) VALUES (?, ?, ?, ?)", entidad, id, timestamp, timestamp)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
					bad = "Fallo al añadir entidad"
					fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
				} else {
					fmt.Fprint(w, "OK")
				}
			} else {
				bad = "Solo un usuario ROOT puede añadir una entidad"
				fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
			}
		}
	}
}

//Función para dar de alta un nuevo almacen
func almacen(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	username := r.FormValue("username")
	almacen := r.FormValue("almacen")
	entidad := r.FormValue("entidad")

	if almacen == "" {
		empty = "El campo no puede estar vacio"
		fmt.Fprintf(w, "<div class='form-group text-warning'>%s</div>", empty)
	} else if entidad == "" {
		empty = "Debe haber almenos una entidad"
		fmt.Fprintf(w, "<div class='form-group text-warning'>%s</div>", empty)
	} else {
		query, err := db.Query("SELECT id, entidad_id FROM usuarios WHERE user = ?", username)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			var id, entidad_id int
			err = query.Scan(&id, &entidad_id)
			if err != nil {
				Error.Println(err)
			}
			if entidad_id == 0 {
				timestamp := time.Now().Unix()
				db_mu.Lock()
				_, err1 := db.Exec("INSERT INTO almacenes (`almacen`, `creador_id`, `timestamp`, `entidad_id`) VALUES (?, ?, ?, ?)", almacen, id, timestamp, entidad)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
					bad = "Fallo al añadir almacen"
					fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
				} else {
					fmt.Fprint(w, "OK")
				}
			} else {
				bad = "Solo un usuario ROOT puede añadir un almacen"
				fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
			}
		}
	}
}

//Función para dar de alta un nuevo pais
func pais(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	username := r.FormValue("username")
	almacen := r.FormValue("almacen")
	pais := r.FormValue("pais")

	if pais == "" {
		empty = "El campo pais no puede estar vacio"
		fmt.Fprintf(w, "<div class='form-group text-warning'>%s</div>", empty)
	} else if almacen == "" {
		empty = "Debe haber almenos un almacen"
		fmt.Fprintf(w, "<div class='form-group text-warning'>%s</div>", empty)
	} else {
		query, err := db.Query("SELECT id, entidad_id FROM usuarios WHERE user = ?", username)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			var id, entidad_id int
			err = query.Scan(&id, &entidad_id)
			if err != nil {
				Error.Println(err)
			}
			if entidad_id == 0 {
				timestamp := time.Now().Unix()
				db_mu.Lock()
				_, err1 := db.Exec("INSERT INTO pais (`pais`, `creador_id`, `timestamp`, `almacen_id`) VALUES (?, ?, ?, ?)", pais, id, timestamp, almacen)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
					bad = "Fallo al añadir pais"
					fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
				} else {
					fmt.Fprint(w, "OK")
				}
			} else {
				bad = "Solo un usuario ROOT puede añadir un pais"
				fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
			}
		}
	}
}

//Función para dar de alta una nueva región
func region(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	username := r.FormValue("username")
	region := r.FormValue("region")
	pais := r.FormValue("pais")

	if region == "" {
		empty = "El campo region no puede estar vacio"
		fmt.Fprintf(w, "<div class='form-group text-warning'>%s</div>", empty)
	} else if pais == "" {
		empty = "Debe haber almenos un almacen"
		fmt.Fprintf(w, "<div class='form-group text-warning'>%s</div>", empty)
	} else {
		query, err := db.Query("SELECT id, entidad_id FROM usuarios WHERE user = ?", username)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			var id, entidad_id int
			err = query.Scan(&id, &entidad_id)
			if err != nil {
				Error.Println(err)
			}
			if entidad_id == 0 {
				timestamp := time.Now().Unix()
				db_mu.Lock()
				_, err1 := db.Exec("INSERT INTO region (`region`, `creador_id`, `timestamp`, `pais_id`) VALUES (?, ?, ?, ?)", region, id, timestamp, pais)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
					bad = "Fallo al añadir region"
					fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
				} else {
					fmt.Fprint(w, "OK")
				}
			} else {
				bad = "Solo un usuario ROOT puede añadir una region"
				fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
			}
		}
	}
}

//Función para dar de alta una nueva provincia
func provincia(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	username := r.FormValue("username")
	region := r.FormValue("region")
	provincia := r.FormValue("provincia")

	if provincia == "" {
		empty = "El campo provincia no puede estar vacio"
		fmt.Fprintf(w, "<div class='form-group text-warning'>%s</div>", empty)
	} else if region == "" {
		empty = "Debe haber almenos una región"
		fmt.Fprintf(w, "<div class='form-group text-warning'>%s</div>", empty)
	} else {
		query, err := db.Query("SELECT id, entidad_id FROM usuarios WHERE user = ?", username)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			var id, entidad_id int
			err = query.Scan(&id, &entidad_id)
			if err != nil {
				Error.Println(err)
			}
			if entidad_id == 0 {
				timestamp := time.Now().Unix()
				db_mu.Lock()
				_, err1 := db.Exec("INSERT INTO provincia (`provincia`, `creador_id`, `timestamp`, `region_id`) VALUES (?, ?, ?, ?)", provincia, id, timestamp, region)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
					bad = "Fallo al añadir provincia"
					fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
				} else {
					fmt.Fprint(w, "OK")
				}
			} else {
				bad = "Solo un usuario ROOT puede añadir una provincia"
				fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
			}
		}
	}
}

//Función para dar de alta una nueva tienda
func tienda(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	username := r.FormValue("username")
	tienda := r.FormValue("tienda")
	provincia := r.FormValue("provincia")
	address := r.FormValue("address")
	phone := r.FormValue("phone")
	extra := r.FormValue("extra")

	if tienda == "" || address == "" || phone == "" {
		empty = "No pueden haber campos vacíos"
		fmt.Fprintf(w, "<div class='form-group text-warning'>%s</div>", empty)
	} else if provincia == "" {
		empty = "Debe haber almenos una provincia"
		fmt.Fprintf(w, "<div class='form-group text-warning'>%s</div>", empty)
	} else {
		query, err := db.Query("SELECT id, entidad_id FROM usuarios WHERE user = ?", username)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			var id, entidad_id int
			err = query.Scan(&id, &entidad_id)
			if err != nil {
				Error.Println(err)
			}
			if entidad_id == 0 {
				timestamp := time.Now().Unix()
				db_mu.Lock()
				_, err1 := db.Exec("INSERT INTO tiendas (`tienda`, `creador_id`, `timestamp`, `provincia_id`, `address`, `phone`, `extra`) VALUES (?, ?, ?, ?, ?, ?, ?)", tienda, id, timestamp, provincia, address, phone, extra)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
					bad = "Fallo al añadir tienda"
					fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
				} else {
					fmt.Fprint(w, "OK")
				}
			} else {
				bad = "Solo un usuario ROOT puede añadir una tienda"
				fmt.Fprintf(w, "<div class='form-group text-danger'>%s</div>", bad)
			}
		}
	}
}
