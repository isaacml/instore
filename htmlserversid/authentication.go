package main

import ()

var username string
var password string

// This function could be used to access to a Database for user/pass authentication procedure
func authentication(user, pass string) bool {
	if user == "Isaac" && pass == "alabama" {
		return true
	} else {
		return false
	}
}
