package main

import (
	"flag"
	"log"

	"github.com/thebaer/burner/auth"
	"github.com/thebaer/burner/database"
	"github.com/thebaer/burner/mail"
)

var host = flag.String("h", "example.com", "Domain this service lives on.")

func main() {
	// Parse configuration flags and validate
	flag.Parse()
	if *host == "example.com" {
		log.Printf("WARNING: Default hostname (example.com) unchanged. Use -h flag to set correct host.")
	}

	// Connect to database
	db, err := database.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// TODO: make port numbers configurable
	go func() {
		auth.Serve(8080)
	}()

	mail.Serve(*host, 2525)
}
