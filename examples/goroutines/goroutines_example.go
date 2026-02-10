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
	"math/rand"
	"sync"
	"time"
)

// function to wait a random set of time and then print a messaged
func randomWaitAction(callType string) {
	// golangci-lint tosses a false positive G404: Use of weak random number generator error so we'll skip that for this line
	seconds := rand.Intn(10) // #nosec
	time.Sleep(time.Duration(seconds) * time.Second)
	fmt.Printf("Function completed after waiting %d seconds for: %s\n", seconds, callType)
}

// this will kick off a bunch of stuff but exit before they completed because they weren't waited for
func runWithoutWait() {
	fmt.Println("Starting goroutines without wait function")
	// lets kick off 10 goroutines but don't wait for them
	for i := 1; i <= 10; i++ {
		go randomWaitAction(fmt.Sprintf("non wait call %d", i))
	}
	fmt.Println("Finished goroutines without wait function")
}

// this will kick off a bunch of stuff and wait for them all to complete
func runWait() {
	fmt.Println("Starting goroutines with wait function")
	// setup a wait group
	var wg sync.WaitGroup
	// kick off 10 goroutines but wait for them with a wait group
	for i := range 10 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			randomWaitAction(fmt.Sprintf("call is waited for %d", i))
		}(i)
	}
	fmt.Println("Done spawning goroutines with wait function")
	// wait until everything is done then exit the function
	wg.Wait()
	fmt.Println("Finished goroutines with wait function")
}

// this will kick off a bunch of stuff and wait for them all to complete but within reason
func runWaitLimiter() {
	fmt.Println("Starting goroutines with wait limit function")
	// setup a wait group
	var wg sync.WaitGroup
	ch := make(chan int, 10)
	// kick off 10 goroutines but wait for them with a wait group
	for i := range 100 {
		wg.Add(1)
		ch <- 1
		go func(i int) {
			defer wg.Done()
			randomWaitAction(fmt.Sprintf("call is waited for %d", i))
			<-ch
		}(i)
	}
	fmt.Println("Done spawning goroutines with wait limit function")
	// wait until everything is done then exit the function
	wg.Wait()
	fmt.Println("Finished goroutines with wait limit function")
}

func main() {
	runWithoutWait()
	runWait()
	runWaitLimiter()
}
