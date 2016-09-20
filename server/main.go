package main

import (
	"bufio"
	"bytes"
	"log"
	"net/http"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

var indexContent []byte

func loadIndexIntoMem() {
	// This is pretty small website so we just load the entire index Response
	// into memory for the fastest serving possible

	buffer := new(bytes.Buffer)
	bufWriter := bufio.NewWriter(buffer)

	//now compile the template and save the bytes to memory
	htmlTemplate, err := template.ParseFiles("../static/index.html")
	if err != nil {
		log.Fatal(err)
	}

	_, err = htmlTemplate.ParseFiles("../static/css/main.css",
		"../static/js/main.js")
	if err != nil {
		log.Fatal(err)
	}
	htmlTemplate.Execute(bufWriter, nil)
	bufWriter.Flush()
	indexContent = buffer.Bytes()

}
func main() {

	loadIndexIntoMem()

	router := httprouter.New()

	//add the static file routes here
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write(indexContent)
	})
	http.ListenAndServe(":8080", router)
}
