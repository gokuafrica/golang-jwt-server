// provides the singleton server object and its respective methods that are used to handle server management
package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gokuafrica/golang-jwt-server/config"
	"github.com/gokuafrica/golang-jwt-server/utils"
	"github.com/gorilla/mux"
)

// has the main server components like mux router, configurations, and http server
type Server struct {
	sm     *mux.Router
	conf   *config.Server
	server *http.Server
}

// singleton server object that is present for the lifetime of the application
var serverInstance *Server

// used to initialize the server singleton object by taking server configs as input
func Initialize(config *config.Server) *Server {
	// log request middleware
	// the reason we are using a struct object rather than directly calling the middleware is because we want a place to preserve the log object
	lr := newLogRequestMiddleware(config.Log)

	// create new servermux to server our endpoints
	sm := mux.NewRouter()

	// as log request is a global middleware, rather than attaching it to any sub router, we attach it to the main mux router itself
	sm.Use(lr.logRequest)

	// create server with appropriate settings
	httpServer := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	server := &Server{
		sm:     sm,
		conf:   config,
		server: httpServer,
	}
	serverInstance = server
	return server
}

// get server configurations
func GetServerConfigs() config.Server {
	return *serverInstance.conf
}

// Start the http server using this method. Returns a channel to block main thread from exiting as long as server is running.
func (s *Server) Start() chan os.Signal {
	// run server as a separate go routine so that it doesn't block main
	go func() {
		err := s.server.ListenAndServe()
		if err != nil {
			s.conf.Log.Fatal(err)
		}
	}()
	// create channel to keep track of server kill/interrupt attempt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	return sigChan
}

// To be called when server has recieved stop/interrupt signal. Recommended to be used just after trying to recieve values with the channel.
func (s *Server) HandleServerStop() {
	s.conf.Log.Println("Attempting shutdown")

	// give 30 seconds to handle existing requests and then shutdown (gracefully)
	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s.server.Shutdown(tc)
}

// get primary mux router
func (s *Server) GetMuxRouter() *mux.Router {
	return s.sm
}

// Generate and set the refresh token in the cookie. Returns error if refresh token configurations are empty.
func (s *Server) setRefreshToken(rw http.ResponseWriter, payload *map[string]string) error {
	refreshToken := s.conf.RefreshToken
	if refreshToken == nil {
		return fmt.Errorf("refresh token configurations are empty")
	}
	token, err := utils.GenerateToken(refreshToken.Secret, refreshToken.Expiry, payload)
	if err != nil {
		return err
	}
	cookie := utils.GenerateRefreshCookie(token, refreshToken.Expiry)
	http.SetCookie(rw, cookie)
	return nil
}

// Used to generate new access tokens and set refresh token in the cookie with given payload. Recommended to be called in the login handler post successful login to generate respective tokens for the user. Method takes in map[string]string (jwt payload) and http.ResponseWriter object (to set refresh token cookie) as parameters. Passing claims payload is mandatory. Passing http.ResponseWriter is optional and passing it signals that refresh token has to be generated along with the access token. In case refresh token configs are not set in the server, an error is returned instead.
func (s *Server) GetNewToken(arguments ...interface{}) (string, error) {
	payload, rw, err := utils.GetTokenArguments(arguments)
	if err != nil {
		return "", err
	}
	token, err := utils.GenerateToken(s.conf.Secret, s.conf.Expiry, payload)
	if err != nil {
		return "", err
	}
	if rw != nil {
		err = s.setRefreshToken(rw, payload)
		if err != nil {
			return "", err
		}
	}
	return token, nil
}
