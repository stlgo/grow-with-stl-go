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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SimpleJSON is a basic struct so we can marshall data -> JSON / unmarshall JSON -> data in the program
type SimpleJSON struct {
	FileDate        *int64  `json:"fileDate,omitempty"`
	FileDateISO8601 *string `json:"fileDateISO8601,omitempty"`
	FileText        *string `json:"fileText,omitempty"`
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

// write a simple txt file
func writeSimpleFile() (*string, error) {
	// create the temp dir
	tmpDir, err := makeTempDir()
	if err != nil {
		return nil, err
	}

	// write a basic text file
	if tmpDir != nil {
		fileName := filepath.Join(*tmpDir, "simpleFile.txt")
		txt := fmt.Sprintf("This text was be written to the file '%s' by this example program on %s", fileName, time.Now().Format("Mon Jan 2 15:04:05 MST 2006"))
		err := os.WriteFile(fileName, []byte(txt), 0o600)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Temp file %s was created and successfully written to\n", fileName)
		return &fileName, nil
	}
	return nil, fmt.Errorf("directory is nil, cannot continue")
}

// read a simple txt file
func readSimpleFile(fileName *string) (*string, error) {
	if fileName != nil {
		bytes, err := os.ReadFile(*fileName)
		if err != nil {
			return nil, err
		}
		txt := string(bytes)
		return &txt, nil
	}
	return nil, fmt.Errorf("file name is nil, cannot continue")
}

// function to run through the simple file functions
func runSimpleFileFunctions() {
	fileName, err := writeSimpleFile()
	if err != nil {
		fmt.Printf("Unable to continue, cannot write a simple file: %s", err)
		os.Exit(-1)
	}

	fileText, err := readSimpleFile(fileName)
	if err != nil {
		fmt.Printf("Unable to continue, cannot read the file %s: %s", *fileName, err)
		os.Exit(-1)
	}

	if fileText != nil {
		fmt.Printf("Text from %s is as follows: \n%s\n", *fileName, *fileText)
	}
}

// create a simple JSON by creating a struct
func createSimpleJSON() (*string, error) {
	fileText := "This text was be written to the file by this example program"

	now := time.Now()
	millis := now.UnixMilli()
	iso8601 := now.UTC().Format("2006-01-02T15:04:05-0700")

	jo := SimpleJSON{
		FileText: &fileText,

		FileDate:        &millis,
		FileDateISO8601: &iso8601,
	}

	fileName, err := jo.persist()
	if err != nil {
		return nil, err
	}

	return fileName, nil
}

// read a simple json file
func readSimpleJSONFile(fileName *string) (*SimpleJSON, error) {
	if fileName != nil {
		jsonBytes, err := os.ReadFile(*fileName)
		if err != nil {
			return nil, err
		}

		// unmarshal the file into a basic JSON Object
		var jo map[string]interface{}
		if err1 := json.Unmarshal(jsonBytes, &jo); err1 != nil {
			return nil, err1
		}

		// print it back out as a generic JSON
		jsonOutBytes, err := json.MarshalIndent(jo, "", "\t")
		if err != nil {
			return nil, err
		}
		fmt.Println(string(jsonOutBytes))

		// unmarshal the json as the struct
		var simpleJSON SimpleJSON
		if err := json.Unmarshal(jsonBytes, &simpleJSON); err != nil {
			return nil, err
		}
		return &simpleJSON, nil
	}
	return nil, fmt.Errorf("file name is nil, cannot continue")
}

// write the simple JSON struct to disk
func (jo SimpleJSON) persist() (*string, error) {
	// create the temp dir
	tmpDir, err := makeTempDir()
	if err != nil {
		return nil, err
	}

	// marshall the SimpleJSON object to a byte array
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

	// you can now interact directly with the struct
	fmt.Printf("%s text %s\n, milliseconds %d, which is easier to use but harder to read than ISO8601 %s",
		*fileName, *simpleJSON.FileText, *simpleJSON.FileDate, *simpleJSON.FileDateISO8601)
}

func main() {
	runSimpleFileFunctions()
	runSimpleJSONFunctions()
}
