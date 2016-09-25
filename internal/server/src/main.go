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
	Config.TemplatePath, _ = filepath.Abs("../../templates")
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

		latestPost, err := getLatestPost()
		if err != nil {
			w.WriteHeader(500)
			return
		}
		data["LatestPost"] = latestPost
		err = globalTemplate.ExecuteTemplate(w, "index", data)
		if err != nil {
			fmt.Println(err)
		}
	})
	router.GET("/posts", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		data := BaseData(r, "CorvusCrypto.com - Posts")
		err := globalTemplate.ExecuteTemplate(w, "posts", data)
		if err != nil {
			fmt.Println(err)
		}
	})
	router.GET("/posts/:postURL", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		url := p.ByName("postURL")
		data := BaseData(r, "CorvusCrypto.com - Posts")
		post, err := getPostByURL(url)
		if err != nil {
			if err == ErrPostNotFound {
				w.WriteHeader(404)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		data["Post"] = post
		err = globalTemplate.ExecuteTemplate(w, "postFull", data)
		if err != nil {
			fmt.Println(err)
		}
	})
	router.GET("/about", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		data := BaseData(r, "CorvusCrypto.com - About Me")
		err := globalTemplate.ExecuteTemplate(w, "about", data)
		if err != nil {
			fmt.Println(err)
		}
	})
	return router
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
