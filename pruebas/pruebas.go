package main

import (
	"golang.org/x/text/encoding/charmap"
	"io"
	"os"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		// handle file open error
	}
	out, err := os.Create("output.txt")
	if err != nil {
		// handler error
	}

	r := charmap.ISO8859_1.NewDecoder().Reader(f)
	io.Copy(out, r)

	out.Close()
	f.Close()
}
