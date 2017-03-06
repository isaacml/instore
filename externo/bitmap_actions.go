package main

import (
	"fmt"
	"github.com/isaacml/instore/libs"
	"net/http"
)

func bitmap_actions(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	query, err := db.Query("SELECT bitmap_acciones FROM usuarios WHERE user = ?", r.FormValue("user"))
	if err != nil {
		Error.Println(err)
	}
	for query.Next() {
		var bitmap_hex string
		err = query.Scan(&bitmap_hex)
		if err != nil {
			Error.Println(err)
		}
		prog_pub := libs.BitmapParsing(bitmap_hex, 1)    //res[0]
		prog_mus := libs.BitmapParsing(bitmap_hex, 2)    //res[1]
		prog_msg := libs.BitmapParsing(bitmap_hex, 4)    //res[2]
		add_mus := libs.BitmapParsing(bitmap_hex, 8)     //res[3]
		msg_auto := libs.BitmapParsing(bitmap_hex, 10)   //res[4]
		msg_normal := libs.BitmapParsing(bitmap_hex, 20) //res[5]
		//Pasamos los valores al html
		fmt.Fprintf(w, "%d;%d;%d;%d;%d;%d", prog_pub, prog_mus, prog_msg, add_mus, msg_auto, msg_normal)
	}
}
