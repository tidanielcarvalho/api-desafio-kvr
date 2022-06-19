package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Log struct{}

func isId(id string) string {
	var isId = "[" + id + "] "
	if id == "" {
		isId = " "
	}
	return isId
}

func (l *Log) Debug(id string, str string) {
	forTestingErr := godotenv.Load("../.env")
	if forTestingErr != nil {
		// retry
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}
	}

	if os.Getenv("LOG_DEBUG") == "enable" {
		fmt.Println("[" + time.Now().Format("2006-01-02 15:04:05.000") + "][DEBUG]" + isId(id) + str)
	}
}

func (l *Log) Info(id string, str string) {
	fmt.Println("[" + time.Now().Format("2006-01-02 15:04:05.000") + "][INFO]" + isId(id) + str)
}

func (l *Log) Warn(id string, str string) {
	fmt.Println("[" + time.Now().Format("2006-01-02 15:04:05.000") + "][WARN]" + isId(id) + str)
}

func (l *Log) Error(id string, str string) {
	fmt.Println("[" + time.Now().Format("2006-01-02 15:04:05.000") + "][ERROR]" + isId(id) + str)
}

func (l *Log) Fatal(id string, str string, err error) {
	fmt.Println("[" + time.Now().Format("2006-01-02 15:04:05.000") + "][FATAL]" + isId(id) + str)
	panic(err)
}
