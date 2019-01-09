package main

import "fmt"
import "strconv"

var logs []string

func LogInit() {
	logs = make([]string, 0, 5)
	LogAppend("Log Started...")
}

func LogAppend(content string) {
	logs = append(logs, content)
}

func LogAppendInt(content int) {
	LogAppend(strconv.Itoa(content))
}

func LogDisplay() {
	for index, logContent := range logs {
		fmt.Println(index, ": ", logContent)
	}
}
