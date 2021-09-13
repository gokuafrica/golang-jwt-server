package server

import (
	"context"
	"log"
	"net/http"
)

// used to store log object
type logRequestMiddleWare struct {
	l *log.Logger
}

type (
	// key used to save and retrieve the logger object
	LogKey struct{}
	// key used to save and retrieve the server object
	ServerKey struct{}
)

// constructor
func newLogRequestMiddleware(l *log.Logger) *logRequestMiddleWare {
	return &logRequestMiddleWare{l: l}
}

// Log request middleware logs all the incoming request URIs along with the respective HTTP method. Automatically attached to the main mux router. It injects logger and server object to the request context.
func (l *logRequestMiddleWare) logRequest(method http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		l.l.Println("incoming request " + r.Method + " " + r.RequestURI)
		ctx := context.WithValue(r.Context(), LogKey{}, l.l)
		r = r.WithContext(ctx)
		ctx = context.WithValue(r.Context(), ServerKey{}, serverInstance)
		r = r.WithContext(ctx)
		method.ServeHTTP(rw, r)
	})
}
