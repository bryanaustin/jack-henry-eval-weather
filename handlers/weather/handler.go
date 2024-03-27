package weather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bryanaustin/jack-henry-eval-weather/config"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const (
	Identifier       = "Weather Handler"
	WeatherAPIMethod = http.MethodGet
	WeatherURL       = "https://api.openweathermap.org/data/2.5/weather"
)

var (
	ErrorArgsRequired     = errors.New("lat and lon arguments are required")
	ErrorArgsInvlaidFloat = errors.New("lat and lon arguments need to be valid float32s")
	ErrorDecodeErrorError = errors.New("unable to decode the error message from a non-200 response")
)

// Handle is the interface for the http server
func Handle(rsp http.ResponseWriter, req *http.Request) {
	lat := req.FormValue("lat")
	lon := req.FormValue("lon")

	// input validation
	if err := validate(lat, lon); err != nil {
		respondError(rsp, err, http.StatusBadRequest)
		return
	}

	// make the request to the weather API
	obj, err := dorequest(req.Context(), lat, lon)
	if err != nil {
		respondError(rsp, err, http.StatusInternalServerError)
		return
	}

	// convert into the desired output
	fin := consolidate(obj)
	enc := json.NewEncoder(rsp)
	err = enc.Encode(fin)
	if err != nil {
		// log because you might have partial output given to the client already
		log.Printf("%s: error encoding response: %s", Identifier, err.Error())
	}
}

// respondError will attempt to return a json error message, on encoding
// failure it will default to a plain text response.
func respondError(rsp http.ResponseWriter, subj error, code int) {
	rsp.Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(ErrorResponse{Message: subj.Error()})
	if err != nil {
		// fallback to plain text error on json error
		body = []byte(subj.Error())
		rsp.Header().Set("Content-Type", "text/plain")
	}
	rsp.WriteHeader(code)
	rsp.Write(body)
}

// validate will validate the input arguments
func validate(lat, lon string) error {
	// validate that something was passed
	if len(lat) < 1 || len(lon) < 1 {
		return ErrorArgsRequired
	}

	// validating that the arguments are valid floats because interviewers love
	// redundancy
	_, laterr := strconv.ParseFloat(lat, 32)
	_, lonerr := strconv.ParseFloat(lon, 32)
	if laterr != nil || lonerr != nil {
		return ErrorArgsInvlaidFloat
	}

	return nil
}

// dorequest will handle the http request
func dorequest(ctx context.Context, lat, lon string) (*OpenWeatherMapResponse, error) {
	// build the URL with params
	u, err := url.Parse(WeatherURL)
	if err != nil {
		// error with context
		return nil, fmt.Errorf("parsing url: %w", err)
	}
	q := url.Values{}
	q.Add("lat", lat)
	q.Add("lon", lon)
	q.Add("appid", config.Current.Openweathermap.Key)
	q.Add("units", "imperial")
	u.RawQuery = q.Encode()

	// create the request with context
	req, err := http.NewRequestWithContext(ctx, WeatherAPIMethod, u.String(), nil)
	if err != nil {
		// error with context
		return nil, fmt.Errorf("building request: %w", err)
	}

	// make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// error with context
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()
	d := json.NewDecoder(resp.Body)

	// non 200 response codes
	if resp.StatusCode != http.StatusOK {
		// try to parse as error message
		obj := new(OpenWeatherMapError)
		err := d.Decode(obj)
		if err != nil {
			// unable to parse error message, return generic
			return nil, ErrorDecodeErrorError
		}
		// able to parse, show message from weather service
		return nil, fmt.Errorf("from weather API: %s", obj.Message)
	}

	// Is a 200 message at this point
	weth := new(OpenWeatherMapResponse)
	if err = d.Decode(weth); err != nil {
		return weth, fmt.Errorf("decoding API response: %w", err)
	}

	return weth, nil
}

// consolidate will convert OpenWeatherMapResponse to Response
// making all the necessary conversions.
func consolidate(i *OpenWeatherMapResponse) (o Response) {
	// set the temperature
	o.Temperature = "moderate"
	// for a generic weather indication the "feels like" seems more apt
	if i.Main.FeelsLike >= 85.0 {
		o.Temperature = "hot"
	} else if i.Main.FeelsLike <= 65.0 {
		o.Temperature = "cold"
	}

	// Set the condition
	if len(i.Weather) >= 1 {
		o.PrimaryCondition = i.Weather[0].Main
	}

	if len(i.Weather) <= 1 {
		// no more conditions to add
		return
	}

	// we have the length, so allocate the exact size
	o.ExtraConditions = make([]string, len(i.Weather)-1)
	for i, x := range i.Weather[1:] {
		o.ExtraConditions[i] = x.Main
	}

	return
}
