package config

// Configuration is the current configuration for the service
type Configuration struct {
	Openweathermap struct {
		Key string
	}
	Server struct {
		Debug  bool
		Listen string
	}
}
