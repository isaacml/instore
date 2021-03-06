package libs

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/todostreaming/ratelimit"
)

/*
DownloadFile: función que se encarga de descargar un fichero desde una URL específica.
	url: http://www.dreamscene.org/wallpapers/Aquadream.jpg
	rutaFichero: /home/isaac/nuevasPruebas/new.jpg
	timeout: segundos que tiene para descargarse el fichero, 2sec
	bitrate: velocidad de bajada del fichero, 100Kb

La función devuelve el número de bytes y un error.
Si el proceso ha ido mal, los resultados devueltos serán 0, err(nombre del error).
*/
func DownloadFile(link, rutaFichero string, timeout time.Duration, bitrate float64) (bytes int64, err error) {
	var errR error
	var existe bool
	//Comprobar que existe el fichero(publicidad/mensaje)
	client := &http.Client{
		Timeout: timeout * time.Second,
	}
	resp, err := client.Get(link)
	if err != nil {
		errR = fmt.Errorf("URL: No puedo bajar la URL")
		return 0, errR
	}
	cod_status := resp.Status //Obtenemos el código de estado de la petición (404 Not Found / 202 OK)
	//Guardamos en una variable de estado si se encuenta o no el fichero
	if cod_status == "404 Not Found" {
		existe = false
	} else {
		existe = true
	}
	//Si existe procedemos a la descarga del fichero
	if existe == true {
		//Creamos un nuevo fichero en la ruta especificada
		escritor, err := os.Create(rutaFichero)
		if err != nil {
			errR = fmt.Errorf("CreateFile: No puedo crear el fichero")
			return 0, errR
		}
		size := resp.Header.Get("Content-Length")
		tamanio, err := strconv.Atoi(size) // tamaño del fichero cogido del Header()
		if err != nil {
			errR = fmt.Errorf("StringToInt: No se ha realizado la conversion")
			return 0, errR
		}
		bucket := ratelimit.NewBucketWithRate(bitrate*1024, 1*1024) // inicializando el bucket
		_, err = io.Copy(escritor, ratelimit.Reader(resp.Body, bucket))
		if err != nil {
			errR = fmt.Errorf("Copy: No puedo copiar el fichero")
			return 0, errR
		}
		defer resp.Body.Close()
		imagen, err := os.Stat(rutaFichero) // tomamos la informacion del fichero recientemente bajado
		if err != nil {
			errR = fmt.Errorf("Stat: Error al acceder al estado")
			return 0, errR
		}
		if imagen.Size() != int64(tamanio) { // comparamos tamaños
			errR = fmt.Errorf("CompareSize: Se ha producido un error en la copia del fichero")
			return 0, errR
		} else {
			bytes = imagen.Size()
			return bytes, nil
		}
	} else {
		errR = fmt.Errorf("FileExt: El servidor externo no tiene ese fichero")
		return 0, errR
	}
}

/*
DownloadPage: función encargada de descargar un página específica.
	url: http://www.streamrus.com/en/index.php
	timeout: segundos que tiene para descargarse el fichero, 2sec
La función devuelve un io.Reader y un error, este io.Reader contiene el contenido de la página.
*/
func DownloadPage(url string, timeout time.Duration) (resp io.Reader, err error) {
	client := &http.Client{
		Timeout: timeout * time.Second,
	}
	respuesta, err := client.Get(url)
	if err != nil {
		err = fmt.Errorf("URL: No puedo bajar la URL")
		return
	}
	res, err := ioutil.ReadAll(respuesta.Body) //leemos el body
	if err != nil {
		err = fmt.Errorf("Body: No puedo leer el Body")
		return
	}
	defer respuesta.Body.Close()
	resp = strings.NewReader(string(res))

	return
}

