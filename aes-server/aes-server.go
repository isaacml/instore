package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	cifrado("/home/isaac/Música/JR/Alejandro.xxx", "/home/isaac/Música/JR/Alejandro.mp3", []byte{11, 22, 33, 44, 55, 66, 77, 88})
	//cifrado("/home/isaac/Música/Tentandome.mp3", "/home/isaac/Música/Tentandome.xxx", []byte{11, 22, 33, 44, 55, 66, 77, 88})
	//cifrado("/home/isaac/Música/La Nina.mp3", "/home/isaac/Música/La Nina.xxx", []byte{11, 22, 33, 44, 55, 66, 77, 88})
	//cifrado("/home/isaac/Música/Mysterious.mp3", "/home/isaac/Música/Mysterious.xxx", []byte{11, 22, 33, 44, 55, 66, 77, 88})
	//cifrado("/home/isaac/Música/J Alvarez.mp3", "/home/isaac/Música/J Alvarez.xxx", []byte{11, 22, 33, 44, 55, 66, 77, 88})
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
