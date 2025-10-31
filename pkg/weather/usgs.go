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

package weather

import (
	"errors"
	"fmt"
	"net/http"

	"stl-go/grow-with-stl-go/pkg/configs"
)

func usgsSiteNumberLookup(lat, long *float64) (*string, error) {
	if lat != nil && long != nil && configs.GrowSTLGo != nil && configs.GrowSTLGo.Weather != nil && configs.GrowSTLGo.Weather.USGSLookup != nil {
		url := fmt.Sprintf("%s/?bBox=%.3f,%.3f,%.3f,%.3f", *configs.GrowSTLGo.Weather.USGSLookup, *lat, *long, *lat+.11, *long+.11)
		response, statusCode, requestErr := configs.HTTPRequest(url, http.MethodGet, nil)
		if requestErr != nil {
			return nil, requestErr
		}
		if statusCode != nil && *statusCode >= 300 {
			if response == nil {
				return nil, fmt.Errorf("bad response from url %s.  HTTP Status Code %d", url, *statusCode)
			}
			return nil, fmt.Errorf("bad response from url %s.  Response '%s'.  HTTP Status Code %d", url, *response, *statusCode)
		}
		return response, nil
	}
	return nil, errors.New("invalid lat & long cannot lookup usgs site")
}
