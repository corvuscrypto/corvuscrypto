package main

import (
	"encoding/json"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

const cfgFileName = "./blogCMSCfg.json"

type config struct {
	TemplatePath  string
	StaticPath    string
	Assets        map[string]interface{}
	LogFilePath   string
	Etag          string
	OwnerUsername string
	OwnerPassword string
	passHash      []byte
}

//Config for ALL the things
var Config *config

func loadConfig() {
	Config = new(config)

	configFile, err := os.OpenFile(cfgFileName, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

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
		Config.LogFilePath = "../logs/cms.log"
		cfgEncoder := json.NewEncoder(configFile)
		cfgEncoder.Encode(Config)
	}
	if Config.OwnerPassword == "" || Config.OwnerUsername == "" {
		log.Fatal("Both the username and password must be set!")
	}
	Config.passHash, err = bcrypt.GenerateFromPassword([]byte(Config.OwnerPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
}
