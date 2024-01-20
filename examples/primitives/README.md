# Primitives

There are 4 types of primitives in Go: string, numeric type, bool, and (possibly according to some docs) error.

## String Types

<https://go.dev/ref/spec#String_types>

Example from [primitives_example.go](primitives_example.go)

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

### String Comparisons

Find if a string is equal but do it in a case sensitive way

```go
if s == "This is a string" {
    fmt.Printf("found case sensitive %s\n", s)
}
```

Output

```bash
found case sensitive This is a string
```

Find if a string is equal but do it in a case insensitive way

```go
if sp != nil && strings.EqualFold(*sp, "this is a string") {
    fmt.Printf("found case insensitive %s\n", *sp)
}
```

Output

```bash
found case insensitive This is a string
```

Find if a string has a prefix

```go
if strings.HasPrefix(s, "This") {
    fmt.Printf("found This as a prefix to %s\n", s)
}
```

Output

```bash
found This as a prefix to This is a string
```

Find if a string has a suffix

```go
if sp != nil && strings.HasSuffix(*sp, "ing") {
    fmt.Printf("found ing as a suffix to %s\n", *sp)
}
```

Output

```bash
found ing as a suffix to This is a string
```

## Numeric Types

<https://go.dev/ref/spec#Numeric_types>

Example from [primitives_example.go](primitives_example.go)

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

Example from [primitives_example.go](primitives_example.go)

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

Example from [primitives_example.go](primitives_example.go)

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

You can see all of this, and more, live in our [primitives_example.go](primitives_example.go) example.  To run this example:

```bash
go run examples/primitives/primitives_example.go
```

Output

```bash
bool variable t's value true
*bool variable pointer tp's address 0xc000102068 value of true
bool variable f's value false
*bool variable pointer fp's address 0xc000102088 value of false
bool variable at 0 index of the bools slice has a value of true
bool variable at 1 index of the bools slice has a value of false
bool variable at 2 index of the bools slice has a value of true
bool variable at 3 index of the bools slice has a value of false
bool variable at 4 index of the bools slice has a value of true
bool variable at 5 index of the bools slice has a value of false
int variable i's value of 12345
*int variable pointer ip's address 0xc0001020a0 value of 12345
int variable at 0 index of the ints slice has a value of 1234
int variable at 1 index of the ints slice has a value of 4321
int variable at 2 index of the ints slice has a value of 12345
int variable at 3 index of the ints slice has a value of 12345
float64 variable f's value of 1234.500000
*float64 variable pointer fp's address 0xc000102100 value of 1234.500000
float64 variable at 0 index of the floats slice has a value of 1234.500000
float64 variable at 1 index of the floats slice has a value of 4321.000000
float64 variable at 2 index of the floats slice has a value of 1234.500000
float64 variable at 3 index of the floats slice has a value of 1234.500000
[]uint8 variable ba's value of This is a byte array
*[]uint8 variable pointer bap's address 0xc000100030 value of This is a byte array
[]uint8 variable at 0 index of the byteArrays slice has a value of This is a byte array
[]uint8 variable at 1 index of the byteArrays slice has a value of This is a byte array
[]uint8 variable at 2 index of the byteArrays slice has a value of This is a byte array
string variable s's value of This is a string
*string variable pointer sp's address 0xc000106020 value of This is a string
string variable at 0 index of the floats slice has a value of This is a string
string variable at 1 index of the floats slice has a value of This is a string
string variable at 2 index of the floats slice has a value of This is a string
found case sensitive This is a string
found case insensitive This is a string
found This as a prefix to This is a string
found ing as a suffix to This is a string
*errors.errorString variable err's value this is an error message
*error variable pointer errPtr's address 0xc000106100 value of this is an error message
*errors.errorString variable at 0 index of the errors slice has a value of this is an error message
*errors.errorString variable at 1 index of the errors slice has a value of this is an error message
*errors.errorString variable at 2 index of the errors slice has a value of this is an error message

Calling function main.booleanTypes
bool variable t's value true
*bool variable pointer tp's address 0xc000102160 value of true
bool variable f's value false
*bool variable pointer fp's address 0xc000102161 value of false
bool variable at 0 index of the bools slice has a value of true
bool variable at 1 index of the bools slice has a value of false
bool variable at 2 index of the bools slice has a value of true
bool variable at 3 index of the bools slice has a value of false
bool variable at 4 index of the bools slice has a value of true
bool variable at 5 index of the bools slice has a value of false

Calling function main.numericTypes
int variable i's value of 12345
*int variable pointer ip's address 0xc000102168 value of 12345
int variable at 0 index of the ints slice has a value of 1234
int variable at 1 index of the ints slice has a value of 4321
int variable at 2 index of the ints slice has a value of 12345
int variable at 3 index of the ints slice has a value of 12345
float64 variable f's value of 1234.500000
*float64 variable pointer fp's address 0xc0001021c8 value of 1234.500000
float64 variable at 0 index of the floats slice has a value of 1234.500000
float64 variable at 1 index of the floats slice has a value of 4321.000000
float64 variable at 2 index of the floats slice has a value of 1234.500000
float64 variable at 3 index of the floats slice has a value of 1234.500000

Calling function main.stringTypes
[]uint8 variable ba's value of This is a byte array
*[]uint8 variable pointer bap's address 0xc000100120 value of This is a byte array
[]uint8 variable at 0 index of the byteArrays slice has a value of This is a byte array
[]uint8 variable at 1 index of the byteArrays slice has a value of This is a byte array
[]uint8 variable at 2 index of the byteArrays slice has a value of This is a byte array
string variable s's value of This is a string
*string variable pointer sp's address 0xc000106160 value of This is a string
string variable at 0 index of the floats slice has a value of This is a string
string variable at 1 index of the floats slice has a value of This is a string
string variable at 2 index of the floats slice has a value of This is a string

Calling function main.stringComparisons
found case sensitive This is a string
found case insensitive This is a string
found This as a prefix to This is a string
found ing as a suffix to This is a string

Calling function main.errorTypes
*errors.errorString variable err's value this is an error message
*error variable pointer errPtr's address 0xc000106260 value of this is an error message
*errors.errorString variable at 0 index of the errors slice has a value of this is an error message
*errors.errorString variable at 1 index of the errors slice has a value of this is an error message
*errors.errorString variable at 2 index of the errors slice has a value of this is an error message
```
