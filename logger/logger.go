package logger

import (
	"log"
	"net/http"
	"time"
)

// Log simple logger middleware to log requests
func Log(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("start: %v - %s\n", time.Now(), r.URL.Path)
		f(w, r)
		log.Printf("end: %v - %s\n", time.Now(), r.URL.Path)
	}
}
