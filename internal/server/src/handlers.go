package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type compressionWriter struct {
	io.Writer
	http.ResponseWriter
}

func (c *compressionWriter) Write(p []byte) (int, error) {
	if c.Header().Get("Content-Type") == "" {
		c.Header().Set("Content-Type", http.DetectContentType(p))
	}
	b, err := c.Writer.Write(p)
	if err != nil {
		fmt.Println(err)
	}
	return b, err
}

func compress(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Encoding", "gzip")
		newWriter := new(compressionWriter)
		newWriter.ResponseWriter = w
		gzw := gzip.NewWriter(w)
		newWriter.Writer = gzw
		handler(newWriter, r, p)
		err := gzw.Close()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func serveStatic(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	filepath := p.ByName("filepath")
	// first reject all requests with .. in the path
	if strings.Contains(filepath, "..") || len(filepath) < 3 {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	var contentType = ""
	switch filepath[len(filepath)-3:] {
	case "css":
		contentType = "text/css"
		break
	case ".js":
		contentType = "text/javascript"
		break
	}

	w.Header().Set("Content-Type", contentType)

	file, err := os.Open(Config.StaticPath + filepath)
	defer file.Close()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	reader := bufio.NewReader(file)
	_, err = reader.WriteTo(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
	}
}

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
	router.GET("/static/*filepath", compress(serveStatic))
	router.GET("/", compress(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	}))
	router.GET("/posts", compress(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	}))
	router.GET("/posts/:postURL", compress(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
	}))
	router.GET("/about", compress(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		data := BaseData(r, "CorvusCrypto.com - About Me")
		err := globalTemplate.ExecuteTemplate(w, "about", data)
		if err != nil {
			fmt.Println(err)
		}
	}))
	return router
}
