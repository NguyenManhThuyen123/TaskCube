package core

import (
	"app/config"
	"fmt"
	"os"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func WriteLog(message string) {
	time.LoadLocation(config.Config("APP_TIME_ZONE"))
	dirPath := "./assets/log"
	fileName := fmt.Sprintf("%s/%s.txt", dirPath, time.Now().Format("2006-01-02"))

	logTime := time.Now().Format("2006-01-02 15:04:05")
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)
	defer f.Close()

	fmt.Println(123123)

	msg := fmt.Sprintf("[%s] | %s \n", logTime, message)
	f.WriteString(msg)

}
