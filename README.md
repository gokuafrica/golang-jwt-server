# golang-jwt-server
#### A simple http server with jwt handling capabilities


## Features

- JWT token generation with expiry
- Out of the box middleware to handle authenticaetd requests and process claims payload
- Refresh token generation and management done through cookies with separate secret and expiry configurations

## Installation

This package requires golang to run.

```go
go get github.com/gokuafrica/golang-jwt-server
```


## Dependencies

**golang-jwt-server** currently depends on the following:

| Module |
| ------ |
| [Gorilla Mux](https://github.com/gorilla/mux) |
| [Golang-JWT](https://github.com/golang-jwt/jwt) |

## Development

Sample main program that uses **golang-jwt-server**

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gokuafrica/golang-jwt-server/config"
	"github.com/gokuafrica/golang-jwt-server/context"
	"github.com/gokuafrica/golang-jwt-server/handlers"
	"github.com/gokuafrica/golang-jwt-server/middleware"
	"github.com/gokuafrica/golang-jwt-server/server"
)

func main() {
	// login handler
	login := func(rw http.ResponseWriter, r *http.Request) {
		/*
			Login validation logic. 
            For eg: take email and password and verify in database.
		*/
		log := context.GetLogger(r)
		server := context.GetServer(r)
		mymap := make(map[string]string)
		mymap["userid"] = "gokuafrica"
		token, err := server.GetNewToken(&mymap, rw)
		if err != nil {
			log.Panic(err)
			return
		}
		fmt.Fprintln(rw, token)
	}

	// an authenticated request
	authenticated := func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "Just an authenticated request")
	}

	// get new jwt from refresh token handler
	refresh := handlers.Refresh

	// set server configs
	config := config.Server{
		Port:   "9090",
		Log:    log.New(os.Stdout, "hello-world", log.LstdFlags),
		Secret: "My Secret",
		Expiry: 600,
		// Refresh token configs are optional. Only use if refresh token to be used.
		RefreshToken: &config.RefreshToken{
			Secret: "My Secret1",
			Expiry: 86400,
			Https:  false,
		},
	}

	// initialise server with given configs
	jwtServer := server.Initialize(&config)

	// get mux router from server
	sm := jwtServer.GetMuxRouter()

	// authenticated router
	authRouter := sm.Path("/").Subrouter()
	// uses the jwt middleware for all authentication
	authRouter.Use(middleware.JWT)
	authRouter.HandleFunc("/", authenticated)

	// add unauthenticated endpoints directly into the main mux router
	sm.HandleFunc("/login", login)
	sm.HandleFunc("/refresh", refresh)

	// start the jwt server with above configurations and endpoints
	sigChan := jwtServer.Start()
	// block main thread using server channel
	<-sigChan
	// handle graceful server shutdown after channel recieves data
	jwtServer.HandleServerStop()
}
```

## License

MIT