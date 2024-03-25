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
	"os"
	"os/signal"
	"syscall"
	"time"

	"stl-go/grow-with-stl-go/pkg/log"
)

func init() {
	go timedTask()
}

// this will print a message ever y 10 seconds on the dot starting at the top of the next minute
func timedTask() {
	log.Infof("Waiting till the top of the minute (%d seconds) until starting the timed task", 60-time.Now().Local().Second())
	time.Sleep(time.Duration(60-time.Now().Local().Second()) * time.Second)
	for range time.NewTicker(10 * time.Second).C {
		log.Info("Something done at this time")
	}
}

func main() {
	log.Info("Starting the timed task program")
	closeWait := make(chan byte, 1)

	// since we're running more or less in the background we need to listen for an OS exit call CTL+C or pid kill
	osInterrupt := make(chan os.Signal, 1)
	signal.Notify(osInterrupt, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-osInterrupt
		log.Info("os interrupt received, exiting")
		os.Exit(0)
	}()

	// just a little busy box thing to keep the program running
	<-closeWait
}
