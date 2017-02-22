package main

import (
	"github.com/isaacml/instore/libs"
	"net/url"
)

func main() {
	v := url.Values{}
	v.Add("nombre", "Isaac ")
	v.Add("pass", "alabaña")

	libs.ClienteUpload("/home/isaac/Documentos/img1.jpg", "http://localhost:9090/upload.cgi?"+v.Encode(), 100, 2)
}
