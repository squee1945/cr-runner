package main

import (
	"encoding/json"
	"fmt"
)

func logInfo(template string, args ...any) {
	logStructured("INFO", template, args...)
}

func logWarn(template string, args ...any) {
	logStructured("WARN", template, args...)
}

func logError(template string, args ...any) {
	logStructured("ERROR", template, args...)
}

type structuredLog struct {
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

func logStructured(severity string, template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	sl := structuredLog{
		Severity: severity,
		Message:  msg,
	}
	content, err := json.Marshal(sl)
	if err != nil {
		fmt.Printf("Failed to log (message below): %v\n", err)
		fmt.Println(msg)
		return
	}
	fmt.Println(string(content))
}
