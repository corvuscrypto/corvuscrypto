package main

import (
	"io/ioutil"
	"log"
	"text/template"
)

var templates map[string]*template.Template

func compileTemplates() {
	templates = make(map[string]*template.Template)
	commonTemplate, err := new(template.Template).ParseFiles(
		Config.TemplatePath+"/header.html",
		Config.TemplatePath+"/navbar.html",
		Config.TemplatePath+"/footer.html",
	)
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir(Config.TemplatePath + "/pages")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			var err error
			name := file.Name()
			templates[name], err = commonTemplate.ParseFiles(Config.TemplatePath + "/pages/" + name)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}
