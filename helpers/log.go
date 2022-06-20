package helpers

import (
	"fmt"
	"os"
	"time"
)

type Log struct {
	colorReset  string
	colorRed    string
	colorGreen  string
	colorYellow string
	colorBlue   string
	colorPurple string
	colorCyan   string
	colorWhite  string
}

func StartColors(l *Log) {
	l.colorReset = string("\033[0m")
	l.colorRed = string("\033[31m")
	l.colorGreen = string("\033[32m")
	l.colorYellow = string("\033[33m")
	l.colorBlue = string("\033[34m")
	l.colorPurple = string("\033[35m")
	l.colorCyan = string("\033[36m")
	l.colorWhite = string("\033[37m")
}

func isId(id string) string {
	var isId = "[" + id + "] "
	if id == "" {
		isId = " "
	}
	return isId
}

func (l *Log) Debug(id string, str string) {
	if os.Getenv("LOG_DEBUG") == "enable" {
		fmt.Println("[" + time.Now().Format("2006-01-02 15:04:05.000") + "][DEBUG]" + isId(id) + str)
	}
}

func (l *Log) Info(id string, str string) {
	fmt.Println("[" + time.Now().Format("2006-01-02 15:04:05.000") + "][INFO]" + isId(id) + str)
}

func (l *Log) Warn(id string, str string) {
	StartColors(l)
	fmt.Println(l.colorYellow + "[" + time.Now().Format("2006-01-02 15:04:05.000") + "][WARN]" + isId(id) + str + l.colorReset)
}

func (l *Log) Error(id string, str string) {
	StartColors(l)
	fmt.Println(l.colorRed + "[" + time.Now().Format("2006-01-02 15:04:05.000") + "][ERROR]" + isId(id) + str + l.colorReset)
}

func (l *Log) Fatal(id string, str string, err error) {
	StartColors(l)
	fmt.Println(l.colorRed + "[" + time.Now().Format("2006-01-02 15:04:05.000") + "][FATAL]" + isId(id) + str + l.colorReset)
	panic(err)
}
