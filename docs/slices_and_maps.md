# Slices and Maps

## Slices

<https://go.dev/tour/moretypes/7>

A slice in go is a dynamically-sized, flexible view into the elements of an array

Example from our [maps_and_slices.go](../examples/maps_and_slices/maps_and_slices.go) file:

### Create a slice

```go
// create an empty slice of integers
var s1 []int

// create and populate a slice all in one line
s2 := []int{5, 4, 3, 2}
```

### Append to a slice

```go
// add some stuff to the slice
s1 = append(s1, 1, 2, 3)
```

### Get something from a specific index of a slice
```go
fmt.Printf("thing at index 1 in slice s1 is %d\n", s1[1])
```

Output

```bash
thing at index 1 in slice s1 is 2
```

### Iterate through the slice

```go
for index, value := range s1 {
    fmt.Printf("slice s1 of type %T has %d at index %d\n", s1, index, value)
}
```

Output

```bash
slice s1 of type []int has 0 at index 1
slice s1 of type []int has 1 at index 2
slice s1 of type []int has 2 at index 3
```

## Maps

<https://go.dev/tour/moretypes/19>

### Create a map

```go
// create an empty map of unspecified type (technically the map is nil as we'll see later)
var m1 map[string]interface{}

// create a map and initialize it all in one go
m2 := map[string]int{
    "one":   1,
    "two":   2,
    "three": 3,
}
```

### Append to a map

```go
// add something to the m1 map
// m1["foo"] = "bar" will cause a panic: assignment to entry in nil map because we did not initialize the map
if m1 == nil {
    m1 = make(map[string]interface{})
}

// now that the map is initialized it can be appended to it and since it is of type interface we can append whatever
m1["foo"] = "bar"
```

### Iterate through a map

```go
for key, value := range m1 {
    fmt.Printf("Key %s has a value of  %s\n", key, value)
}
```

Output

```bash
Key one has a value of  1
Key two has a value of  2
Key three has a value of  3
```

### Get a specific item from a map

```go
if value, ok := m2["one"]; ok {
    fmt.Printf("Value %d was found for key 'one' in the m2 map", value)
}
```

Output

```bash
Value 1 was found for key 'one' in the m2 map
```
