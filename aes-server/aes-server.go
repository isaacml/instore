package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	cifrado("/home/isaac/Música/Alejandro.cif", "/home/isaac/Música/Nuevo.mp3", []byte{11, 22, 33, 44, 55, 66, 77, 88})
}

func cifrado(origen, destino string, key []byte) error {
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
	ioutil.WriteFile(destino, container, 0666)
	return fail
}
