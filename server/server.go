package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {

	http.HandleFunc("/upload.cgi", upload)

	log.Fatal(http.ListenAndServe(":9090", nil)) // Servidor HTTP multihilo
}

// upload logic
func upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("nombre")
	if err != nil {
		fmt.Printf("Error Multipart: %s\n", err)
		return
	}
	//fmt.Printf("%v-%v\n", r.FormValue("nombre"), r.FormValue("pass"))
	defer file.Close()
	fmt.Fprintf(w, "%v", handler.Header)
	nom := strings.Split(handler.Filename, "/")
	f, err := os.OpenFile("/home/isaac/nuevas pruebas/"+nom[4], os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("Error server upload: %s\n", err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
}
