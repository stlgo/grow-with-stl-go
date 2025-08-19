/*
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package locations

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"stl-go/grow-with-stl-go/pkg/configs"
)

type latLongLookup struct {
	Properties *forecastURL `json:"properties,omitempty"`
}

type forecastURL struct {
	ForecastURL *string `json:"forecast,omitempty"`
}

type forecastLookup struct {
	Properties *periods `json:"properties,omitempty"`
}

type periods struct {
	Periods []*Forecast `json:"periods,omitempty"`
}

// Forecast will be held with the zip code lookups for display purposes
type Forecast struct {
	Number                     *int                        `json:"number,omitempty"`
	Name                       *string                     `json:"name,omitempty"`
	StartTime                  *string                     `json:"startTime,omitempty"`
	EndTime                    *string                     `json:"endTime,omitempty"`
	IsDaytime                  *bool                       `json:"isDaytime,omitempty"`
	Temperature                *int                        `json:"temperature,omitempty"`
	TemperatureUnit            *string                     `json:"temperatureUnit,omitempty"`
	TemperatureTrend           *string                     `json:"temperatureTrend,omitempty"`
	ProbabilityOfPrecipitation *probabilityOfPrecipitation `json:"probabilityOfPrecipitation,omitempty"`
	WindSpeed                  *string                     `json:"windSpeed,omitempty"`
	WindDirection              *string                     `json:"windDirection,omitempty"`
	Icon                       *string                     `json:"icon,omitempty"`
	ShortForecast              *string                     `json:"shortForecast,omitempty"`
	DetailedForecast           *string                     `json:"detailedForecast,omitempty"`
}

type probabilityOfPrecipitation struct {
	UnitCode *string `json:"unitCode,omitempty"`
	Value    *int    `json:"value,omitempty"`
}

func getWeather() error {
	if configs.GrowSTLGo != nil && configs.GrowSTLGo.Users != nil {
		lookups := make(map[string]*ZipCode)
		for _, user := range configs.GrowSTLGo.Users {
			if user.Location != nil {
				ZipCodeCacheMutex.Lock()
				zip, ok := ZipcodeLookup[*user.Location]
				ZipCodeCacheMutex.Unlock()
				if ok {
					lookups[*user.Location] = zip
				}
			}
		}
		return getWeatherHelper(lookups)
	}
	return errors.New("cannot get weather, configuration is invalid")
}

func getWeatherHelper(lookups map[string]*ZipCode) error {
	if configs.GrowSTLGo != nil && configs.GrowSTLGo.WeatherAPI != nil {
		for _, zip := range lookups {
			if zip.Latitude != nil && zip.Longitude != nil {
				url, urlErr := url.JoinPath(*configs.GrowSTLGo.WeatherAPI, "points", fmt.Sprintf("%f,%f", *zip.Latitude, *zip.Longitude))
				if urlErr != nil {
					return urlErr
				}
				response, statusCode, requestErr := configs.HTTPRequest(url, http.MethodGet, nil)
				if requestErr != nil {
					return requestErr
				}
				if statusCode != nil && *statusCode >= 300 {
					if response == nil {
						return fmt.Errorf("bad response from url %s.  HTTP Status Code %d", url, *statusCode)
					}
					return fmt.Errorf("bad response from url %s.  Response '%s'.  HTTP Status Code %d", url, *response, *statusCode)
				}
				var lookup latLongLookup
				unmarshalErr := json.Unmarshal([]byte(*response), &lookup)
				if unmarshalErr != nil {
					return unmarshalErr
				}

				if lookup.Properties != nil && lookup.Properties.ForecastURL != nil {
					if forecastErr := getForecast(lookup.Properties.ForecastURL, zip); forecastErr != nil {
						return forecastErr
					}
				}
			}
		}
		return nil
	}
	return errors.New("invalid configuration cannot get weather")
}

func getForecast(url *string, zip *ZipCode) error {
	if url != nil && zip != nil {
		response, statusCode, requestErr := configs.HTTPRequest(*url, http.MethodGet, nil)
		if requestErr != nil {
			return requestErr
		}
		if statusCode != nil && *statusCode >= 300 {
			if response == nil {
				return fmt.Errorf("bad response from url %s.  HTTP Status Code %d", *url, *statusCode)
			}
			return fmt.Errorf("bad response from url %s.  Response '%s'.  HTTP Status Code %d", *url, *response, *statusCode)
		}
		var lookup forecastLookup
		unmarshalErr := json.Unmarshal([]byte(*response), &lookup)
		if unmarshalErr != nil {
			return unmarshalErr
		}
		if lookup.Properties != nil && len(lookup.Properties.Periods) > 0 {
			zip.Mutex.Lock()
			zip.Forecast = lookup.Properties.Periods[0]
			zip.Mutex.Unlock()
		}
		return nil
	}
	return errors.New("invalid url or zip cannot get weather forecast")
}