/*
ClienteUpload: esta función va a realizar la subida de un fichero por parte del cliente, a un servidor específico.
	filename: /home/isaac/Documentos/img1.jpg
	targetUrl: url del servidor, junto con sus parametros, http://localhost:9090/upload.cgi?parametro1=par1&parametro2=par2
	bitrate: velocidad de bajada del fichero, 100Kb
	timeout: segundos que tiene para descargarse el fichero, 2sec
La función devuelve un err=nil, si todo ha ido bien o un err = 'Nombre del error' si ha ido mal.
*/
func ClienteUpload(filename, targetUrl string, bitrate float64, timeout time.Duration) (resp *http.Response, err error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("file", filename)
	if err != nil {
		err = fmt.Errorf("error writing to buffer")
		return
	}
	//openfile
	fh, err := os.Open(filename)
	if err != nil {
		err = fmt.Errorf("error opening file")
		return
	}
	bucket := ratelimit.NewBucketWithRate(bitrate*1024, 1*1024) // inicializando el bucket
	//iocopy
	_, err = io.Copy(fileWriter, ratelimit.Reader(fh, bucket))
	if err != nil {
		err = fmt.Errorf("error copying file")
		return
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	client := &http.Client{
		Timeout: timeout * time.Second,
	}
	res, err := client.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		err = fmt.Errorf("failure by the post")
		return
	}
	defer res.Body.Close()
	res.Header.Set("Content-Type", contentType)
	resp = res
	err = nil
	return
}

/*
SendFile: esta función determina que fichero se va a publicar por NATS
	ruta: nombre de la ruta del fichero a subir --> /home/isaac/Documentos/
	name: nombre del fichero con su extension --> google.jpg
La función devuelve un array de bytes con el nombre y el contenido binario del fichero.
*/
func SendFile(ruta, name string) []byte {
	file, err1 := ioutil.ReadFile(ruta + name)
	if err1 != nil {
		fmt.Println("No puedo leer el fichero")
	}
	arr := []byte(fmt.Sprintf("SENDFILE[filename]%s[binary]%s", name, file))
	return arr
}

/*
GetFile: esta función va a tomar los parámetros [filename], [binary] generados por SendFile(),
para generar el fichero en una nueva ruta.
	ruta: nombre de la ruta del fichero donde se va a copiar --> /home/isaac/Documentos/
	data: datos binarios del fichero
*/
func GetFile(ruta, data string) {
	arr := strings.Split(data, "[filename]")
	if arr[0] == "SENDFILE" {
		cmd := strings.Split(arr[1], "[binary]")
		//Convierto el string a un array de bytes
		binario := []byte(cmd[1])
		ioutil.WriteFile(ruta+cmd[0], binario, 0666)
	}
}

/*
GenerateFORM: esta función va a tomar un link o url a donde se va enviar una serie de parámetros.
	link: link o URL donde se envian los parámetros --> http://localhost:9999/login.cgi
	param: de este tipo --> "user;"+r.FormValue(name_username)
	donde user es el nombre del parámetro y r.FormValue(name_username) es el parámetro extraido del formulario.
*/
func GenerateFORM(link string, param ...interface{}) (r string) {
	v := url.Values{}
	for _, val := range param {
		values := strings.Split(val.(string), ";")
		v.Set(values[0], values[1])
	}
	client := &http.Client{}
	res, err := client.PostForm(link, v)
	if err != nil {
		err = fmt.Errorf("failure by the post")
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("fail to read")
		return
	}
	return string(body)
}

/*
GenerateSelectOrg: esta función va a generar los selects de destino para las distintas organizaciones.
	resultado: es el resultado obtenido de los selects. Formato: id;org1;org2::
	org: nos indica a que organización pertenece el select. Formato: entidad:.:1
Nos devuelve el select formado correctamente y un estado que contiene la organización a la que pertenece.
*/
func GenerateSelectOrg(resultado, org string) (seleccion, estado string) {
	var arr_org []string
	var clear_res string
	if strings.Contains(resultado, "@@") {
		res := strings.Split(resultado, "@@")
		clear_res = res[1]
	} else {
		clear_res = resultado
	}
	seleccion = "<option value='" + org + ":.:0'>...</option>"
	arr := strings.Split(clear_res, "::")
	for _, val := range arr {
		if val != "" {
			arr_org = strings.Split(val, ";")
			seleccion += fmt.Sprintf("<option value='"+org+":.:%s'>%s</option>", arr_org[0], arr_org[1])
		}
	}
	estado = arr_org[2]
	return
}

