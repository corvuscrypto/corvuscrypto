package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
)

//Config for ALL the things
var Config struct {
	TemplatePath  string
	StaticPath    string
	Assets        map[string]interface{}
	CertFile      string
	PrivateKey    string
	OwnerUsername string
	//When loaded, pw will be hashed just in case somehow someone gets
	//access to the memory (uuuuunlikely but just in case)
	OwnerPassword []byte
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
	//TODO add in hashing of config password

}
func getRouter() *httprouter.Router {

	router := httprouter.New()
	router.GET("/drafts", checkAuth(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		err := globalTemplate.ExecuteTemplate(w, "drafts", nil)
		if err != nil {
			fmt.Println(err)
		}
	}))
	router.GET("/dashboard", checkAuth(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	}))
	router.POST("/login", login)
	return router
}

func main() {
	loadConfig()
	initCipher()
	router := getRouter()
	server := new(http.Server)
	server.Handler = router
	server.Addr = ":8081"
	//server.ListenAndServeTLS("","")
	server.ListenAndServe()
}
