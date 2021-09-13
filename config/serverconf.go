// stores the config data types for the project
package config

import "log"

// refresh token is saved in the cookie and sent on every request
type RefreshToken struct {
	// refresh token secret. Don't put the same as access token secret
	Secret string

	// set refresh token expiry time
	Expiry int64

	// if set to true, would set secure flag in cookie to true
	Https bool
}

// server config object stores all the relevant server details
type Server struct {
	Port   string
	Log    *log.Logger
	Secret string
	Expiry int64

	// This is an optional field. Not passing refresh token configs would log out the user on expiration of the access token rather than generate a new access token from the refresh token
	RefreshToken *RefreshToken
}
