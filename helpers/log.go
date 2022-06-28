package helpers

import (
	"fmt"
	"os"
	"time"
)

type Log struct{}

func (l Log) colorReset() string {
	return string("\033[0m")
}
func (l Log) colorRed() string {
	return string("\033[31m")
}
func (l Log) colorGreen() string {
	return string("\033[32m")
}
func (l Log) colorYellow() string {
	return string("\033[33m")
}
func (l Log) colorBlue() string {
	return string("\033[34m")
}
func (l Log) colorPurple() string {
	return string("\033[35m")
}
func (l Log) colorCyan() string {
	return string("\033[36m")
}
func (l Log) colorWhite() string {
	return string("\033[37m")
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
	fmt.Println(l.colorYellow() + "[" + time.Now().Format("2006-01-02 15:04:05.000") + "][WARN]" + isId(id) + str + l.colorReset())
}

func (l *Log) Error(id string, str string) {
	fmt.Println(l.colorRed() + "[" + time.Now().Format("2006-01-02 15:04:05.000") + "][ERROR]" + isId(id) + str + l.colorReset())
}

func (l *Log) Fatal(id string, str string, err error) {
	fmt.Println(l.colorRed() + "[" + time.Now().Format("2006-01-02 15:04:05.000") + "][FATAL]" + isId(id) + str + l.colorReset())
	panic(err)
}
