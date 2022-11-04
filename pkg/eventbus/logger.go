package eventbus

import (
	"fmt"
	"log"
)

type LogFunc func(msg string)
type LogError func(err error)

type Logger struct {
	infoFunc    LogFunc
	warningFunc LogFunc
	errorFunc   LogError
}

func printMessage(msg string) {
	log.Print(msg)
}

func printError(err error) {
	log.Print(fmt.Sprint(err))
}

var logger = Logger{
	infoFunc:    printMessage,
	warningFunc: printMessage,
	errorFunc:   printError,
}

func ConfigureLogging(infoLogger LogFunc, warningLogger LogFunc, errorLogger LogError) {
	logger.infoFunc = infoLogger
	logger.warningFunc = warningLogger
	logger.errorFunc = errorLogger
}

func (l Logger) Info(msg string) {
	if l.infoFunc != nil {
		l.infoFunc(msg)
	}
}

func (l Logger) Warning(msg string) {
	if l.warningFunc != nil {
		l.warningFunc(msg)
	}
}

func (l Logger) Error(err error) {
	if l.errorFunc != nil {
		l.errorFunc(err)
	}
}
