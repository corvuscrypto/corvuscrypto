package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

func parsePostForm(r *http.Request) (post *Post, ok bool) {
	post.Body = r.PostFormValue("body")
	if post.Body == "" {
		return nil, false
	}
	var err error
	post.Date, err = time.Parse(r.PostFormValue("date"), time.RFC3339)
	if err != nil {
		return nil, false
	}
	if r.PostFormValue("publish") != "" {
		post.Publish = true
	} else {
		post.Publish = r.PostFormValue("publish") == "true"
	}
	post.Summary = r.PostFormValue("summary")
	if post.Summary == "" {
		return nil, false
	}
	post.Tags = strings.Split(r.PostFormValue("tags"), ",")
	if len(post.Tags) == 0 {
		return nil, false
	}
	//normalize all tags
	for i := range post.Tags {
		post.Tags[i] = strings.ToLower(strings.TrimSpace(post.Tags[i]))
	}
	post.Title = r.PostFormValue("title")
	if post.Title == "" {
		return nil, false
	}
	post.URL = strings.ToLower(strings.TrimSpace(r.PostFormValue("url")))
	valid, err := regexp.MatchString(`[0-9a-z\-]{10,}`, post.URL)
	if post.URL == "" || !valid {
		return nil, false
	}
	ok = true
	return
}

func editViewHandler(w http.ResponseWriter, r *http.Request, url string, errors bool) {
	data := BaseData(r, "Edit Post")
	var err error
	data["Error"] = errors
	data["Post"], err = GetPostByURL(url)
	if err != nil {
		http.Redirect(w, r, "/newpost", http.StatusFound)

	}
	err = globalTemplate.ExecuteTemplate(w, "editPost", data)
	if err != nil {
		fmt.Println(err)
	}
}

func getRouter() *httprouter.Router {

	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir(Config.StaticPath))
	router.GET("/", checkAuth(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		data := BaseData(r, "CMS index")
		err := globalTemplate.ExecuteTemplate(w, "cmsIndex", data)
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
	router.POST("/drafts/:url", checkAuth(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		postURL := p.ByName("url")
		post, ok := parsePostForm(r)
		if !ok {
			return
		}
		if err := UpdatePost(postURL, post); err != nil {
			editViewHandler(w, r, postURL, true)
			return
		}
		http.Redirect(w, r, r.URL.Path, http.StatusFound)
	}))
	router.GET("/drafts/:url", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		editViewHandler(w, r, p.ByName("url"), false)
	})
	router.GET("/login", loginView)
	router.POST("/login", login)
	return router
}

func main() {
	loadConfig()
	initCipher()
	compileTemplates()
	initDBSession()
	router := getRouter()
	server := new(http.Server)
	server.Handler = router
	server.Addr = ":8081"
	//server.ListenAndServeTLS("","")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
