package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func getRouter() *httprouter.Router {

	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir(Config.StaticPath))
	router.GET("/", checkAuth(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		err := globalTemplate.ExecuteTemplate(w, "cmsIndex", nil)
		if err != nil {
			fmt.Println(err)
		}
	}))
	router.GET("/dashboard", checkAuth(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	}))
	router.GET("/newpost", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		data := BaseData(r, "New Post")
		err := globalTemplate.ExecuteTemplate(w, "editPost", data)
		if err != nil {
			fmt.Println(err)
		}
	})
	//Not very RESTful but eh.
	router.POST("/drafts/save", nil)
	router.GET("/drafts/:id/edit", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		postID := strings.TrimSpace(p.ByName("id"))
		if postID == "" || len(postID) != 24 {
			http.Redirect(w, r, "/newpost", http.StatusFound)
			return
		}
		data := BaseData(r, "Edit Post")
		err := globalTemplate.ExecuteTemplate(w, "editPost", data)
		if err != nil {
			fmt.Println(err)
		}
	})
	router.GET("/login", loginView)
	router.POST("/login", login)
	return router
}

func main() {
	loadConfig()
	initCipher()
	compileTemplates()
	router := getRouter()
	server := new(http.Server)
	server.Handler = router
	server.Addr = ":8081"
	//server.ListenAndServeTLS("","")
	server.ListenAndServe()
}
