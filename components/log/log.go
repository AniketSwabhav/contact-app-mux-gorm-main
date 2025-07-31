package log

import "fmt"

type Logger interface {
	Print(value ...interface{})
}

type Log struct{}

func NewLog() *Log {
	return &Log{}
}

func (l *Log) Print(value ...interface{}) {
	fmt.Println(value)
}
