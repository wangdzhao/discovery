package util

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// LogInfo logs info messages
func LogInfo(format string, a ...interface{}) {
	mylog("INFO", format, a...)
}

// LogError logs error messages
func LogError(format string, a ...interface{}) {
	mylog("ERROR", format, a...)
}

// log is a helper function for logging
func mylog(level, format string, a ...interface{}) {
	logDir := "data"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.Mkdir(logDir, 0755)
	}
	logFile := filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}
	defer f.Close()

	msg := fmt.Sprintf("%s [%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), level, fmt.Sprintf(format, a...))
	_, _ = f.WriteString(msg)
	fmt.Print(msg)
}
