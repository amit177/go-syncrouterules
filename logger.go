package main

import (
	"fmt"
	"os"
	"time"
)

type logLevel int

// The log levels that exist
const (
	INFO  logLevel = 0
	WARN  logLevel = 1
	DEBUG logLevel = 2
	ERROR logLevel = 3
	FATAL logLevel = 4
)

const debugMode = false

// LogMessage formats the inputted message and beauty prints it
func LogMessage(level logLevel, location string, message string) {
	if level == DEBUG && !debugMode {
		return
	}
	dt := time.Now()

	if (level == INFO || level == WARN) && !debugMode {
		fmt.Printf("\033[1;35m[%s] %s[%s]\033[0m %s\n", dt.Format("01-02-2006 15:04"), getLevelColor(level), level, message)
	} else {
		fmt.Printf("\033[1;35m[%s] %s[%s] \033[1;37m[%s]\033[0m %s\n", dt.Format("01-02-2006 15:04"), getLevelColor(level), level, location, message)
	}

	if level == FATAL {
		os.Exit(1)
	}
}

// the function returns the color according to the requested level
func getLevelColor(level logLevel) string {
	str := ""
	switch level {
	case INFO:
		str = "\033[1;96m" //light cyan
	case WARN:
		str = "\033[1;33m" //yellow
	case DEBUG:
		str = "\033[1;90m" //dark gray
	case ERROR:
		str = "\033[1;91m" //light red
	case FATAL:
		str = "\033[1;31m" //red
	}

	return str
}

// a string representation of the level
func (level logLevel) String() string {
	switch level {
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case DEBUG:
		return "DEBUG"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}
