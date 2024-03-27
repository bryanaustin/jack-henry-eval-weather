package config

import (
	"flag"
	"log"
	"os"
	"strconv"
)

const (
	// Environmental variable constants
	ConfigDebug             = "WEATHER_LISTEN"
	ConfigListen            = "WEATHER_DEBUG"
	ConfigOpenweathermapKey = "WEATHER_OPENWEATHERMAP_KEY"

	// Defaults
	DefaultListen = ":8080"
)

// Current is the current configuration that should be available to the using application
var Current Configuration

// Init will initialize the configuration from application args and env vars.
func Init() {
	var c Configuration

	// Configure the flag library from defaults first
	flag.BoolVar(&c.Server.Debug, "debug", c.Server.Debug, "enable debug logging")
	flag.StringVar(&c.Server.Listen, "listen", DefaultListen, "listen address for the server")
	flag.StringVar(&c.Openweathermap.Key, "openweathermap-key", c.Openweathermap.Key, "API key for Open Weather Map")
	flag.Parse()

	// Parse ConfigDebug env var if it wasn't set from flag
	if v := os.Getenv(ConfigDebug); len(v) > 0 {
		x, err := strconv.ParseBool(v)
		if err != nil {
			log.Fatalf("error while parsing %q: %s", ConfigDebug, err)
		}
		c.Server.Debug = x
	}

	// Parse ConfigListen env var if it wasn't set from flag
	if v := os.Getenv(ConfigListen); c.Server.Listen != DefaultListen && len(v) > 0 {
		c.Server.Listen = v
	}

	// Parse ConfigOpenweathermapKey env var if it wasn't set from flag
	if v := os.Getenv(ConfigOpenweathermapKey); len(c.Openweathermap.Key) < 1 && len(v) > 0 {
		c.Openweathermap.Key = v
	}

	// Se this as the new current
	Current = c
}

// SanityCheck will check to see if the service is capable of running. If errors will
// be printed to `log` and os.Exit(1) will be called if it is not.
// The purpose of this check is to ensure the application does not start up in a zombie state
// where it starts up but cannot fulfill requests.
func SanityCheck() {
	var fault bool

	// This would be caught by the http.ListenAndServe and prevent the application from
	// running. I wouldn't normally put this check here but interviews be crazy.
	if len(Current.Server.Listen) < 1 {
		log.Println("empty listen address provided")
		fault = true
	}

	if len(Current.Openweathermap.Key) < 1 {
		log.Println("no Open Weather Map key provided")
		fault = true
	}

	if fault {
		os.Exit(1)
	}
}
