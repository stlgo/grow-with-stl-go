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

package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// StructJSON is a basic struct so we can marshall data -> JSON / unmarshall JSON -> data in the program
type StructJSON struct {
	FileDate        *int64          `json:"fileDate,omitempty"`
	FileDateISO8601 *string         `json:"fileDateISO8601,omitempty"`
	FileText        *string         `json:"fileText,omitempty"`
	SomeArray       *[]int          `json:"someArray,omitempty"`
	NestedMap       *map[string]any `json:"nestedMap,omitempty"`
}

// write the struct based JSON struct to disk
func (jo StructJSON) persist() (*string, error) {
	// create the temp dir
	tmpDir, err := makeTempDir()
	if err != nil {
		return nil, err
	}

	if tmpDir != nil {
		fileName := filepath.Join(*tmpDir, "structBased.json.gz")

		// create the gzip file
		fi, err := os.Create(fileName)
		if err != nil {
			return nil, err
		}
		defer fi.Close()

		// create the gzip file writer
		gzw := gzip.NewWriter(fi)
		defer gzw.Close()

		// create the buffered writer
		bfw := bufio.NewWriter(gzw)
		defer bfw.Flush()

		// marshall the StructJSON object to a byte array
		jsonBytes, err := json.Marshal(jo)
		if err != nil {
			return nil, err
		}

		numBytes, err := bfw.Write(jsonBytes)
		if err != nil {
			return nil, err
		}

		fmt.Printf("%d bytes were successfully written to %s\n", numBytes, fileName)
		return &fileName, nil
	}
	return nil, fmt.Errorf("the tmp directory is nil, cannot continue")
}

// Generic JSON with string as a Key and a generic interface as the value
func createSimpleJSON() (*string, error) {
	jo := map[string]any{
		"fileText":        "This text was be written to the file by this example program for simple JSON",
		"fileDate":        time.Now().UnixMilli(),
		"fileDateISO8601": time.Now().UTC().Format("2006-01-02T15:04:05-0700"),
		"someArray":       []int{1, 2, 3, 4},
		"nestedMap": map[string]any{
			"foo":       "bar",
			"someArray": []float32{1.2, 2.4, 3.6, 4.8},
		},
	}

	fmt.Printf("Simple JSON created: %v\n", jo)

	// create the temp dir
	tmpDir, err := makeTempDir()
	if err != nil {
		return nil, err
	}

	// marshall the StructJSON object to a byte array
	jsonBytes, err := json.MarshalIndent(jo, "", "\t")
	if err != nil {
		return nil, err
	}

	// write the byte array to the file
	fileName := filepath.Join(*tmpDir, "simple.json")
	err = os.WriteFile(fileName, jsonBytes, 0o600)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Temp JSON file %s was created and successfully written to\n", fileName)
	return &fileName, nil
}

// read a struct based JSON file
func readSimpleJSONFile(fileName *string) (map[string]any, error) {
	if fileName != nil {
		jsonBytes, err := os.ReadFile(*fileName)
		if err != nil {
			return nil, err
		}

		// unmarshal the file into a basic JSON Object
		var jo map[string]any
		if err1 := json.Unmarshal(jsonBytes, &jo); err1 != nil {
			return nil, err1
		}

		// print it back out as a generic JSON
		jsonOutBytes, err := json.MarshalIndent(jo, "", "\t")
		if err != nil {
			return nil, err
		}
		fmt.Println(string(jsonOutBytes))

		return jo, nil
	}
	return nil, fmt.Errorf("file name is nil, cannot continue")
}

// a helper method to run the functions to create & use a simple JSON file
func runSimpleJSONFunctions() {
	fileName, err := createSimpleJSON()
	if err != nil {
		fmt.Printf("Unable to continue, cannot write a simple JSON file: %s", err)
		os.Exit(-1)
	}

	simpleJSON, err := readSimpleJSONFile(fileName)
	if err != nil {
		fmt.Printf("Unable to continue, cannot read the simple json file %s: %s", *fileName, err)
		os.Exit(-1)
	}

	// you can now interact directly with the simple JSON
	for key, value := range simpleJSON {
		fmt.Printf("Simple JSON key %s has a value of %v\n", key, value)
	}

	// you can also interact with specific keys in the map
	if value, ok := simpleJSON["fileDateISO8601"]; ok {
		fmt.Printf("Simple JSON value of key \"fileDateISO8601\" %s\n", value)
	}

	// nothing will happen here because the key "keyDoesNotExist" doesn't exist
	if value, ok := simpleJSON["keyDoesNotExist"]; ok {
		fmt.Printf("value of keyDoesNotExist %s\n", value)
	}
}

// create a tmp runtime directory to write and get files to
func makeTempDir() (*string, error) {
	tmpDir, err := os.MkdirTemp("", "stl-go")
	if err != nil {
		return nil, err
	}
	fmt.Printf("Temp dir %s was created\n", tmpDir)
	return &tmpDir, nil
}

// create a struct based JSON object
func createStructJSON() (*string, error) {
	fileText := "This text was be written to the file by this example program for struct based JSON"

	// because we're using pointers in our struct we need to create the variables first
	now := time.Now()
	millis := now.UnixMilli()
	iso8601 := now.UTC().Format("2006-01-02T15:04:05-0700")
	someArray := []int{1, 2, 3, 4}
	nesteMap := map[string]any{
		"foo":       "bar",
		"someArray": []float32{1.2, 2.4, 3.6, 4.8},
	}

	// we use the addresses when creating the object
	jo := StructJSON{
		FileText:        &fileText,
		FileDate:        &millis,
		FileDateISO8601: &iso8601,
		SomeArray:       &someArray,
		NestedMap:       &nesteMap,
	}

	// write the file out
	fileName, err := jo.persist()
	if err != nil {
		return nil, err
	}

	return fileName, nil
}

// read a struct based JSON file
func readStructJSONFile(fileName *string) (*StructJSON, error) {
	if fileName != nil {
		// crack open the file
		f, err := os.Open(*fileName)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		// create a gzip file reader on the open file handler
		gzr, err := gzip.NewReader(f)
		if err != nil {
			return nil, err
		}
		defer gzr.Close()

		jsonBytes, err := io.ReadAll(gzr)
		if err != nil {
			return nil, err
		}

		// output the string we read in from the file
		fmt.Printf("Data read from %s is:\n%s\n", *fileName, string(jsonBytes))

		// unmarshal the json as the struct
		var structJSON StructJSON
		if err := json.Unmarshal(jsonBytes, &structJSON); err != nil {
			return nil, err
		}
		return &structJSON, nil
	}
	return nil, fmt.Errorf("file name is nil, cannot continue")
}

func runStructJSONFunctions() {
	fileName, err := createStructJSON()
	if err != nil {
		fmt.Printf("Unable to continue, cannot write a struct based JSON file: %s", err)
		os.Exit(-1)
	}

	structJSON, err := readStructJSONFile(fileName)
	if err != nil {
		fmt.Printf("Unable to continue, cannot read the struct based json file %s: %s", *fileName, err)
		os.Exit(-1)
	}

	// you can now interact directly with the struct
	fmt.Printf("%s text %s\n, milliseconds %d, which is easier to use but harder to read than ISO8601 %s",
		*fileName, *structJSON.FileText, *structJSON.FileDate, *structJSON.FileDateISO8601)
}

func main() {
	runSimpleJSONFunctions()
	runStructJSONFunctions()
}