/*
BackDestOrg: esta función nos va a permitir retroceder en un destino.
	estado_destino: Estado del destino actual. Formato: entidad.almacen.pais.region.provincia.*
	num_backs: Nos indica el numero de saltos que queremos dar hacia atrás.
Nos devuelve un string con el nuevo estado, el cual será enviado a la base de datos.
*/
func BackDestOrg(estado_destino string, num_backs int) (resultado string) {
	var res string
	separator := strings.Split(estado_destino, ".")
	arr := separator[:len(separator)-num_backs]
	for _, v := range arr {
		res += v + "."
	}
	resultado = res
	return
}

/*
DeleteSplitsChars: Función que elimina los puntos, puntos comas y dobles puntos.
	cadena: Se le pasa el valor de un input. Ex: r.FormValue("user")
Nos devuelve un string limpio.
*/
func DeleteSplitsChars(cadena string) (resultado string) {
	var correct_res string
	r := strings.NewReplacer(".", "", ":", "", ";", "")
	if strings.Contains(cadena, ".") || strings.Contains(cadena, ":") || strings.Contains(cadena, ";") {
		correct_res = r.Replace(cadena)
	} else {
		correct_res = cadena
	}
	resultado = correct_res
	return
}

/*
BitmapParsing: Parsea el valor tomado de la base de datos a INT y le aplica la máscara correspondiente.
	bitmap_hex: Contiene el valor del bitmap parseado a INT.
	mascara: Mascara aplicada al valor del bitmap_hex.
Devuelve:  = 0  accion activada; != 0 no tiene la acción activada
*/
func BitmapParsing(bitmap_hex string, mascara int64) (res int64) {
	bitmap_parsed, err := strconv.ParseInt(bitmap_hex, 16, 16)
	if err != nil {
		err = fmt.Errorf("fail to parsing")
	}
	res = bitmap_parsed & mascara
	return
}

/*
DomainGenerator: genera los dominios necesarios a partir del dominio de la tienda.
	dom_tienda: entidad.almacen.pais.region.provincia.tienda.
Devuelve: un array de strings con todos los dominios.
*/
func DomainGenerator(dom_tienda string) []string {
	var dom_provincia, dom_region, dom_pais, dom_almacen, dom_entidad string
	var asterisco = "*"
	var list_dom []string

	sep := strings.Split(dom_tienda, ".")
	//entidad.almacen.pais.region.provincia.tienda
	list_dom = append(list_dom, dom_tienda)
	//entidad.almacen.pais.region.provincia.*
	borrado1 := sep[:len(sep)-1]
	for _, v1 := range borrado1 {
		dom_provincia += v1 + "."
	}
	dom_provincia += asterisco
	list_dom = append(list_dom, dom_provincia)
	//entidad.almacen.pais.region.*
	borrado2 := sep[:len(sep)-2]
	for _, v2 := range borrado2 {
		dom_region += v2 + "."
	}
	dom_region += asterisco
	list_dom = append(list_dom, dom_region)
	//entidad.almacen.pais.*
	borrado3 := sep[:len(sep)-3]
	for _, v3 := range borrado3 {
		dom_pais += v3 + "."
	}
	dom_pais += asterisco
	list_dom = append(list_dom, dom_pais)
	//entidad.almacen.*
	borrado4 := sep[:len(sep)-4]
	for _, v4 := range borrado4 {
		dom_almacen += v4 + "."
	}
	dom_almacen += asterisco
	list_dom = append(list_dom, dom_almacen)
	//entidad.*
	borrado5 := sep[:len(sep)-5]
	for _, v5 := range borrado5 {
		dom_entidad += v5 + "."
	}
	dom_entidad += asterisco
	list_dom = append(list_dom, dom_entidad)

	return list_dom
}

/*
RemoveDuplicates: Borra los datos duplicados en un slice.
	domains: contiene el listado de dominios(puede contener duplicados).
Devuelve un slice limpio sin duplciados
*/

func RemoveDuplicates(domains []string) []string {
	encountered := map[string]bool{}
	result := []string{}
	for value := range domains {
		if encountered[domains[value]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[domains[value]] = true
			// Append to result slice.
			result = append(result, domains[value])
		}
	}
	return result
}

/*
LimpiarMatriz: Limpia de carácteres nulos la matriz salida de windows
	matriz: nombre de la matriz
*/
func LimpiarMatriz(matriz []byte) []byte {
	var matriz_limpiada []byte
	for _, v := range matriz {
		if v != 0 {
			matriz_limpiada = append(matriz_limpiada, v)
		}
	}
	return matriz_limpiada
}

