package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func serverErrorPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	data := BaseData(r, "CorvusCrypto.com - internal server error")
	err := globalTemplate.ExecuteTemplate(w, "500", data)
	if err != nil {
		fmt.Println(err)
	}
}

func getRouter() *httprouter.Router {
	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir(Config.StaticPath))
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		data := BaseData(r, "CorvusCrypto.com - just another coder blog")

		latestPost, err := getLatestPost()
		if err != nil && err != ErrPostNotFound {
			serverErrorPage(w, r)
			fmt.Println(err)
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
		var posts []*Post

		last, err := strconv.Atoi(r.FormValue("last"))
		if err != nil {
			last = 0
		}
		q := strings.ToLower(strings.Trim(r.FormValue("q"), " "))
		if q != "" {
			searchTerms := strings.Split(q, " ")
			posts, err = searchPosts(searchTerms, last)
			if err != nil {
				w.WriteHeader(500)
				return
			}
		} else {
			posts, err = getAllPosts(last)
			if err != nil {
				w.WriteHeader(500)
				return
			}
		}

		data["Posts"] = posts
		err = globalTemplate.ExecuteTemplate(w, "posts", data)
		if err != nil {
			fmt.Println(err)
		}
	})
	router.GET("/posts/:postURL", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		url := p.ByName("postURL")
		post, err := getPostByURL(url)
		if err != nil {
			if err == ErrPostNotFound {
				w.WriteHeader(404)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		data := BaseData(r, "CorvusCrypto.com - "+strings.Title(post.Title))
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
