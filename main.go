package main

import (
	"github.com/bryanaustin/jack-henry-eval-weather/config"
	"log"
	"net/http"
)

// main is the entry point for the application
func main() {
	// configure service
	config.Init()

	// check for issues
	config.SanityCheck()

	// register handlers
	initHandlers()

	// run the server and capture errors
	err := http.ListenAndServe(config.Current.Server.Listen, nil)

	// print error and exit
	log.Println("http ListenAndServe:", err)
}
