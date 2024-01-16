# Primitives

There are 4 types of primitives in Go: string, numeric type, bool, and (possibly according to some docs) error.

## String Types

<https://go.dev/ref/spec#String_types>

Example from [primitives_example.go](../examples/primitives/primitives_example.go)

```go
s := "This is a string" // string
sp := &s                // pointer to a string
fmt.Printf("%T variable s's value of %s\n", s, s)
fmt.Printf("%T variable pointer sp's address %p value of %s\n", sp, sp, *sp)
```

Output

```bash
string variable s's value of This is a string
*string variable pointer sp's address 0xc000026170 value of This is a string
```

## Numeric Types

<https://go.dev/ref/spec#Numeric_types>

Example from [primitives_example.go](../examples/primitives/primitives_example.go)

```go
i := 12345 // integer
ip := &i   // pointer to an integer
fmt.Printf("%T variable i's value of %d\n", i, i)
fmt.Printf("%T variable pointer ip's address %p value of %d\n", ip, ip, *ip)
```

Output

```bash
int variable i's value of 12345
*int variable pointer ip's address 0xc00000a198 value of 12345
```

## Boolean Types

<https://go.dev/ref/spec#Boolean_types>

Example from [primitives_example.go](../examples/primitives/primitives_example.go)

```go
t := true // true boolean
tp := &t  // pointer to a boolean
fmt.Printf("%T variable t's value %t\n", t, t)
fmt.Printf("%T variable pointer tp's address %p value of %t\n", tp, tp, *tp)
```

Output

```bash
bool variable t's value true
*bool variable pointer tp's address 0xc00000a190 value of true
```

## Error Types

<https://pkg.go.dev/builtin#error>

Example from [primitives_example.go](../examples/primitives/primitives_example.go)

```go
err := fmt.Errorf("this is an error message")
errPtr := &err // pointer to an error
fmt.Printf("%T variable err's value %s\n", err, err)
fmt.Printf("%T variable pointer errPtr's address %p value of %s\n", errPtr, errPtr, *errPtr)
```

Output

```bash
*errors.errorString variable err's value this is an error message
*error variable pointer errPtr's address 0xc000026110 value of this is an error message
```

## Example

You can see all of this, and more, live in our [primitives_example.go](../examples/primitives/primitives_example.go) example.  To run this example:

```bash
go run examples/primitives/primitives_example.go
```

Output

```bash
bool variable t's value true
*bool variable pointer tp's address 0xc00000a098 value of true
bool variable f's value false
*bool variable pointer fp's address 0xc00000a0b8 value of false
bool variable at 0 index of the bools slice has a value of true
bool variable at 1 index of the bools slice has a value of false
bool variable at 2 index of the bools slice has a value of true
bool variable at 3 index of the bools slice has a value of false
bool variable at 4 index of the bools slice has a value of true
bool variable at 5 index of the bools slice has a value of false
int variable i's value of 12345
*int variable pointer ip's address 0xc00000a0d0 value of 12345
int variable at 0 index of the ints slice has a value of 1234
int variable at 1 index of the ints slice has a value of 4321
int variable at 2 index of the ints slice has a value of 12345
int variable at 3 index of the ints slice has a value of 12345
float64 variable f's value of 1234.500000
*float64 variable pointer fp's address 0xc00000a130 value of 1234.500000
float64 variable at 0 index of the floats slice has a value of 1234.500000
float64 variable at 1 index of the floats slice has a value of 4321.000000
float64 variable at 2 index of the floats slice has a value of 1234.500000
float64 variable at 3 index of the floats slice has a value of 1234.500000
[]uint8 variable ba's value of This is a byte array
*[]uint8 variable pointer bap's address 0xc000008048 value of This is a byte array
[]uint8 variable at 0 index of the byteArrays slice has a value of This is a byte array
[]uint8 variable at 1 index of the byteArrays slice has a value of This is a byte array
[]uint8 variable at 2 index of the byteArrays slice has a value of This is a byte array
string variable s's value of This is a string
*string variable pointer sp's address 0xc000026070 value of This is a string
string variable at 0 index of the floats slice has a value of This is a string
string variable at 1 index of the floats slice has a value of This is a string
string variable at 2 index of the floats slice has a value of This is a string
*errors.errorString variable err's value this is an error message
*error variable pointer errPtr's address 0xc000026110 value of this is an error message
*errors.errorString variable at 0 index of the errors slice has a value of this is an error message
*errors.errorString variable at 1 index of the errors slice has a value of this is an error message
*errors.errorString variable at 2 index of the errors slice has a value of this is an error message

Calling function main.booleanTypes
bool variable t's value true
*bool variable pointer tp's address 0xc00000a190 value of true
bool variable f's value false
*bool variable pointer fp's address 0xc00000a191 value of false
bool variable at 0 index of the bools slice has a value of true
bool variable at 1 index of the bools slice has a value of false
bool variable at 2 index of the bools slice has a value of true
bool variable at 3 index of the bools slice has a value of false
bool variable at 4 index of the bools slice has a value of true
bool variable at 5 index of the bools slice has a value of false

Calling function main.numericTypes
int variable i's value of 12345
*int variable pointer ip's address 0xc00000a198 value of 12345
int variable at 0 index of the ints slice has a value of 1234
int variable at 1 index of the ints slice has a value of 4321
int variable at 2 index of the ints slice has a value of 12345
int variable at 3 index of the ints slice has a value of 12345
float64 variable f's value of 1234.500000
*float64 variable pointer fp's address 0xc00000a1f8 value of 1234.500000
float64 variable at 0 index of the floats slice has a value of 1234.500000
float64 variable at 1 index of the floats slice has a value of 4321.000000
float64 variable at 2 index of the floats slice has a value of 1234.500000
float64 variable at 3 index of the floats slice has a value of 1234.500000

Calling function main.stringTypes
[]uint8 variable ba's value of This is a byte array
*[]uint8 variable pointer bap's address 0xc000008138 value of This is a byte array
[]uint8 variable at 0 index of the byteArrays slice has a value of This is a byte array
[]uint8 variable at 1 index of the byteArrays slice has a value of This is a byte array
[]uint8 variable at 2 index of the byteArrays slice has a value of This is a byte array
string variable s's value of This is a string
*string variable pointer sp's address 0xc000026170 value of This is a string
string variable at 0 index of the floats slice has a value of This is a string
string variable at 1 index of the floats slice has a value of This is a string
string variable at 2 index of the floats slice has a value of This is a string

Calling function main.errorTypes
*errors.errorString variable err's value this is an error message
*error variable pointer errPtr's address 0xc000026220 value of this is an error message
*errors.errorString variable at 0 index of the errors slice has a value of this is an error message
*errors.errorString variable at 1 index of the errors slice has a value of this is an error message
*errors.errorString variable at 2 index of the errors slice has a value of this is an error message
```
