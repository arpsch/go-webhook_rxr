package logger

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	PANIC Level = iota
	ERROR
	WARN
	INFO
	DEBUG
)

type Logger struct {
}

// Log support leveled and structured logging
func (l *Logger) Log(level Level, msg string, val ...interface{}) {
	fmt.Printf(level.String()+": "+msg, val...)
	fmt.Printf("\n")
}

type Level uint32

// String stringify the constant
func (l Level) String() string {
	switch l {
	case PANIC:
		return "Panic"
	case ERROR:
		return "Error"
	case WARN:
		return "Warn"
	case INFO:
		return "Info"
	case DEBUG:
		return "Debug"
	}
	return "unknown"
}

// LogRequest simple logger middleware to log requests
func LogRequest(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("start: %v - %s\n", time.Now(), r.URL.Path)
		f(w, r)
	}
}
