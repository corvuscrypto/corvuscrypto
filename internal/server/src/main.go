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
	if Config.CertFile != "" {
		server.ListenAndServeTLS(Config.CertFile, Config.KeyFile)
	}
	server.ListenAndServe()
}
