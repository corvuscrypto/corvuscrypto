package main

import (
	"html/template"
	"io/ioutil"
	"log"
)

var templates map[string]*template.Template

var tempFuncs template.FuncMap

func compileTemplates() {

	tempFuncs = make(template.FuncMap)
	tempFuncs["eq"] = func(a, b interface{}) bool {
		return a == b
	}

	templates = make(map[string]*template.Template)
	commonTemplate, err := template.New("").ParseFiles(
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
			templates[name] = templates[name].Funcs(tempFuncs)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}