/*
RemoveIndex: Borrar espacios vacios dentro de un array
	arr: nombre del array
	index: identificador del valor vacio que queremos borrar
Devuelve un array limpio, sin valores nulos
*/
func RemoveIndex(arr []string, index int) []string {
	return append(arr[:index], arr[index+1:]...)
}

/*
Cifrado: Funcion que cifra o descifra un fichero existente.
	origen:  fichero origen
	destino: fichero destino
	key: 	 patrón por el que va a ser cifrado.
Devuelve un error en caso de que algo no haya salido bien(en formato err y string).
*/
func Cifrado(origen, destino string, key []byte) (error, string) {
	var fail error
	var outString string
	p := make([]byte, 8) //Va a contener el archivo origen en bloques de 8 bytes
	var container []byte //Va almacenar los datos del fichero de destino
	file, err := os.OpenFile(origen, os.O_RDONLY, 0666)
	if err != nil {
		fail = fmt.Errorf("Error en la apertura")
		outString = "BAD"
	}
	lector := bufio.NewReader(file)
	for {
		num, err := lector.Read(p)
		if err != nil {
			fail = fmt.Errorf("Fin de lectura")
			outString = "END"
			break
		}
		if num <= 0 {
			break
		} else {
			for i := 0; i < num; i++ {
				container = append(container, p[i]^key[i])
			}
		}
	}
	//Escribimos el fichero
	err = ioutil.WriteFile(destino, container, 0666)
	if err != nil {
		fail = fmt.Errorf("Error en escritura")
		outString = "BAD"
	} else {
		outString = "GOOD"
	}
	return fail, outString
}

/*
FileCopier: Esta funcion es usada internamente por la tienda y se encarga de copiar directorios de musica o cifrados.
	contenedor:  es el array que contiene todas las carpetas que se quieren copiar.
	destino: ruta donde se guardaran los directorios de musica seleccionados "C:\instore\\Music\"
*/
func FileCopier(contenedor []string, destino string) {
	//Obtenemos todos los ficheros para la copia y los procesamos.
	for _, val := range contenedor {
		canciones := strings.Split(val, "\\")
		//Nombre la carpeta que contiene la cancion
		song_dir := canciones[len(canciones)-2]
		//ruta de destino + nombre de directorio de la cancion
		all_song_dir := destino + song_dir
		//Comprobamos si existe esa carpeta en el directorio de Musica(C:\instore\\Music)
		ext_dir := Existencia(all_song_dir)
		if ext_dir == false {
			//Si no existe, se crea
			os.Mkdir(all_song_dir, os.FileMode(0777))
		}
		//Comprobamos si existe dicha cancion en el directorio de Musica
		song_name := canciones[len(canciones)-1]
		//ruta de destino y directorio de cancion + cancion
		all_song_name := all_song_dir + "\\" + song_name
		ext_song := Existencia(all_song_name)
		if ext_song == false {
			//Abrimos el fichero origen
			original, err := os.Open(val)
			if err != nil {
				err = fmt.Errorf("FileCopier: fail to open origen")
				return
			}
			//Creamos el fichero destino
			copia, err := os.Create(all_song_name)
			if err != nil {
				err = fmt.Errorf("FileCopier: fail to open destino")
				return
			}
			//Realizamos la copia
			io.Copy(copia, original)
		}
	}
}

/*
Existencia: Comprueba la existencia de un fichero o directorio
	name:  ruta completa del fichero o directorio
Devuelve un bool con el resultado
*/
func Existencia(ruta string) bool {
	var existe bool
	_, err := os.Stat(ruta)
	if err != nil {
		if os.IsNotExist(err) {
			existe = false
		}
	} else {
		existe = true
	}
	return existe
}

