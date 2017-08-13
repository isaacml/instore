package main

import (
	"fmt"
	//"github.com/isaacml/instore/libs"
	"net/http"
	//"strconv"
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
	// Recupera campos del form tanto GET como POST
	r.ParseForm() 
	//FORMATO DE SALIDA: usuario;nom_completo;contraseña;status
	var output string 
	//VARIABLE DE FORMULARIO
	user := r.FormValue("user")
	name_user := r.FormValue("name_user")
	pass := r.FormValue("pass")
	father := r.FormValue("padre")
	input_entidad := r.FormValue("input_entidad")
	//Comprobamos que ninguno de los campos esté vacio
	if user == "" || name_user == "" || pass == "" || input_entidad == "" {
		output = fmt.Sprintf("%s;%s;;<div class='form-group text-warning'>Hay campos vacios</div>", user, name_user)
	} else {
		var existe_usuario int
		//Comprobamos si existe o NO el usuario en base de datos
		err1 := db.QueryRow("SELECT count(*) FROM usuarios WHERE user = ?", user).Scan(&existe_usuario)
		if err1 != nil {
			Error.Println(err1)
		}
		//Usuario no existe, continuamos...
		if existe_usuario == 0 {
			var id_admin, bitmap int
			//Tomamos el identificador del padre
			err2 := db.QueryRow("SELECT id FROM usuarios WHERE user = ?", father).Scan(&id_admin)
			if err2 != nil {
				Error.Println(err2)
			}
			//Generar el bitmap de acciones en hexadecimal
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
			bitmap_hex := fmt.Sprintf("%x", bitmap)  //Se guarda el valor del bitmap en hexadecimal
			//Insertamos los datos en BD
			db_mu.Lock()
			_, err1 := db.Exec("INSERT INTO usuarios (`user`, `old_user`, `pass`, `nombre_completo`, `entidad_id`, `padre_id`, `bitmap_acciones`) VALUES (?,?,?,?,?,?,?)",
				user, user, pass, name_user, input_entidad, id_admin, bitmap_hex)
			db_mu.Unlock()
			if err1 != nil {
				Error.Println(err1)
				output = fmt.Sprintf(";;;<div class='form-group text-danger'>Fallo al añadir usuario</div>")
			} else {
				output = fmt.Sprintf(";;;<div class='form-group text-success'>Usuario añadido correctamente</div>")
			}
		} else {
			//ERROR: el usuario ya existe
			output = fmt.Sprintf(";%s;;<div class='form-group text-danger'>El usuario ya existe, prueba con otro</div></div>", name_user)
		}
	}
	fmt.Fprint(w, output)
}

//Función para dar de alta una nueva entidad
func entidad(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	username := r.FormValue("username")
	entidad := r.FormValue("entidad")
	var output string 
	if entidad == "" {
		output = "<div class='form-group text-warning'>El campo no puede estar vacio</div>"
	} else {
		query, err := db.Query("SELECT id, padre_id FROM usuarios WHERE user = ?", username)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			var id, padre_id, id_admin int
			err = query.Scan(&id, &padre_id)
			if err != nil {
				Error.Println(err)
			}
			//Hacemos un select para obtener el id del usuario super-admin
			err = db.QueryRow("SELECT id FROM usuarios WHERE padre_id = 0 AND entidad_id = 0").Scan(&id_admin)
			if err != nil {
				Error.Println(err)
			}
			//Si es un usuario super-admin o un usuario que tiene creador super-admin, le permitimos crear entidades
			if padre_id == 0 || padre_id == id_admin {
				timestamp := time.Now().Unix()
				db_mu.Lock()
				_, err1 := db.Exec("INSERT INTO entidades (`nombre`, `creador_id`, `timestamp`, `last_access`) VALUES (?, ?, ?, ?)", entidad, id, timestamp, timestamp)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
					output = "<div class='form-group text-danger'>Fallo al añadir entidad</div>"
				} else {
					output = "OK"
				}
			} else {
				output = "<div class='form-group text-danger'>Solo un usuario ROOT puede añadir una entidad</div>"
			}
		}
	}
	fmt.Fprint(w, output)
}

