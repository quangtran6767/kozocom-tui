package config

import (
	"log"
	"os"
)

var DebugLog *log.Logger

func InitLogger() {
	f, err := os.OpenFile("/tmp/kozocom-tui.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	DebugLog = log.New(f, "[DEBUG]", log.Ltime|log.Lshortfile)
}
