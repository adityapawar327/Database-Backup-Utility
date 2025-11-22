package utils

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func InitLogger() {
	file, err := os.OpenFile("backup_tool.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func LogInfo(message string) {
	if InfoLogger != nil {
		InfoLogger.Println(message)
	}
	log.Println("INFO:", message)
}

func LogError(message string) {
	if ErrorLogger != nil {
		ErrorLogger.Println(message)
	}
	log.Println("ERROR:", message)
}
