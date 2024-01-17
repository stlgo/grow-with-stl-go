# Structures

<https://go.dev/tour/moretypes/2>

In golang a structure (struct) is a collection of things you want to store off and use.  It would be considered, in object oriented terms, a class.

Example from our [struct_example.go](../examples/structures/struct_example.go) file:

## This is a basic struct with a single string field

```go
type genus struct {
    name string
}
```

## This is a more complex struct utilizing other structs and primitives

```go
type plant struct {
    genus
    species

    cultivar   string
    commonName string
    hybrid     bool
}
```

## Receiver functions

Structs can have [receiver function](https://go.dev/tour/methods/4) that will take an action on the specific object.  In this case it will return the values in a formatted string for the type "plant"

```go
func (p plant) info() string {
    return fmt.Sprintf("Genus %s, Species %s, Cultivar %s, Common name: %s, is hybrid %t\n", p.genus.name, p.species.name, p.cultivar, p.commonName, p.hybrid)
}
```

## Put it all together and this is what we get

```go
// define a single plant
p := plant{
    genus:      genus{"Solanum"},
    species:    species{"lycopersicum"},
    cultivar:   "Cherokee Purple",
    commonName: "tomato",
    hybrid:     true,
}

// output the info for the plant created above
fmt.Println(p.info())
```

Output

```bash
Genus Solanum, Species lycopersicum, Cultivar Cherokee Purple, Common name: tomato, is hybrid true
```

## Interface

An [interface](https://go.dev/tour/methods/9) is a defined set of method signatures.  So long as your type complies it can be used as the interface type.

Example from our [struct_example.go](../examples/structures/struct_example.go) file:

### This is a simple interface

```go
type item interface {
    info() string
}
```

because the plant struct (above) implements the info method that returns a string it can be used by the interface

```go
p := plant{
    genus:      genus{"Solanum"},
    species:    species{"lycopersicum"},
    cultivar:   "Cherokee Purple",
    commonName: "tomato",
    hybrid:     true,
}

itfs := []item{p}
for _, value := range itfs {
    fmt.Println(value.info())
}
```

Output

```bash
Genus Solanum, Species lycopersicum, Cultivar Cherokee Purple, Common name: tomato, is hybrid true
```
