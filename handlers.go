package main

import (
	"github.com/bryanaustin/jack-henry-eval-weather/handlers/weather"
	"net/http"
)

// initHandlers will list all of the handlers.
// Please keep the list sorted by path.
func initHandlers() {
	http.HandleFunc("GET /weather", weather.Handle)
}