//Función para dar de alta un nuevo almacen
func almacen(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	username := r.FormValue("username")
	almacen := r.FormValue("almacen")
	entidad := r.FormValue("entidad")
	var output string
	if almacen == "" {
		output = "<div class='form-group text-warning'>El campo no puede estar vacio</div>"
	} else if entidad == "" {
		output = "<div class='form-group text-warning'>Debe haber almenos una entidad</div>"
	} else {
		query, err := db.Query("SELECT id, padre_id FROM usuarios WHERE user = ?", username)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			var id, padre_id, id_admin int
			err = query.Scan(&id, &padre_id)
			if err != nil {
				Error.Println(err)
			}
			//Hacemos un select para obtener el id del usuario super-admin
			err = db.QueryRow("SELECT id FROM usuarios WHERE padre_id = 0 AND entidad_id = 0").Scan(&id_admin)
			if err != nil {
				Error.Println(err)
			}
			//Si es un usuario super-admin o un usuario que tiene creador super-admin, le permitimos crear almacenes
			if padre_id == 0 || padre_id == id_admin {
				timestamp := time.Now().Unix()
				db_mu.Lock()
				_, err1 := db.Exec("INSERT INTO almacenes (`almacen`, `creador_id`, `timestamp`, `entidad_id`) VALUES (?, ?, ?, ?)", almacen, id, timestamp, entidad)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
					output = "<div class='form-group text-danger'>Fallo al añadir almacen</div>"
				} else {
					output = "OK"
				}
			} else {
				output = "<div class='form-group text-danger'>Solo un usuario ROOT puede añadir un almacen</div>"
			}
		}
	}
	fmt.Fprint(w, output)
}

//Función para dar de alta un nuevo pais
func pais(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	username := r.FormValue("username")
	almacen := r.FormValue("almacen")
	pais := r.FormValue("pais")
	var output string
	if pais == "" {
		output = "<div class='form-group text-warning'>El campo pais no puede estar vacio</div>"
	} else if almacen == "" {
		output = "<div class='form-group text-warning'>Debe haber almenos un almacen</div>"
	} else {
		query, err := db.Query("SELECT id, padre_id FROM usuarios WHERE user = ?", username)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			var id, padre_id, id_admin int
			err = query.Scan(&id, &padre_id)
			if err != nil {
				Error.Println(err)
			}
			//Hacemos un select para obtener el id del usuario super-admin
			err = db.QueryRow("SELECT id FROM usuarios WHERE padre_id = 0 AND entidad_id = 0").Scan(&id_admin)
			if err != nil {
				Error.Println(err)
			}
			//Si es un usuario super-admin o un usuario que tiene creador super-admin, le permitimos crear paises
			if padre_id == 0 || padre_id == id_admin {
				timestamp := time.Now().Unix()
				db_mu.Lock()
				_, err1 := db.Exec("INSERT INTO pais (`pais`, `creador_id`, `timestamp`, `almacen_id`) VALUES (?, ?, ?, ?)", pais, id, timestamp, almacen)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
					output = "<div class='form-group text-danger'>Fallo al añadir pais</div>"
				} else {
					output = "OK"
				}
			} else {
				output = "<div class='form-group text-danger'>Solo un usuario ROOT puede añadir un pais</div>"
			}
		}
	}
	fmt.Fprint(w, output)
}

//Función para dar de alta una nueva región
func region(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	username := r.FormValue("username")
	region := r.FormValue("region")
	pais := r.FormValue("pais")
	var output string
	if region == "" {
		output = "<div class='form-group text-warning'>El campo region no puede estar vacio</div>"
	} else if pais == "" {
		output = "<div class='form-group text-warning'>Debe haber almenos un almacen</div>"
	} else {
		query, err := db.Query("SELECT id, padre_id FROM usuarios WHERE user = ?", username)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			var id, padre_id, id_admin int
			err = query.Scan(&id, &padre_id)
			if err != nil {
				Error.Println(err)
			}
			//Hacemos un select para obtener el id del usuario super-admin
			err = db.QueryRow("SELECT id FROM usuarios WHERE padre_id = 0 AND entidad_id = 0").Scan(&id_admin)
			if err != nil {
				Error.Println(err)
			}
			//Si es un usuario super-admin o un usuario que tiene creador super-admin, le permitimos crear regiones
			if padre_id == 0 || padre_id == id_admin {
				timestamp := time.Now().Unix()
				db_mu.Lock()
				_, err1 := db.Exec("INSERT INTO region (`region`, `creador_id`, `timestamp`, `pais_id`) VALUES (?, ?, ?, ?)", region, id, timestamp, pais)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
					output = "<div class='form-group text-danger'>Fallo al añadir region</div>"
				} else {
					output = "OK"
				}
			} else {
				output = "<div class='form-group text-danger'>Solo un usuario ROOT puede añadir una region</div>"
			}
		}
	}
	fmt.Fprintf(w, output)
}

