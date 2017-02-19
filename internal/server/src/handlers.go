package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/julienschmidt/httprouter"
)

type compressionWriter struct {
	io.Writer
	http.ResponseWriter
	overrideCompression bool
}

type nopWriter struct {
	io.Writer
}

func (n *nopWriter) Write(b []byte) (int, error) {
	return 0, nil
}

func (c *compressionWriter) Write(p []byte) (int, error) {
	if c.Header().Get("Content-Type") == "" {
		c.Header().Set("Content-Type", http.DetectContentType(p))
	}
	var b int
	var err error
	if c.overrideCompression {
		b, err = c.ResponseWriter.Write(p)
	} else {
		c.Header().Set("Content-Encoding", "gzip")
		b, err = c.Writer.Write(p)
	}
	if err != nil {
		fmt.Println(err)
	}
	return b, err
}

func (c *compressionWriter) WriteHeader(status int) {
	if status != http.StatusFound || status != http.StatusOK || status != http.StatusNotModified {
		c.overrideCompression = true
	}
	c.ResponseWriter.WriteHeader(status)
}

func cached(r *http.Request) bool {
	if r.Header.Get("If-None-Match") == Config.Etag {
		return true
	}
	return false
}

func compress(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		newWriter := new(compressionWriter)
		newWriter.ResponseWriter = w
		gzw := gzip.NewWriter(w)
		newWriter.Writer = gzw
		handler(newWriter, r, p)
		if newWriter.overrideCompression {
			gzw.Reset(new(nopWriter))
		}
		err := gzw.Close()
		if err != nil && err != http.ErrBodyNotAllowed {
			fmt.Println(err)
		}
	}
}

func serveStatic(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	//do a cache check
	if cached(r) {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	filepath := p.ByName("filepath")

	var contentType = "text/plain"
	switch filepath[len(filepath)-3:] {
	case "css":
		contentType = "text/css"
		break
	case ".js":
		contentType = "text/javascript"
		break
	}

	w.Header().Set("Content-Type", contentType)

	// first reject all requests with .. in the path
	if strings.Contains(filepath, "..") || len(filepath) < 3 {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	file, err := os.Open(Config.StaticPath + filepath)
	defer file.Close()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//if we get this far set the cache header values
	w.Header().Set("Cache-Control", "max-age=31104000")
	w.Header().Set("ETag", Config.Etag)

	reader := bufio.NewReader(file)
	_, err = reader.WriteTo(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
	}
}

func serverErrorPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
	data := BaseData(r, "CorvusCrypto.com - internal server error")
	err := globalTemplate.ExecuteTemplate(w, "500", data)
	if err != nil {
		fmt.Println(err)
	}
}

func panicHandler(w http.ResponseWriter, r *http.Request, e interface{}) {
	LogError(e)
	serverErrorPage(w, r)
}

func notFoundPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	data := BaseData(r, "CorvusCrypto.com - 404")
	err := globalTemplate.ExecuteTemplate(w, "404", data)
	if err != nil {
		fmt.Println(err)
	}
}

//NotFoundHandler deals with the 404 pages
type NotFoundHandler struct{}

func (n NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serverErrorPage(w, r)
}

func getRouter() *httprouter.Router {
	router := httprouter.New()
	router.NotFound = new(NotFoundHandler)
	router.PanicHandler = panicHandler
	router.GET("/static/*filepath", compress(serveStatic))
	router.GET("/favicon.ico", compress(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		http.Redirect(w, r, "static/images/favicon.ico", http.StatusFound)
	}))
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
		var err error

		last := bson.ObjectId(r.FormValue("last"))
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
			notFoundPage(w, r)
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
