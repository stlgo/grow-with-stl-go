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
	"reflect"
	"runtime"
	"strings"
)

// booleanTypes is a function that flexes basic boolean operations
func booleanTypes() {
	t := true // true boolean
	tp := &t  // pointer to a boolean
	fmt.Printf("%T variable t's value %t\n", t, t)
	fmt.Printf("%T variable pointer tp's address %p value of %t\n", tp, tp, *tp)

	f := false // false boolean
	fp := &f   // pointer to a boolean
	fmt.Printf("%T variable f's value %t\n", f, f)
	fmt.Printf("%T variable pointer fp's address %p value of %t\n", fp, fp, *fp)

	bools := []bool{true, false, t, f, *tp, *fp} // slice (array) of boolean type values
	for index, value := range bools {
		fmt.Printf("%T variable at %d index of the bools slice has a value of %t\n", value, index, value)
	}

	var errBool *bool
	// this will cause - panic: runtime error: invalid memory address or nil pointer dereference
	// fmt.Printf("%t this will error", *errBool)
	// to utilize pointers correctly, regardless of type, it has to have a nil check first
	if errBool != nil {
		fmt.Printf("%t this will not error", *errBool)
	}
}

// numericTypes is a function that flexes basic numeric operations
func numericTypes() {
	i := 12345 // integer
	ip := &i   // pointer to an integer
	fmt.Printf("%T variable i's value of %d\n", i, i)
	fmt.Printf("%T variable pointer ip's address %p value of %d\n", ip, ip, *ip)

	ints := []int{1234, 4321, i, *ip} // slice (array) of integer type values
	for index, value := range ints {
		fmt.Printf("%T variable at %d index of the ints slice has a value of %d\n", value, index, value)
	}

	var errInt *int
	// this will cause - panic: runtime error: invalid memory address or nil pointer dereference
	// fmt.Printf("%d this will error", *errInt)
	// to utilize pointers correctly, regardless of type, it has to have a nil check first
	if errInt != nil {
		fmt.Printf("%d this will not error", *errInt)
	}

	f := 1234.5 // float
	fp := &f    // pointer to a float
	fmt.Printf("%T variable f's value of %f\n", f, f)
	fmt.Printf("%T variable pointer fp's address %p value of %f\n", fp, fp, *fp)

	floats := []float64{1234.5, 4321.0, f, *fp} // slice (array) of float type values
	for index, value := range floats {
		fmt.Printf("%T variable at %d index of the floats slice has a value of %f\n", value, index, value)
	}

	var errFloat *float64
	// this will cause - panic: runtime error: invalid memory address or nil pointer dereference
	// fmt.Printf("%f this will error", *errFloat)
	// to utilize pointers correctly, regardless of type, it has to have a nil check first
	if errFloat != nil {
		fmt.Printf("%f this will not error", *errFloat)
	}
}

// stringTypes is a function that flexes basic string operations
func stringTypes() {
	ba := []byte("This is a byte array") // byte array
	bap := &ba                           // pointer to a byte array
	fmt.Printf("%T variable ba's value of %s\n", ba, ba)
	fmt.Printf("%T variable pointer bap's address %p value of %s\n", bap, bap, *bap)

	byteArrays := [][]byte{[]byte("This is a byte array"), ba, *bap} // slice (array) of byte array type values
	for index, value := range byteArrays {
		fmt.Printf("%T variable at %d index of the byteArrays slice has a value of %s\n", value, index, value)
	}

	var errBa *[]byte
	// this will cause - panic: runtime error: invalid memory address or nil pointer dereference
	// fmt.Printf("%s this will error", *errBa)
	// to utilize pointers correctly, regardless of type, it has to have a nil check first
	if errBa != nil {
		fmt.Printf("%s this will not error", *errBa)
	}

	s := "This is a string" // string
	sp := &s                // pointer to a string
	fmt.Printf("%T variable s's value of %s\n", s, s)
	fmt.Printf("%T variable pointer sp's address %p value of %s\n", sp, sp, *sp)

	stringSlice := []string{"This is a string", s, *sp} // slice (array) of string type values note strings is a reserved word
	for index, value := range stringSlice {
		fmt.Printf("%T variable at %d index of the floats slice has a value of %s\n", value, index, value)
	}

	var errString *string
	// this will cause - panic: runtime error: invalid memory address or nil pointer dereference
	// fmt.Printf("%s this will error", *errString)
	// to utilize pointers correctly, regardless of type, it has to have a nil check first
	if errString != nil {
		fmt.Printf("%s this will not error", *errString)
	}
}

func stringComparisons() {
	s := "This is a string" // string
	sp := &s                // pointer to a string

	// string comparison case sensitive
	if s == "This is a string" {
		fmt.Printf("found case sensitive %s\n", s)
	}

	// string comparison case insensitive
	if sp != nil && strings.EqualFold(*sp, "this is a string") {
		fmt.Printf("found case insensitive %s\n", *sp)
	}

	// string has a prefix
	if strings.HasPrefix(s, "This") {
		fmt.Printf("found This as a prefix to %s\n", s)
	}

	// string has a suffix
	if sp != nil && strings.HasSuffix(*sp, "ing") {
		fmt.Printf("found ing as a suffix to %s\n", *sp)
	}
}

// errorTypes is a function that flexes basic error operations
func errorTypes() {
	err := fmt.Errorf("this is an error message")
	errPtr := &err // pointer to an error
	fmt.Printf("%T variable err's value %s\n", err, err)
	fmt.Printf("%T variable pointer errPtr's address %p value of %s\n", errPtr, errPtr, *errPtr)

	errors := []error{fmt.Errorf("this is an error message"), err, *errPtr} // slice (array) of error type values
	for index, value := range errors {
		fmt.Printf("%T variable at %d index of the errors slice has a value of %s\n", value, index, value)
	}

	var errVar *error
	// this will cause - panic: runtime error: invalid memory address or nil pointer dereference
	// fmt.Printf("%s this will error", *errVar)
	// to utilize pointers correctly, regardless of type, it has to have a nil check first
	if errVar != nil {
		fmt.Printf("%t this will not error", *errVar)
	}
}

func main() {
	// you can call the functions individually
	booleanTypes()
	numericTypes()
	stringTypes()
	stringComparisons()
	errorTypes()

	// or you can have a bit of fun and let the program do it for us with a slice of functions
	functions := []func(){booleanTypes, numericTypes, stringTypes, stringComparisons, errorTypes}
	for _, function := range functions {
		// using reflection to get the name of the function it's calling
		// the name in the slice isn't a string it's a function reference
		fmt.Printf("\nCalling function %s\n", runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name())
		function()
	}
}
