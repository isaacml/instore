package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func check_actions(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.FormValue("action") == "create_pub" {
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
			bitmap_parsed, err := strconv.ParseInt(bitmap_hex, 16, 32)
			if err != nil {
				Error.Println(err)
			}
			res := bitmap_parsed & 1
			fmt.Fprint(w, res)
		}
	}
	if r.FormValue("action") == "prog_mus" {
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
			bitmap_parsed, err := strconv.ParseInt(bitmap_hex, 16, 32)
			if err != nil {
				Error.Println(err)
			}
			res := bitmap_parsed & 2
			fmt.Fprint(w, res)
		}
	}
}
