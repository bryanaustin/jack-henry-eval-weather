package weather

import (
	"reflect"
	"testing"
)

// TestConsolidationHot tests the hot condition
func TestConsolidationHot(t *testing.T) {
	t.Parallel()
	input := &OpenWeatherMapResponse{
		Main: OpenWeatherMapMain{
			FeelsLike: 100.0,
		},
		Weather: []OpenWeatherMapWeather{
			{
				Main: "Snow",
			},
		},
	}
	expected := Response{
		Temperature:      "hot",
		PrimaryCondition: "Snow",
		ExtraConditions:  nil,
	}
	gotten := consolidate(input)
	if !reflect.DeepEqual(expected, gotten) {
		t.Error("Objects not equal")
	}
}

// TestConsolidationCold tests the cold condition
func TestConsolidationCold(t *testing.T) {
	t.Parallel()
	input := &OpenWeatherMapResponse{
		Main: OpenWeatherMapMain{
			FeelsLike: -1.888,
		},
		Weather: []OpenWeatherMapWeather{
			{
				Main: "Cloudy",
			},
		},
	}
	expected := Response{
		Temperature:      "cold",
		PrimaryCondition: "Cloudy",
		ExtraConditions:  nil,
	}
	gotten := consolidate(input)
	if !reflect.DeepEqual(expected, gotten) {
		t.Error("Objects not equal")
	}
}

// TestConsolidationModerate tests the moderate condition
func TestConsolidationModerate(t *testing.T) {
	t.Parallel()
	input := &OpenWeatherMapResponse{
		Main: OpenWeatherMapMain{
			FeelsLike: 77.77,
		},
		Weather: []OpenWeatherMapWeather{
			{
				Main: "Clear",
			},
		},
	}
	expected := Response{
		Temperature:      "moderate",
		PrimaryCondition: "Clear",
		ExtraConditions:  nil,
	}
	gotten := consolidate(input)
	if !reflect.DeepEqual(expected, gotten) {
		t.Error("Objects not equal")
	}
}

// TestConsolidationExtras tests the extra conditions
func TestConsolidationExtras(t *testing.T) {
	t.Parallel()
	input := &OpenWeatherMapResponse{
		Main: OpenWeatherMapMain{
			FeelsLike: 65.01,
		},
		Weather: []OpenWeatherMapWeather{
			{
				Main: "Cloudy",
			},
			{
				Main: "Tornado",
			},
			{
				Main: "Cats & Dogs",
			},
		},
	}
	expected := Response{
		Temperature:      "moderate",
		PrimaryCondition: "Cloudy",
		ExtraConditions:  []string{"Tornado", "Cats & Dogs"},
	}
	gotten := consolidate(input)
	if !reflect.DeepEqual(expected, gotten) {
		t.Error("Objects not equal")
	}
}
