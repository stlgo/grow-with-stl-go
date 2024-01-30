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
	"fmt"

	"stl-go/grow-with-stl-go/pkg/log"
)

var (
	logLevels = map[int]string{
		1: "fatal",
		2: "error",
		3: "warn",
		4: "info",
		5: "debug",
		6: "trace",
	}

	logFunctions = []func(string){
		outputPrintln,
		outputLoggingTrace,
		outputLoggingDebug,
		outputLoggingInfo,
		outputLoggingWarn,
		outputLoggingError,
	}

	logFunctionsWithFatal = []func(string){
		outputPrintln,
		outputLoggingTrace,
		outputLoggingDebug,
		outputLoggingInfo,
		outputLoggingWarn,
		outputLoggingError,
		outpuLoggingFatal,
	}
)

func outputPrintln(someMessage string) {
	fmt.Printf("Example output with message: %s\n", someMessage)
}

func outputLoggingTrace(someMessage string) {
	log.Tracef("Example trace output with message: %s", someMessage)
}

func outputLoggingDebug(someMessage string) {
	log.Debugf("Example debug output with message: %s", someMessage)
}

func outputLoggingInfo(someMessage string) {
	log.Infof("Example info output with message: %s", someMessage)
}

func outputLoggingWarn(someMessage string) {
	log.Warnf("Example warn output with message: %s", someMessage)
}

func outputLoggingError(someMessage string) {
	log.Errorf("Example error output with message: %s", someMessage)
}

func outpuLoggingFatal(someMessage string) {
	log.Fatalf("Example fatal output with message: %s", someMessage)
}

func testWithoutFatal() {
	defer log.FunctionTimer()() // this will display how long this took to run at the end of the execution
	logLevelKeys := []int{6, 5, 4, 3, 2}

	for _, logLevel := range logLevelKeys {
		if levelName, ok := logLevels[logLevel]; ok && len(levelName) > 0 {
			log.LogLevel = logLevel
			for _, function := range logFunctions {
				function(fmt.Sprintf("Log attempt for level %s - log level %d", levelName, logLevel))
			}
		}
	}

	// set logging back to trace so we can see the output of the function timer
	log.LogLevel = 6
}

func testWithFatal() {
	defer log.FunctionTimer()() // this will display how long this took to run at the end of the execution but will not work with a fatal
	// maps are randomized by go so we need to iterate through a slice of keys
	logLevelKeys := []int{6, 5, 4, 3, 2, 1}

	for _, logLevel := range logLevelKeys {
		if levelName, ok := logLevels[logLevel]; ok && len(levelName) > 0 {
			log.LogLevel = logLevel
			for _, function := range logFunctionsWithFatal {
				function(fmt.Sprintf("Log attempt for level %s - log level %d", levelName, logLevel))
			}
		}
	}

	// set logging back to trace so we can see the output of the function timer
	// this function will exit prior to getting to this point
	log.LogLevel = 6
}

func main() {
	// all of these function calls will work and the function timer will fire
	log.Info("Sending all log attempts without fatal")
	testWithoutFatal()

	// this will exit the program once it reaches the first fatal log statement
	log.Info("Sending all log attempts with fatal")
	testWithFatal()
}