//Función para dar de alta una nueva provincia
func provincia(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // recupera campos del form tanto GET como POST
	username := r.FormValue("username")
	region := r.FormValue("region")
	provincia := r.FormValue("provincia")
	var output string
	if provincia == "" {
		output = "<div class='form-group text-warning'>El campo provincia no puede estar vacio</div>"
	} else if region == "" {
		output = "<div class='form-group text-warning'>Debe haber almenos una región</div>"
	} else {
		query, err := db.Query("SELECT id, padre_id FROM usuarios WHERE user = ?", username)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			var id, padre_id, id_admin int
			err = query.Scan(&id, &padre_id)
			if err != nil {
				Error.Println(err)
			}
			//Hacemos un select para obtener el id del usuario super-admin
			err = db.QueryRow("SELECT id FROM usuarios WHERE padre_id = 0 AND entidad_id = 0").Scan(&id_admin)
			if err != nil {
				Error.Println(err)
			}
			//Si es un usuario super-admin o un usuario que tiene creador super-admin, le permitimos crear provincias
			if padre_id == 0 || padre_id == id_admin {
				timestamp := time.Now().Unix()
				db_mu.Lock()
				_, err1 := db.Exec("INSERT INTO provincia (`provincia`, `creador_id`, `timestamp`, `region_id`) VALUES (?, ?, ?, ?)", provincia, id, timestamp, region)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
					output = "<div class='form-group text-danger'>Fallo al añadir provincia</div>"
				} else {
					output = "OK"
				}
			} else {
				output = "<div class='form-group text-danger'>Solo un usuario ROOT puede añadir una provincia</div>"
			}
		}
	}
	fmt.Fprint(w, output) 
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
	var output string
	if tienda == "" || address == "" || phone == "" {
		output = "<div class='form-group text-warning'>No pueden haber campos vacíos</div>"
	} else if provincia == "" {
		output = "<div class='form-group text-warning'>Debe haber almenos una provincia</div>"
	} else {
		query, err := db.Query("SELECT id, padre_id FROM usuarios WHERE user = ?", username)
		if err != nil {
			Error.Println(err)
		}
		for query.Next() {
			var id, padre_id, id_admin int
			err = query.Scan(&id, &padre_id)
			if err != nil {
				Error.Println(err)
			}
			//Hacemos un select para obtener el id del usuario super-admin
			err = db.QueryRow("SELECT id FROM usuarios WHERE padre_id = 0 AND entidad_id = 0").Scan(&id_admin)
			if err != nil {
				Error.Println(err)
			}
			//Si es un usuario super-admin o un usuario que tiene creador super-admin, le permitimos crear tiendas
			if padre_id == 0 || padre_id == id_admin {
				timestamp := time.Now().Unix()
				db_mu.Lock()
				_, err1 := db.Exec("INSERT INTO tiendas (`tienda`, `creador_id`, `timestamp`, `provincia_id`, `address`, `phone`, `extra`) VALUES (?, ?, ?, ?, ?, ?, ?)", tienda, id, timestamp, provincia, address, phone, extra)
				db_mu.Unlock()
				if err1 != nil {
					Error.Println(err1)
					output = "<div class='form-group text-danger'>Fallo al añadir tienda</div>"
				} else {
					output = "OK"
				}
			} else {
				output = "<div class='form-group text-danger'>Solo un usuario ROOT puede añadir una tienda</div>"
			}
		}
	}
	fmt.Fprint(w, output)
}
