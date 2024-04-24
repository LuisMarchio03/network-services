package logger

import (
	"log"
	"os"

	"github.com/fatih/color"
)

var (
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
)

// InitLogger inicializa o logger
func InitLogger() {
	infoLogger = log.New(os.Stdout, color.BlueString("[INFO] "), log.Ldate|log.Ltime|log.Lmsgprefix)
	warningLogger = log.New(os.Stdout, color.YellowString("[WARNING] "), log.Ldate|log.Ltime|log.Lmsgprefix)
	errorLogger = log.New(os.Stderr, color.RedString("[ERROR] "), log.Ldate|log.Ltime|log.Lmsgprefix)
}

// Info loga uma mensagem de informação
func Info(message string) {
	infoLogger.Println(message)
}

// Warning loga uma mensagem de aviso
func Warning(message string) {
	warningLogger.Println(message)
}

// Error loga uma mensagem de erro
func Error(message string, err error) {
	if err != nil {
		errorLogger.Printf("%s: %v", message, err)
	} else {
		errorLogger.Println(message)
	}
}
