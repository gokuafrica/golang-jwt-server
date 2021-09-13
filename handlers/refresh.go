// provides out of the box handlers for common use-cases
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gokuafrica/golang-jwt-server/constants"
	"github.com/gokuafrica/golang-jwt-server/server"
	"github.com/gokuafrica/golang-jwt-server/utils"
)

// This is the common refresh token handler method. When access token expires and a valid refresh token is present in the cookie, use this handler to expose an endpoint to use the refresh token to get a new access token. The new access token is populated with the previous payload.
func Refresh(rw http.ResponseWriter, r *http.Request) {
	serverConf := server.GetServerConfigs()
	cookie, err := r.Cookie(constants.REFRESH_TOKEN)
	if err != nil {
		http.Error(rw, "REFRESH TOKEN NOT FOUND", http.StatusUnauthorized)
		return
	}
	claimsMap, err := utils.ValidateToken(cookie.Value, serverConf.RefreshToken.Secret)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}
	tokenString, err := utils.GenerateToken(serverConf.Secret, serverConf.Expiry, claimsMap)
	if err != nil {
		http.Error(rw, "Error ocurred while generating jwt", http.StatusInternalServerError)
		return
	}
	token := struct {
		Token string `json:"token"`
	}{tokenString}
	err = json.NewEncoder(rw).Encode(token)
	if err != nil {
		http.Error(rw, "Error ocurred while sending jwt back", http.StatusInternalServerError)
		return
	}
}
