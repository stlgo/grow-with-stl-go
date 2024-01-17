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

// genus struct definition
type genus struct {
	name string
}

// species struct definition
type species struct {
	name string
}

// plant struct definition
type plant struct {
	genus
	species

	cultivar   string
	commonName string
	hybrid     bool
}

// receiver function for plant type
func (p plant) info() string {
	return fmt.Sprintf("Genus %s, Species %s, Cultivar %s, Common name: %s, is hybrid %t\n", p.genus.name, p.species.name, p.cultivar, p.commonName, p.hybrid)
}

// dog struct definition
type dog struct {
	genus
	species

	breed string
	name  string
	age   int
}

// receiver function for dog type
func (d dog) info() string {
	return fmt.Sprintf("Genus %s, Species %s, breed %s, name: %s, age %d\n", d.genus.name, d.species.name, d.breed, d.name, d.age)
}

// item interface definition
type item interface {
	info() string
}

var (
	// create and populate a map of plant structs
	plants = map[string]plant{
		"Red of Florence": {
			genus:      genus{"Allium"},
			species:    species{"cepa"},
			cultivar:   "Red of Florence",
			commonName: "Red Onion",
			hybrid:     false,
		},
		"Zapotec Jalapeno": {
			genus:      genus{"Capsicum"},
			species:    species{"annum"},
			cultivar:   "Zapotec Jalapeno",
			commonName: "jalapeno",
			hybrid:     false,
		},
		"Plum Regal": {
			genus:      genus{"Solanum"},
			species:    species{"lycopersicum"},
			cultivar:   "Plum Regal",
			commonName: "tomato",
			hybrid:     true,
		},
	}

	// create and populate a map of dog structs
	dogs = map[string]dog{
		"Charlie": {
			genus:   genus{"Canis"},
			species: species{"familiaris"},
			breed:   "English Pointer",
			name:    "Charlie",
			age:     13,
		},
		"Chief": {
			genus:   genus{"Canis"},
			species: species{"familiaris"},
			breed:   "Golden Retriever",
			name:    "Chief",
			age:     1,
		},
	}
)

func main() {
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

	// create a slice of item interface, assign it the plant and invoke the method signature defined by the interface
	itfs := []item{p}
	for _, value := range itfs {
		fmt.Println(value.info())
	}

	// create a map of items that we'll use to invoke info on later
	interfaceItems := make(map[string]item)

	// iterate through the map of plants and output their values
	fmt.Println("Printing plants")
	for key, value := range plants {
		fmt.Printf("%s is the key in plants with a value of: %s", key, value.info())
		interfaceItems[key] = value
	}

	// iterate through the map of dogs and output their values
	fmt.Println("\nPrinting dogs")
	for key, value := range dogs {
		fmt.Printf("%s is the key in dogs with a value of: %s", key, value.info())
		interfaceItems[key] = value
	}

	// iterate through the map of items and output their values
	fmt.Println("\nPrinting the interface")
	for key, value := range interfaceItems {
		fmt.Printf("%s is the key in interfaceItems with a value of: %s", key, value.info())
	}
}
