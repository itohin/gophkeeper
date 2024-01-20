package main

import (
	"github.com/itohin/gophkeeper/cmd"
	"github.com/sirupsen/logrus"
	"log"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

func NewLogger() Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	return logger
}

type User struct {
	nic   string
	login string
	token string
}

func main() {
	l := NewLogger()
	l.Info("hello from main")
	err := cmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
