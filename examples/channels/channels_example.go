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
	"os"
	"os/signal"
	"syscall"
	"time"
)

var continueProgram = true

func longRunningThing() {
	fmt.Println("Starting long running thing")
	for {
		if continueProgram {
			time.Sleep(2)
			if time.Now().Second()%5 == 0 {
				fmt.Println(time.Now())
			}
		} else {
			break
		}
	}
	fmt.Println("Exiting long running thing")
}

func launch() {
	fmt.Println("Starting program")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Cleaning up processes prior to exit")
		continueProgram = false
		time.Sleep(2)
		os.Exit(0)
	}()

	longRunningThing()
}

func main() {
	launch()
}
