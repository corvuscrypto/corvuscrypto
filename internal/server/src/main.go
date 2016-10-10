package main

import (
	"net/http"
	"path/filepath"
)

//config for ALL the things
var Config struct {
	TemplatePath string
	StaticPath   string
	Assets       map[string]interface{}
}

func loadConfig() {
	Config.StaticPath, _ = filepath.Abs("../../../static")
	Config.TemplatePath, _ = filepath.Abs("../../templates")
	Config.Assets = map[string]interface{}{
		"CSS": map[string]string{
			"URL":     "/static/css/main.css",
			"Version": "1",
		},
	}
}

func main() {
	loadConfig()
	initializeDBSession()
	compileTemplates()
	router := getRouter()
	server := new(http.Server)
	server.Handler = router
	server.Addr = ":8080"
	server.ListenAndServe()
}