/*
MusicToPlay: Esta función determina los ficheros que va a reproducir el player de la tienda.
	ruta:  Ruta del directorio que va a contener los ficheros de música
	st:    Estado de la música cifrada (0: solo cif / 1: cif y no cif)
Devuelve un array con todos los ficheros a reproducir.
*/
func MusicToPlay(ruta string, st int) []string {
	var arr_music []string
	var cmd []byte
	var gen_bat string
	var msg_file *os.File
	msg_file, _ = os.Create("music_to_play.bat")
	defer msg_file.Close()
	if st == 0 {
		//Se obtienen los ficheros del directorio y subdirectorios (solo música cif)
		gen_bat = "dir /s /b \"" + ruta + "*.xxx\""
	} else if st == 1 {
		//Se obtienen los ficheros del directorio y subdirectorios (cif / no cif)
		gen_bat = "dir /s /b \"" + ruta + "*.mp3\" & dir /s /b \"" + ruta + "*.xxx\""
	}
	msg_file.WriteString(gen_bat)
	//Una vez creado el fichero, lo ejecutamos
	cmd, _ = exec.Command("cmd", "/c", "music_to_play.bat").CombinedOutput()
	ficheros := strings.Split(string(cmd), "\r\n")
	for _, val := range ficheros {
		if strings.Contains(val, ruta) {
			if !strings.Contains(val, "dir /s /b") {
				//Se agregan cada una de las canciones al contenedor de música
				arr_music = append(arr_music, val)
			}
		}
	}
	return arr_music
}

/*
MainDomain: Obtener el dominio principal de la tienda
	file: Nombre del fichero que contiene el dominio (configShop)
Devuelve el dominio principal de la tienda al completo
*/
func MainDomain(filename string) string {
	var dom string
	fr, err := os.Open(filename)
	defer fr.Close()
	if err == nil {
		reader := bufio.NewReader(fr)
		for {
			linea, rerr := reader.ReadString('\n')
			if rerr != nil {
				break
			}
			linea = strings.TrimRight(linea, "\r\n")
			item := strings.Split(linea, " = ")
			//se obtiene el dominio principal
			if item[0] == "shopdomain" {
				dom = item[1]
			}
		}
	}
	return dom
}

/*
MyCurrentDate: generamos una fecha actual propia
Salida --> 20070405
*/
func MyCurrentDate() string {
	fecha_actual := time.Now()
	string_fecha := fmt.Sprintf("%4d%02d%02d", fecha_actual.Year(), int(fecha_actual.Month()), fecha_actual.Day())
	return string_fecha
}

/*
MyCurrentClock: generamos una hora actual propia
Salida --> 10:09
*/
func MyCurrentClock() string {
	//Obtenemos la hora local
	hh, mm, _ := time.Now().Clock()
	clock := fmt.Sprintf("%02d:%02d", hh, mm)
	return clock
}

/*
DaysIn: los dias que tiene un mes específico
	m: mes
	year: año
Salida --> numero de dias
*/
func DaysIn(m time.Month, year int) int64 {
	return int64(time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day())
}

/*
ToInt: convierte un string numérico en un entero int
	cant: valor de la cadena
Salida --> valor convertido a entero
*/
func ToInt(cant string) (res int) {
	res, _ = strconv.Atoi(cant)
	return
}

/*
FechaCreacion: Función para obtener la fecha y hora de creación
	timestamp: cantidad de segundos
Nos devuelve la fecha y la hora
*/
func FechaCreacion(timestamp int64) string {
	var anio, mes, dia string
	localtime := time.Unix(timestamp, 0)
	toString := localtime.String()
	inSlice := strings.Split(toString, " ")
	fmt.Sscanf(inSlice[0], "%4s-%2s-%2s", &anio, &mes, &dia)
	out := fmt.Sprintf("%s/%s/%s - %s", dia, mes, anio, inSlice[1])
	return out
}

/*
FechaSQLtoNormal: pasamos de un formato de fecha SQL a uno normal
	fecha: la fecha en este formato --> 20070405
Salida de la fecha convertida: 05/04/2007
*/
func FechaSQLtoNormal(fecha string) string {
	var anio, mes, dia string
	fmt.Sscanf(fecha, "%4s%2s%2s", &anio, &mes, &dia)
	out := fmt.Sprintf("%s/%s/%s", dia, mes, anio)
	return out
}

/*
FechaNormaltoSQL: pasamos de un formato de fecha normal al formato SQL
	fecha: la fecha en este formato --> 05/04/2007
Salida de la fecha convertida: 20070405
*/
func FechaNormaltoSQL(fecha string) string {
	var anio, mes, dia string
	fmt.Sscanf(fecha, "%2s/%2s/%4s", &dia, &mes, &anio)
	out := fmt.Sprintf("%s%s%s", anio, mes, dia)
	return out
}

