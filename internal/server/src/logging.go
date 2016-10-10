package main

import (
	"log"
	"os"
)

//BlogLogger is the logger to be used for logging all errors and such. Can get away with not using a syslog spec'd
//logging system since this is just a personal blog.
var BlogLogger *log.Logger

func createLogger() {
	logFile, err := os.OpenFile(Config.LogFilePath, os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Fatal(err)
	}

	BlogLogger = log.New(logFile, "[Blog Server]", log.Lshortfile|log.Ltime|log.Ldate)
}

//LogError prepends an [ERROR] prefix so that we know it was an error, nothing we are debugging
func LogError(v ...interface{}) {
	BlogLogger.Println("[ERROR]", v)
}

//LogDebug prepends an [DEBUG] prefix so that we know it was debugging
func LogDebug(v ...interface{}) {
	BlogLogger.Println("[DEBUG]", v)
}
