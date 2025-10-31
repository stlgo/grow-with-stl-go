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
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"stl-go/grow-with-stl-go/pkg/configs"
	"stl-go/grow-with-stl-go/pkg/log"
)

type weatherStation struct {
	Agency               *string  `json:"agency,omitempty"`
	SiteNumber           *string  `json:"site_number,omitempty"`
	SiteName             *string  `json:"site_name,omitempty"`
	SiteType             *string  `json:"site_type,omitempty"`
	Latitude             *float64 `json:"latitude,omitempty"`
	Longitude            *float64 `json:"longitude,omitempty"`
	CoordinateAccuracy   *string  `json:"coordinate_accuracy,omitempty"`
	CoordinateDatum      *string  `json:"coordinate_datum,omitempty"`
	GageAltitude         *int     `json:"gage_altitude,omitempty"`
	GageAltitudeAccuracy *string  `json:"gage_altitude_accuracy,omitempty"`
	GageAltitudeDatum    *string  `json:"gage_altitude_datum,omitempty"`
	HydrologicUnitCode   *string  `json:"hydrologic_unit_code,omitempty"`
}

func usgsSiteNumberLookup(lat, long *float64) ([]*weatherStation, error) {
	if lat != nil && long != nil && configs.GrowSTLGo != nil && configs.GrowSTLGo.Weather != nil && configs.GrowSTLGo.Weather.USGSLookup != nil {
		url := fmt.Sprintf("%s/?bBox=%.3f,%.3f,%.3f,%.3f", *configs.GrowSTLGo.Weather.USGSLookup, *lat, *long, *lat+.11, *long+.11)
		response, statusCode, requestErr := configs.HTTPRequest(url, http.MethodGet, nil)
		if requestErr != nil {
			return nil, requestErr
		}
		if statusCode != nil && *statusCode >= 300 || response == nil {
			if response == nil {
				return nil, fmt.Errorf("bad response from url %s.  HTTP Status Code %d", url, *statusCode)
			}
			return nil, fmt.Errorf("bad response from url %s.  Response '%s'.  HTTP Status Code %d", url, *response, *statusCode)
		}

		csvReader := csv.NewReader(strings.NewReader(*response))
		csvReader.Comma = '\t'         // set to read tab delimited data
		csvReader.FieldsPerRecord = -1 // fix wrong number of fields errors

		rows, csvErr := csvReader.ReadAll()
		if csvErr != nil {
			return nil, csvErr
		}

		stations := []*weatherStation{}
		first2Records := 0
		for _, record := range rows {
			if len(record) == 12 {
				// first 2 records are headers and controls
				if first2Records < 2 {
					first2Records++
					continue
				}
				stations = append(stations, record2weatherStation(record))
			}
		}
		return stations, nil
	}
	return nil, errors.New("invalid lat & long cannot lookup usgs site")
}

func record2weatherStation(record []string) *weatherStation {
	// remove all the spaces from the individual records
	for index, datum := range record {
		record[index] = strings.TrimSpace(datum)
	}

	var latitude, longitude *float64
	if len(record[4]) > 0 {
		l, latitudeErr := strconv.ParseFloat(record[4], 64)
		if latitudeErr != nil {
			log.Error(latitudeErr)
		}
		latitude = &l
	}

	if len(record[5]) > 0 {
		l, longitudeErr := strconv.ParseFloat(record[5], 64)
		if longitudeErr != nil {
			log.Error(longitudeErr)
		}
		longitude = &l
	}

	var altitude *int
	if len(record[8]) > 0 {
		a, altitudeErr := strconv.Atoi(record[8])
		if altitudeErr != nil {
			log.Error(altitudeErr)
		}
		altitude = &a
	}

	return &weatherStation{
		Agency:               &record[0],
		SiteNumber:           &record[1],
		SiteName:             &record[2],
		SiteType:             &record[3],
		Latitude:             latitude,
		Longitude:            longitude,
		CoordinateAccuracy:   &record[6],
		CoordinateDatum:      &record[7],
		GageAltitude:         altitude,
		GageAltitudeAccuracy: &record[9],
		GageAltitudeDatum:    &record[10],
		HydrologicUnitCode:   &record[11],
	}
}