/*
LoadSettingsLin: esta función va a abrir un fichero, leer los datos que contiene y guardarlos en un mapa (PARA LINUX)
	filename: ruta completa donde se encuentra nuestro fichero("serverext.reg")
	mapa: donde guardamos los datos extraidos del fichero
*/
func LoadSettingsLin(filename string, mapa map[string]string) {
	fr, err := os.Open(filename)
	defer fr.Close()
	if err == nil {
		reader := bufio.NewReader(fr)
		for {
			linea, rerr := reader.ReadString('\n')
			if rerr != nil {
				break
			}
			linea = strings.TrimRight(linea, "\n")
			item := strings.Split(linea, " = ")
			if len(item) == 2 {
				mapa[item[0]] = item[1]
			}
		}
	}
}

/*
LoadDomains: abre el fichero de dominios, lee los dominios que contiene y los guarda en un mapa
	filename: ruta donde se encuentra nuestro fichero(configshop.reg)
*/
func LoadDomains(filename string, arr map[string]string) {
	cont := 1
	fr, err := os.Open(filename)
	defer fr.Close()
	if err == nil {
		reader := bufio.NewReader(fr)
		for {
			linea, rerr := reader.ReadString('\n')
			if rerr != nil {
				break
			}
			linea = strings.TrimRight(linea, "\r\n")
			item := strings.Split(linea, " = ")
			if len(item) == 2 {
				if _, ok := arr[item[0]]; ok {
					clave := fmt.Sprintf("%s%d", item[0], cont)
					arr[clave] = item[1]
					cont++
				} else {
					arr[item[0]] = item[1]
				}
			}
		}
	}
}

func MostrarHoras(hora string) string {
	var str string
	if hora == "" {
		for i := 0; i <= 23; i++ {
			str += fmt.Sprintf("<option value='%02d'>%02d</option>", i, i)
		}
	} else {
		arr_hora := strings.Split(hora, ":")
		hora_sql := ToInt(arr_hora[0])
		for i := 0; i <= 23; i++ {
			if hora_sql == i {
				str += fmt.Sprintf("<option value='%02d' selected>%02d</option>", i, i)
			} else {
				str += fmt.Sprintf("<option value='%02d'>%02d</option>", i, i)
			}
		}
	}
	return str
}

func MostrarMinutos(hora string) string {
	var str string
	if hora == "" {
		for i := 0; i <= 59; i++ {
			str += fmt.Sprintf("<option value='%02d'>%02d</option>", i, i)
		}
	} else {
		arr_hora := strings.Split(hora, ":")
		mins_sql := ToInt(arr_hora[1])
		for i := 0; i <= 59; i++ {
			if mins_sql == i {
				str += fmt.Sprintf("<option value='%02d' selected>%02d</option>", i, i)
			} else {
				str += fmt.Sprintf("<option value='%02d'>%02d</option>", i, i)
			}
		}
	}
	return str
}

/*
Hour2min: toma una hora (hh:mm) y la convierte en minutos
	hh: hora
	mm: minutos
Devuelve la cantidad todal de minutos
*/
func Hour2min(hh int, mm int) int {
	mins := hh*60 + mm
	return mins
}

/*
Min2hour: convierte minutos totales en una hora(hh:mm)
	mm: minutos totales
Devuelve la hora y minutos correspondientes
*/
func Min2hour(mm int) (int, int) {
	hh := int(mm / 60)
	min := mm % 60
	return hh, min
}

func St_Prog_Music(db *sql.DB, db_mu sync.Mutex) (string, error) {
	var err error
	var cont int
	var st_prog string
	db_mu.Lock()
	db.QueryRow("SELECT count(estado) FROM st_prog_music").Scan(&cont)
	db_mu.Unlock()
	if cont == 0 {
		st_prog = ""
	} else {
		err = db.QueryRow("SELECT estado FROM st_prog_music").Scan(&st_prog)
		if err != nil {
			err = fmt.Errorf("Fail Prog: fail to read status prog")
		}
	}
	return st_prog, err
}
