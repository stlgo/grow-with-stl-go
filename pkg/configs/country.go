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
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"

	"stl-go/grow-with-stl-go/pkg/log"
)

// Country is set for our download of zip codes & latitude / longitudes
type Country struct {
	Country *string `json:"country,omitempty"`
	URL     *string `json:"url,omitempty"`
	ZipFile *string `json:"zip_file,omitempty"`
}

func (c *Config) checkCountry() error {
	if c != nil {
		if c.Country == nil {
			u, urlErr := url.Parse("https://download.geonames.org/export/zip/")
			if urlErr != nil {
				return urlErr
			}
			urlStr := u.String()
			country := "US"
			c.Country = &Country{
				Country: &country,
				URL:     &urlStr,
			}
			rewriteConfig = true
			return nil
		} else if c.Country.Country != nil && c.Country.URL != nil {
			_, urlErr := url.Parse(*c.Country.URL)
			if urlErr != nil {
				return urlErr
			}
			return nil
		}
	}
	return errors.New("invalid config cannot check country")
}

// GetCountryData will get the configured country's city, state, zip, latitude / longitude information
func (c *Country) GetCountryData() (*string, error) {
	defer log.FunctionTimer()()
	if c != nil && c.URL != nil && c.Country != nil && GrowSTLGo.DataDir != nil {
		fileName := filepath.Join(*GrowSTLGo.DataDir, fmt.Sprintf("%s.zip", *c.Country))
		url, urlErr := url.JoinPath(*c.URL, fmt.Sprintf("%s.zip", *c.Country))
		if urlErr != nil {
			return nil, urlErr
		}
		statusCode, downloadErr := DownloadFile(url, http.MethodGet, fileName, nil)
		if downloadErr != nil {
			return nil, downloadErr
		}
		if statusCode != nil {
			if *statusCode < 300 {
				c.ZipFile = &fileName
				return &fileName, nil
			}
			return nil, fmt.Errorf("error code returned from endpoint.  HTTP Status: %d", *statusCode)
		}
		return nil, errors.New("error code returned from get country data endpoint")
	}
	return nil, errors.New("invalid country configuration cannot retrieve data")
}
