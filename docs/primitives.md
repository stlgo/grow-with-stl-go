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
