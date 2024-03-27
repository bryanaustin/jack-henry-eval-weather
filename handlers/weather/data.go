package weather

// OpenWeatherMapWeather is the "weather" stanza in the response struct
type OpenWeatherMapWeather struct {
	Main string `json:"main"`
}

type OpenWeatherMapMain struct {
	FeelsLike float32 `json:"feels_like"`
}

// OpenWeatherMapResponse is the response to the /data/2.5/weather path
type OpenWeatherMapResponse struct {
	Main    OpenWeatherMapMain      `json:"main"`
	Weather []OpenWeatherMapWeather `json:"weather"`
}

type OpenWeatherMapError struct {
	Code    string `json:"cod"`
	Message string `json:"message"`
}

// Response is what an OK response is from this handler
type Response struct {
	Temperature      string
	PrimaryCondition string
	ExtraConditions  []string
}

type ErrorResponse struct {
	Message string
}
