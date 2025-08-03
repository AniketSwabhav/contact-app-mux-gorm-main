package log

import "fmt"

type Logger interface {
	Print(value ...interface{})
	Error(args ...interface{})
}

type Log struct{}

func NewLog() *Log {
	return &Log{}
}

func (l *Log) Print(value ...interface{}) {
	fmt.Println(value...)
}

func (l *Log) Error(args ...interface{}) {
	fmt.Println(args...)
}
