package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler, logger Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}
		next.ServeHTTP(recorder, r)

		msg := fmt.Sprintf("%s [%s] %s %s %s %d %d \"%s\"",
			r.RemoteAddr,
			now.UTC().Format(time.RFC3339),
			r.Method,
			r.URL.Path+r.URL.RawQuery,
			r.Proto,
			recorder.Status,
			time.Since(now).Milliseconds(),
			r.Header.Get("user-agent"),
		)
		logger.Info(msg)
	})
}

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}
