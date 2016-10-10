package main

import (
	"encoding/json"
	"log"
	"os"
)

const blogServerConfig = "./blogServerCfg.json"

type config struct {
	TemplatePath string
	StaticPath   string
	Assets       map[string]interface{}
	LogFilePath  string
}

//Config for ALL the things
var Config *config

func loadConfig() {
	Config = new(config)

	configFile, err := os.OpenFile(blogServerConfig, os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}

	cfgDecoder := json.NewDecoder(configFile)
	err = cfgDecoder.Decode(Config)
	if err != nil {
		//create defaults
		Config.StaticPath = "../../../static"
		Config.TemplatePath = "../../templates"
		Config.Assets = map[string]interface{}{
			"CSS": map[string]string{
				"URL":     "/static/css/main.css",
				"Version": "1",
			},
		}
		Config.LogFilePath = "../logs/blog.log"
		cfgEncoder := json.NewEncoder(configFile)
		cfgEncoder.Encode(Config)
	}
}
