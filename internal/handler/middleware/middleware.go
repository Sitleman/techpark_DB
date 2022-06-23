package mw

import (
	"fmt"
	"net/http"
	"time"
)

func TimeLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start).Microseconds()
		if duration > 15000 {
			fmt.Println(duration, " ", r.Method, " ", r.URL)
		}
	})
}
