package main

import "net/http"

func main() {
	loadConfig()
	createLogger()
	initializeDBSession()
	compileTemplates()
	router := getRouter()
	server := new(http.Server)
	server.Handler = router
	server.Addr = ":8080"
	server.ListenAndServe()
}
