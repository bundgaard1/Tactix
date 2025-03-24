package engine

import (
	"os"
	"time"
)

var Logging struct {
	initialized bool
	file        *os.File
}

func InitLogging() {
	// filename := time.Now().Format("2006-01-02-15-04-05")
	filename := "haha"
	if _, err := os.Stat("log"); os.IsNotExist(err) {
		os.Mkdir("log", 0777)
	}
	file, err := os.Create("log/" + filename + ".log")
	if err != nil {
		panic(err)
	}
	Logging.file = file
	Logging.initialized = true

}

func Log(message string) {
	if !Logging.initialized {
		return
	}
	now := time.Now().Format("15:04:05.000")
	Logging.file.WriteString("[" + now + "] " + message + "\n")
}
