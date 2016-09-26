package main

import "net/http"

func main() {
	server := new(http.Server)
	server.Addr = ":8081"
	server.ListenAndServeTLS("", "")
}
