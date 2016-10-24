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
	post = &Post{}
	ok = true
	post.Body = r.PostFormValue("body")
	if post.Body == "" {
		fmt.Println("body")
		ok = false
	}
	var err error
	post.Date, err = time.Parse("01/02/2006", r.PostFormValue("date"))
	if err != nil {
		post.Date = time.Now()
		ok = false
	}
	if r.PostFormValue("unpublish") != "" {
		post.Publish = false
	} else if r.PostFormValue("publish") != "" {
		post.Publish = true
	} else {
		post.Publish = r.PostFormValue("published") == "true"
	}
	post.Summary = r.PostFormValue("summary")
	if post.Summary == "" {
		ok = false
	}
	post.Tags = strings.Split(r.PostFormValue("tags"), ",")
	if len(post.Tags) == 0 {
		ok = false
	}
	//normalize all tags
	for i := range post.Tags {
		post.Tags[i] = strings.ToLower(strings.TrimSpace(post.Tags[i]))
	}
	post.Title = r.PostFormValue("title")
	if post.Title == "" {
		ok = false
	}
	post.URL = strings.ToLower(strings.TrimSpace(r.PostFormValue("url")))
	valid, err := regexp.MatchString(`[0-9a-z\-]{10,}`, post.URL)
	if post.URL == "" || !valid {
		ok = false
	}

	post.SearchTags = append(post.SearchTags, strings.Split(strings.ToLower(post.Title), " ")...)
	return
}

func getRouter() *httprouter.Router {

	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir(Config.StaticPath))
	router.GET("/", checkAuth(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		data := BaseData(r, "CMS index")
		var err error
		data["Posts"], err = GetPosts(true)
		err = globalTemplate.ExecuteTemplate(w, "cmsIndex", data)
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
	router.POST("/newpost", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		post, ok := parsePostForm(r)
		if !ok {
			data := BaseData(r, "New Post")
			data["Error"] = true
			data["Post"] = post
			err := globalTemplate.ExecuteTemplate(w, "editPost", data)
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		InsertNewPost(post)
		http.Redirect(w, r, "/", http.StatusFound)
	})
	//Not very RESTful but eh.
	router.POST("/posts/:url", checkAuth(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		data := BaseData(r, "Edit Post")
		postURL := p.ByName("url")
		post, ok := parsePostForm(r)
		data["Post"] = post
		if !ok {
			data["Error"] = true
			err := globalTemplate.ExecuteTemplate(w, "editPost", data)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			err := UpdatePost(postURL, post)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/posts/"+post.URL, http.StatusFound)
		}
	}))
	router.GET("/posts/:url", checkAuth(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		data := BaseData(r, "Edit Post")
		postURL := p.ByName("url")
		post, err := GetPostByURL(postURL)
		if err != nil {
			data["Error"] = true
			http.Redirect(w, r, "/newpost", http.StatusFound)
		}

		data["Post"] = post
		err = globalTemplate.ExecuteTemplate(w, "editPost", data)
		if err != nil {
			fmt.Println(err)
		}
	}))
	router.GET("/drafts", checkAuth(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		data := BaseData(r, "Drafts")
		posts, err := GetPosts(false)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		fmt.Println(posts)

		data["Posts"] = posts
		err = globalTemplate.ExecuteTemplate(w, "drafts", data)
		if err != nil {
			fmt.Println(err)
		}
	}))
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
