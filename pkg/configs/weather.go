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

package configs

import (
	"errors"

	"stl-go/grow-with-stl-go/pkg/utils"
)

// Weather contains the lookups for the weather functions
type Weather struct {
	USGSLookup *string `json:"usgs_lookup,omitempty"`
	WeatherAPI *string `json:"weather_api,omitempty"`
}

func (c *Config) checkWeather() error {
	if c != nil {
		if c.Weather == nil {
			weatherAPI, err := utils.URLParserHelper("https://api.weather.gov")
			if err != nil {
				return err
			}
			usgsLookup, err := utils.URLParserHelper("https://maps.waterdata.usgs.gov/mapper/nwis/site/")
			if err != nil {
				return err
			}
			c.Weather = &Weather{
				USGSLookup: usgsLookup,
				WeatherAPI: weatherAPI,
			}
			rewriteConfig = true
		}
	}
	return errors.New("invalid config cannot check Weather")
}
