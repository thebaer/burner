package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/thebaer/burner/validate"
)

// Serve starts an HTTP server that handles auth requests from nginx.
func Serve(port int) error {
	if port <= 0 {
		return errors.New("auth server: Invalid port number.")
	}
	serverPort = port

	mailInfo := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	mailInfo.Printf("Starting mail auth server on :%d", serverPort)

	http.HandleFunc("/auth", authHandler)
	http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", serverPort), nil)

	return nil
}

var (
	// Port that the auth server will run on.
	serverPort int

	// Regular expression for matching / finding a valid To address.
	smtpEmailReg = regexp.MustCompile("<(.+)>")
)

// authHandler works with nginx to determine whether or not a receipient email
// address is valid. If it is, running mail server's information is passed
// back.
func authHandler(w http.ResponseWriter, r *http.Request) {
	toHeader := r.Header.Get("Auth-SMTP-To")
	if toHeader == "" {
		w.Header().Set("Auth-Status", "Unrecognized receipient.")
		w.Header().Set("Auth-Error-Code", "550")
		return
	}

	to := smtpEmailReg.FindStringSubmatch(toHeader)[1]
	if to == "" {
		w.Header().Set("Auth-Status", "Unrecognized receipient.")
		w.Header().Set("Auth-Error-Code", "550")
		return
	}
	if err := validate.Email(to); err != nil {
		// Email address validation failed
		w.Header().Set("Auth-Status", err.Error())
		w.Header().Set("Auth-Error-Code", "550")
		return
	}

	// Email passed validation, send back mail server information
	w.Header().Set("Auth-Status", "OK")
	w.Header().Set("Auth-Server", "127.0.0.1")
	w.Header().Set("Auth-Port", fmt.Sprintf("%d", serverPort))
}
