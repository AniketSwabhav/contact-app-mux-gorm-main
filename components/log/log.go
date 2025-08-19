package log

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Printf(format string, args ...interface{})
	Print(value ...interface{})
	Error(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Fatalf(format string, args ...interface{})
}

var logger = logrus.New()

type Log struct{}

func NewLog() *Log {
	return &Log{}
}

func GetLogger() Logger {
	return logger
}

func (l *Log) Print(value ...interface{}) {
	fmt.Println(value...)
}

func (l *Log) Error(args ...interface{}) {
	fmt.Println(args...)
}
