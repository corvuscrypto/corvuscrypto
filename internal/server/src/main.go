package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
)

//config for ALL the things
var Config struct {
	TemplatePath string
	StaticPath   string
	Assets       map[string]interface{}
}

func loadConfig() {
	Config.StaticPath, _ = filepath.Abs("../../../static")
	Config.TemplatePath = Config.StaticPath + "/html"
	Config.Assets = map[string]interface{}{
		"CSS": map[string]string{
			"URL":     "/static/css/main.css",
			"Version": "1",
		},
	}
}

func getRouter() *httprouter.Router {
	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir(Config.StaticPath))
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		data := BaseData(r, "CorvusCrypto.com - just another coder blog")

		data["LatestPost"] = struct {
			Title,
			Summary string
		}{
			"Test title for great justice!",
			"Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
		}

		err := globalTemplate.ExecuteTemplate(w, "index", data)
		if err != nil {
			fmt.Println(err)
		}
	})
	return router
}

func main() {
	loadConfig()
	compileTemplates()
	router := getRouter()
	server := new(http.Server)
	server.Handler = router
	server.Addr = ":8080"
	server.ListenAndServe()
}
