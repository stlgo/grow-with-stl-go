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
	"bytes"
	"errors"
	"io"
	"os"

	"stl-go/grow-with-stl-go/pkg/configs"
	"stl-go/grow-with-stl-go/pkg/log"
	"stl-go/grow-with-stl-go/pkg/utils"
)

// Init is different than the standard init because it is called outside of the object load
func Init() error {
	return getLocations()
}

func getLocations() error {
	// we're going to get our zip code based locations on the country here: https://download.geonames.org/export/zip/
	fileName, err := configs.GrowSTLGo.Country.GetCountryData()
	if err != nil {
		return err
	}
	go func() {
		if err := extractLocationData(fileName); err != nil {
			log.Error(err)
		}
	}()
	return nil
}

// unzipFile will extract a zip file
func extractLocationData(fileName *string) error {
	if fileName != nil {
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
			// Open the file
			rc, zipFileErr := zf.Open()
			if zipFileErr != nil {
				return zipFileErr
			}
			defer rc.Close()

			// Print the file name and contents
			_, zipContentsErr := readZipContents(rc)
			if zipContentsErr != nil {
				return zipContentsErr
			}
			// if str != nil {
			// 	// fmt.Printf("File: %s\nContents: %s\n\n", zf.Name, *str)

			// }
		}
		log.Tracef("finished extracting location information from the %s archive", *fileName)
		return nil
	}
	return errors.New("invalid input file")
}

func readZipContents(r io.Reader) (*string, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r); err != nil {
		return nil, err
	}

	return utils.StringPointer(buf.String()), nil
}
