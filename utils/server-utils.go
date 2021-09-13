// provides main utility methods used across the server and the project
package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gokuafrica/golang-jwt-server/constants"
	myerrors "github.com/gokuafrica/golang-jwt-server/errors"
	"github.com/golang-jwt/jwt"
)

// error thrown when server.GetNewToken method recieves invalid arguments
var errFooInvalidArguments = fmt.Errorf("invalid arguments")

// extract payload and responsewriter data from the interface array
func GetTokenArguments(arguments []interface{}) (*map[string]string, http.ResponseWriter, error) {
	if length := len(arguments); length < 1 || length > 2 {
		return nil, nil, errFooInvalidArguments
	}
	payload, ok := arguments[0].(*map[string]string)
	if !ok {
		return nil, nil, errFooInvalidArguments
	}
	if length := len(arguments); length == 1 {
		return payload, nil, nil
	}
	rw, ok := arguments[1].(http.ResponseWriter)
	if !ok {
		return nil, nil, errFooInvalidArguments
	}
	return payload, rw, nil
}

// generate new jwt token using given secret, expiry and claims payload
func GenerateToken(secret string, expiry int64, payload *map[string]string) (string, error) {
	claims := jwt.MapClaims{}
	for key, val := range *payload {
		claims[key] = val
	}
	claims["exp"] = time.Now().Add(time.Second * time.Duration(expiry)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// set refresh token cookie
func GenerateRefreshCookie(token string, expiry int64) *http.Cookie {
	return &http.Cookie{
		Name:     constants.REFRESH_TOKEN,
		Value:    token,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Second * time.Duration(expiry)),
	}
}

// Validate given JWT token. Returns claims payload on successful validation, else returns UnauthorizedError
func ValidateToken(tokenstring string, secret string) (*map[string]string, error) {
	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret
		return []byte(secret), nil
	})
	if nil != token && token.Valid {
		// insert claims
		claimsMap := make(map[string]string)
		claims := token.Claims.(jwt.MapClaims)
		for key, value := range claims {
			var val string
			switch v := value.(type) {
			case string:
				val = v
			case float64:
				val = fmt.Sprintf("%f", v)
			case int64:
				val = fmt.Sprintf("%d", v)
			}
			claimsMap[key] = val
		}
		return &claimsMap, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, myerrors.GetUnauthorizedError("MALFORMED TOKEN")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return nil, myerrors.GetUnauthorizedError("EXPIRED TOKEN")
		} else {
			return nil, myerrors.GetUnauthorizedError("INVALID TOKEN")
		}
	} else {
		return nil, myerrors.GetUnauthorizedError("INVALID TOKEN")
	}
}
