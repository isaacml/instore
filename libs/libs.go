package libs

import (
	"bytes"
	"fmt"
	"github.com/todostreaming/ratelimit"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"bufio"
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
	seleccion = "<option value='" + org + ":.:0'>...</option>"
	arr := strings.Split(resultado, "::")
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
	bitmap_parsed, err := strconv.ParseInt(bitmap_hex, 16, 32)
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
Cifrado: Funcion que cifra o descifra un fichero existente.
	origen:  fichero origen
	destino: fichero destino
	key: 	 patrón por el que va a ser cifrado.
Devuelve un error en caso de que algo no haya salido bien.
*/
func Cifrado(origen, destino string, key []byte) error {
	var fail error
	p := make([]byte, 8) //Va a contener el archivo origen en bloques de 8 bytes
	var container []byte //Va almacenar los datos del fichero de destino
	file, err := os.OpenFile(origen, os.O_RDONLY, 0666)
	if err != nil {
		fail = fmt.Errorf("Error en la apertura")
	}
	lector := bufio.NewReader(file)
	for {
		num, err := lector.Read(p)
		if err != nil {
			fail = fmt.Errorf("Fin de lectura")
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
	ioutil.WriteFile(destino, container, 0666)
	return fail
}
