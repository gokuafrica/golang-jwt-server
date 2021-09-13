// provides methods to retrieve values from the request context
package context

import (
	"log"
	"net/http"

	"github.com/gokuafrica/golang-jwt-server/middleware"
	"github.com/gokuafrica/golang-jwt-server/server"
)

// Get JWT claims/payload from the request. Authorization header is processed and the claims are inserted in the JWT middleware.
func GetClaims(r *http.Request) *map[string]string {
	return r.Context().Value(middleware.JWTClaimsKey{}).(*map[string]string)
}

// Get server object from the request context. It is inserted artificially in the logger middleware.
func GetServer(r *http.Request) *server.Server {
	return r.Context().Value(server.ServerKey{}).(*server.Server)
}

// Get logger object from the request context based on initial server configuration. It is inserted artificially in the logger middleware.
func GetLogger(r *http.Request) *log.Logger {
	return r.Context().Value(server.LogKey{}).(*log.Logger)
}
