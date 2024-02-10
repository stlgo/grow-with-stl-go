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
	"crypto/sha256"
	"errors"
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
		// sha256 hash compare first
		if err := compareSHA256Sum(fileName); err != nil {
			return nil, err
		}

		// read the file
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
	// create a file
	fileName, err := writeSimpleFile()
	if err != nil || fileName == nil {
		fmt.Printf("Unable to continue, cannot write a simple file: %s", err)
		os.Exit(-1)
	}

	// sum has to be calculated after the file buffers have flushed and the file handlers have closed
	if err = writeSHA256Sum(fileName); err != nil {
		fmt.Printf("Unable to continue, cannot write a the sum of a simple file: %s", err)
		os.Exit(-1)
	}

	// read the file
	fileText, err := readSimpleFile(fileName)
	if err != nil {
		fmt.Printf("Unable to continue, cannot read the file %s: %s", *fileName, err)
		os.Exit(-1)
	}

	// output the file
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

// read a compressed file
func readGzipFile(fileName *string) (*string, error) {
	if fileName != nil {
		// sha256 hash compare first
		if err := compareSHA256Sum(fileName); err != nil {
			return nil, err
		}

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

// helper method to run through the compressed file functions
func runCompressedFileFunctions() {
	// create the file
	fileName, err := writeGzipFile()
	if err != nil {
		fmt.Printf("Unable to continue, cannot write a compressed file: %s", err)
		os.Exit(-1)
	}

	// sum has to be calculated after the file buffers have flushed and the file handlers have closed
	if err = writeSHA256Sum(fileName); err != nil {
		fmt.Printf("Unable to continue, cannot write a the sum of a compressed file: %s", err)
		os.Exit(-1)
	}

	// read the file
	fileText, err := readGzipFile(fileName)
	if err != nil {
		fmt.Printf("Unable to continue, cannot read the file %s: %s", *fileName, err)
		os.Exit(-1)
	}

	// output the file
	if fileText != nil {
		fmt.Printf("Text from %s is as follows: \n%s\n", *fileName, *fileText)
	}
}

// compare the sha256 hash from a file to what's on disk
func compareSHA256Sum(fileName *string) error {
	if fileName != nil {
		shaSum, err := getSHA256Sum(fileName)
		if err != nil || shaSum == nil {
			fmt.Printf("Unable to continue, cannot sha sum of file: %s", err)
			os.Exit(-1)
		}

		sumFileName := fmt.Sprintf("%s.sha256", *fileName)
		bytes, err := os.ReadFile(sumFileName)
		if err != nil {
			return err
		}
		shaSumFromFile := string(bytes)
		if *shaSum != shaSumFromFile {
			return fmt.Errorf("file %s has a different sha256 hash %s from the stored hash %s", *fileName, *shaSum, shaSumFromFile)
		}

		fmt.Printf("File %s hash is the same as the one stored in %s\n", *fileName, sumFileName)
		return nil
	}
	return errors.New("nil filename cannot compare sha256 sums")
}

// write the sha256 sum of the file to disk along side the original file
func writeSHA256Sum(fileName *string) error {
	if fileName != nil {
		shaSum, err := getSHA256Sum(fileName)
		if err != nil || shaSum == nil {
			fmt.Printf("Unable to continue, cannot sha sum of file: %s", err)
			os.Exit(-1)
		}

		fmt.Printf("SHA256 sum of %s is %s\n", *fileName, *shaSum)
		sha256SumFileName := fmt.Sprintf("%s.sha256", *fileName)
		err = os.WriteFile(sha256SumFileName, []byte(*shaSum), 0o600)
		if err != nil {
			return err
		}
		fmt.Printf("SHA256 Sum of temp file %s was created and successfully written to %s\n", *fileName, sha256SumFileName)
		return nil
	}
	return errors.New("nil filename cannot write sha256 sum")
}

// calculate the sha256 sum from a file
func getSHA256Sum(fileName *string) (*string, error) {
	if fileName != nil {
		f, err := os.Open(*fileName)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		hash := sha256.New()
		if _, err := io.Copy(hash, f); err != nil {
			return nil, err
		}

		sum := fmt.Sprintf("%x", hash.Sum(nil))
		return &sum, nil
	}
	return nil, errors.New("nil filename cannot calculate sha256 sum")
}

func main() {
	runSimpleFileFunctions()
	runCompressedFileFunctions()
}
