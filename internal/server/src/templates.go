package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"text/template"
	"time"
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

	tempFuncs["date"] = func(a interface{}) interface{} {
		date := a.(time.Time)
		return fmt.Sprintf("%02d/%02d/%d", date.Month(), date.Day(), date.Year())
	}
	tempFuncs["title"] = strings.Title
	tempFuncs["has"] = strings.Contains
	tempFuncs["join"] = strings.Join

	globalTemplate = template.New("globalCommon").Funcs(tempFuncs)

	walkAndCompile("")

}
