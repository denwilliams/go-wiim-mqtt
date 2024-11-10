package logging

import (
	"io"
	"log"
	"os"
)

var (
	debugLogger   *log.Logger
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
)

const DefaultFlags = log.Ldate | log.Ltime

func Init(out io.Writer, flag int) {
	so := out
	se := out

	if so == nil {
		so = os.Stdout
	}
	if se == nil {
		se = os.Stderr
	}

	debugLogger = log.New(so, "DEBG: ", flag)
	infoLogger = log.New(so, "INFO: ", flag)
	warningLogger = log.New(so, "WARN: ", flag)
	errorLogger = log.New(se, "ERRO: ", flag)
}

func Debug(format string, v ...interface{}) {
	// debugLogger.Printf(format+"\n", v...)
}

func Info(format string, v ...interface{}) {
	infoLogger.Printf(format+"\n", v...)
}

func Warn(format string, v ...interface{}) {
	warningLogger.Printf(format+"\n", v...)
}

func Error(format string, v ...interface{}) {
	errorLogger.Printf(format+"\n", v...)
}
