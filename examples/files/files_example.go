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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

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

// helper method to run through the simple file functions
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

// why write an uncompressed file when you can write a compressed one instead
func writeGzipFile() (*string, error) {
	// create the temp dir
	tmpDir, err := makeTempDir()
	if err != nil {
		return nil, err
	}

	// write a basic text file
	if tmpDir != nil {
		fileName := filepath.Join(*tmpDir, "gzipFile.txt.gz")
		txt := fmt.Sprintf("This text was be written to the file '%s' by this example program on %s", fileName, time.Now().Format("Mon Jan 2 15:04:05 MST 2006"))

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

		numBytes, err := bfw.WriteString(txt)
		if err != nil {
			return nil, err
		}

		fmt.Printf("%d bytes were written to %s\n", numBytes, fileName)
		return &fileName, nil
	}
	return nil, fmt.Errorf("the tmp directory is nil, cannot continue")
}

func readGzipFile(fileName *string) (*string, error) {
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

		bytes, err := io.ReadAll(gzr)
		if err != nil {
			return nil, err
		}

		txt := string(bytes)
		return &txt, nil
	}
	return nil, fmt.Errorf("file name is nil, cannot continue")
}

func runCompressedFileFunctions() {
	fileName, err := writeGzipFile()
	if err != nil {
		fmt.Printf("Unable to continue, cannot write a simple file: %s", err)
		os.Exit(-1)
	}

	fileText, err := readGzipFile(fileName)
	if err != nil {
		fmt.Printf("Unable to continue, cannot read the file %s: %s", *fileName, err)
		os.Exit(-1)
	}

	if fileText != nil {
		fmt.Printf("Text from %s is as follows: \n%s\n", *fileName, *fileText)
	}
}

func main() {
	runSimpleFileFunctions()
	runCompressedFileFunctions()
}
