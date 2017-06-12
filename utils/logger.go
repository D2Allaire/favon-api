package utils

import (
	"io"
	"log"
)

var InfoLog *log.Logger
var ErrorLog *log.Logger

func InitLog(infoHandle io.Writer, errorHandle io.Writer) {

	InfoLog = log.New(infoHandle,
		"LOG: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	ErrorLog = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
