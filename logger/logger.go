package logger

import (
	"fmt"
	"log"
	"os"
	"io"
)

var LogFile *os.File
var Logchan chan string

func init() {
	// logging
	var err error
	LogFile, err = os.OpenFile("logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error while opening log file!")
		panic(err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, LogFile))
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func LogHandler() {
	defer LogFile.Close()
	for {
		select {
		case entry := <-Logchan:
			log.Println(entry)
		}
	}
}
