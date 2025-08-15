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
	"archive/zip"
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"stl-go/grow-with-stl-go/pkg/configs"
	"stl-go/grow-with-stl-go/pkg/log"
	"stl-go/grow-with-stl-go/pkg/utils"
)

var (
	// ZipcodeCache will expose the zip code info for use in other packages
	ZipcodeCache = make(map[int]*ZipCode)
	// ZipCodeCacheMutex is the controlling mutex for ZipCodeCache
	ZipCodeCacheMutex sync.Mutex
)

// ZipCode contains the extracted city & location details for the data extracted from https://download.geonames.org/export/zip/
type ZipCode struct {
	// the record is in this order
	Country           *string  `json:"country,omitempty"`
	ZipCode           *int     `json:"zip_code,omitempty"`
	City              *string  `json:"city,omitempty"`
	State             *string  `json:"state,omitempty"`
	StateAbbreviation *string  `json:"state_abbreviation,omitempty"`
	County            *string  `json:"county,omitempty"`
	Latitude          *float64 `json:"latitude,omitempty"`
	Longitude         *float64 `json:"longitude,omitempty"`
}

// Init is different than the standard init because it is called outside of the object load
func Init() error {
	go timedTask()
	return getLocations()
}

func getLocations() error {
	// we're going to get our zip code based locations on the country here: https://download.geonames.org/export/zip/
	if configs.GrowSTLGo != nil && configs.GrowSTLGo.Country != nil && configs.GrowSTLGo.Country.Country != nil {
		go func() {
			fileName, err := configs.GrowSTLGo.Country.GetCountryData()
			if err != nil {
				log.Error(err)
				return
			}
			if err := extractLocationData(fileName, configs.GrowSTLGo.Country.Country); err != nil {
				log.Error(err)
			}
		}()

		return nil
	}
	return errors.New("invalid country config cannot get locations")
}

// unzipFile will extract a zip file
func extractLocationData(fileName, country *string) error {
	defer log.FunctionTimer()()
	if fileName != nil && country != nil {
		log.Tracef("attempting to extract location information from the %s archive", *fileName)
		f, fileOpenErr := os.Open(*fileName)
		if fileOpenErr != nil {
			return fileOpenErr
		}
		defer f.Close()

		// Get the file info and create a new ZIP reader
		fi, fileStatError := f.Stat()
		if fileStatError != nil {
			return fileStatError
		}
		zr, zipReaderErr := zip.NewReader(f, fi.Size())
		if zipReaderErr != nil {
			return zipReaderErr
		}

		// Iterate through the files in the ZIP archive
		for _, zf := range zr.File {
			if strings.EqualFold(*country, strings.TrimSuffix(zf.Name, filepath.Ext(zf.Name))) {
				// Open the file
				rc, zipFileErr := zf.Open()
				if zipFileErr != nil {
					return zipFileErr
				}
				defer rc.Close()

				csvReader := csv.NewReader(rc)
				csvReader.Comma = '\t'
				rows, csvErr := csvReader.ReadAll()
				if csvErr != nil {
					return csvErr
				}

				for _, record := range rows {
					if len(record) == 12 {
						zipCode, zipCodeErr := strconv.Atoi(record[1])
						if zipCodeErr != nil {
							log.Error(zipCode)
							continue
						}
						latitude, latitudeErr := strconv.ParseFloat(record[9], 64)
						if latitudeErr != nil {
							log.Error(latitudeErr)
							continue
						}
						longitude, longitudeErr := strconv.ParseFloat(record[10], 64)
						if longitudeErr != nil {
							log.Error(longitudeErr)
							continue
						}
						ZipCodeCacheMutex.Lock()
						ZipcodeCache[zipCode] = &ZipCode{
							Country:           utils.StringPointer(record[0]),
							ZipCode:           &zipCode,
							City:              utils.StringPointer(record[2]),
							State:             utils.StringPointer(record[3]),
							StateAbbreviation: utils.StringPointer(record[4]),
							County:            utils.StringPointer(record[5]),
							Latitude:          &latitude,
							Longitude:         &longitude,
						}
						ZipCodeCacheMutex.Unlock()
					}
				}
			}
		}
		log.Tracef("finished extracting location information from the %s archive.  There are %s zip codes in the cache",
			*fileName, utils.FormatNumber(len(ZipcodeCache)))
		return nil
	}
	return errors.New("invalid input file")
}

func timedTask() {
	seconds := time.Duration((60-time.Now().Local().Second())+120) * time.Second
	log.Infof("location timed tasks will start at %s", time.Now().Add(seconds).Format(configs.HHMMSS))

	for range time.NewTicker(1 * time.Minute).C {
		minute := time.Now().Minute()
		if minute%10 == 0 {
			go func() {
				if err := getLocations(); err != nil {
					log.Errorf("error refreshing locations: %s", err)
				}
			}()
		}
	}
}
