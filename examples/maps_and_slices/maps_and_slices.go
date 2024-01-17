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

import "fmt"

func slices() {
	// create an empty slice of integers
	var s1 []int

	// create and populate a slice all in line
	s2 := []int{5, 4, 3, 2}

	// add some stuff to the slice
	s1 = append(s1, 1, 2, 3)

	// get something at a specific index in the slice
	fmt.Printf("thing at index 1 in slice s1 is %d\n", s1[1])
	fmt.Printf("thing at index 1 in slice s2 is %d\n", s2[1])

	// get something out of range this will cause a panic: runtime error: index out of range [127] with length 3
	// fmt.Printf("thing at index 127 in slice s1 is %d\n", s1[127])

	// to do this we must first test if it's ok
	if len(s1) >= 127 {
		fmt.Printf("thing at index 127 in slice s1 is %d\n", s1[127])
	}

	// print out the values of the slice
	for index, value := range s1 {
		fmt.Printf("slice s1 of type %T has %d at index %d\n", s1, index, value)
	}

	for index, value := range s2 {
		fmt.Printf("slice s2 of type %T has %d at index %d\n", s1, index, value)
	}
}

func maps() {
	// create an empty map of unspecified type (technically the map is nil as we'll see later)
	var m1 map[string]interface{}

	// add something to the m1 map
	// m1["foo"] = "bar" will cause a panic: assignment to entry in nil map because we did not initialize the map
	if m1 == nil {
		m1 = make(map[string]interface{})
	}

	// now that the map is initialized it can be appended to it and since it is of type interface we can append whatever
	m1["foo"] = "bar"
	m1["bar"] = 1
	m1["glitch"] = 1.234

	// create a map and initialize it all in one go
	m2 := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	// attempting to put a random string into a strongly typed map will toss an cannot use "four" (untyped string constant) as int value in assignment
	// m2["four"] = "four"

	// iterate through a map
	for key, value := range m2 {
		fmt.Printf("Key %s has a value of  %d\n", key, value)
	}

	// get a specific thing from a map
	if value, ok := m2["one"]; ok {
		fmt.Printf("Value %d was found for key 'one' in the m2 map", value)
	}

	// get a specific thing from a map that's not there
	if value, ok := m2["not there"]; ok {
		// this won't print anything
		fmt.Printf("Value %d was found for key 'not there' in the m2 map", value)
	}
}

func main() {
	slices()
	maps()
}
