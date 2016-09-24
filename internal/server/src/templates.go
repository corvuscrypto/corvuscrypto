package main

import (
	"io/ioutil"
	"log"
	"text/template"
)

var globalTemplate *template.Template

var tempFuncs template.FuncMap

func walkAndCompile(subdir string) {
	files, err := ioutil.ReadDir(Config.TemplatePath + subdir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() {
			walkAndCompile(subdir + "/" + file.Name())
		} else {
			var err error
			name := file.Name()
			globalTemplate, err = globalTemplate.ParseFiles(Config.TemplatePath + subdir + "/" + name)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func compileTemplates() {

	tempFuncs = make(template.FuncMap)
	tempFuncs["eq"] = func(a, b interface{}) bool {
		return a == b
	}

	globalTemplate = template.New("globalCommon")

	walkAndCompile("")
	globalTemplate = globalTemplate.Funcs(tempFuncs)

}
