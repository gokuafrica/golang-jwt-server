// provides out of the box middlewares for common use-cases
package middleware

import (
	"context"
	"net/http"

	"github.com/gokuafrica/golang-jwt-server/server"
	"github.com/gokuafrica/golang-jwt-server/utils"
)

// just created a key to be used to get save and retrieve the jwt claims
type JWTClaimsKey struct{}

// This is the jwt token middleware. Use this for authenticated requests. In case of valid JWT present, claims are added to the request context.
func JWT(method http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if length := len(tokenString); length > 7 {
			tokenString = tokenString[7:]
		} else {
			tokenString = ""
		}
		claimsMap, err := utils.ValidateToken(tokenString, server.GetServerConfigs().Secret)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), JWTClaimsKey{}, &claimsMap)
		r = r.WithContext(ctx)
		// continue the handler
		method.ServeHTTP(rw, r)
	})
}
