package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func getRouter() *httprouter.Router {
	router := httprouter.New()

	router.GET("", checkAuth(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		err := globalTemplate.ExecuteTemplate(w, "drafts", nil)
		if err != nil {
			fmt.Println(err)
		}
	}))
	return router
}

func main() {

	server := new(http.Server)
	server.Addr = ":8081"
	server.ListenAndServeTLS("", "")
}
