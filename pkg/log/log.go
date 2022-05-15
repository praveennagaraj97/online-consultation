package logger

import (
	"log"
	"os"
)

// Pritty print error log with red color and stops the application.
func ErrorLogFatal(msg interface{}) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Fatal(msg)
}

// Pritty print success log with green color
func PrintLog(msg string) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	logger.Println(msg)
}
