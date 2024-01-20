# JSON

[JSON](https://www.json.org/json-en.html) (JavaScript Object Notation) is a lightweight data-interchange format.  JSON is ubiquitous.  It's used for interprocess communication, application configurations, data storage in database and on disk.  It has largely supplanted the use of XML in web services and some might say it's a thing we as an industry what passes for correct.

JSON at its core is key value pairs ([tuples](https://en.wikipedia.org/wiki/Tuple)).  The keys are generally strings or integers and values can be primitives or complex styles.  You can also have arrays of JSON Objects and arrays of primitives.

Examples from our [json_example.go](json_example.go) file:

## Simple JSON

### Create temp dir

The makeTempDir function is referenced several places in [this example](json_example.go).  It is used to create an OS defined temp directory to read and write files to.

```go
func makeTempDir() (*string, error) {
    tmpDir, err := os.MkdirTemp("", "stl-go")
    if err != nil {
        return nil, err
    }
    fmt.Printf("Temp dir %s was created\n", tmpDir)
    return &tmpDir, nil
}
```

### Create a simple JSON object

```go
jo := map[string]any{
    "fileText":        "This text was be written to the file by this example program for simple JSON",
    "fileDate":        time.Now().UnixMilli(),
    "fileDateISO8601": time.Now().UTC().Format("2006-01-02T15:04:05-0700"),
    "someArray":       []int{1, 2, 3, 4},
    "nestedMap": map[string]any{
        "foo":       "bar",
        "someArray": []float32{1.2, 2.4, 3.6, 4.8},
    },
}

fmt.Printf("Simple JSON created: %v\n", jo)
```

Output

```bash
Simple JSON created: map[fileDate:1705782118924 fileDateISO8601:2024-01-20T20:21:58+0000 fileText:This text was be written to the file by this example program for simple JSON nestedMap:map[foo:bar someArray:[1.2 2.4 3.6 4.8]] someArray:[1 2 3 4]]
```

### Write the simple JSON object to file

```go
func createSimpleJSON() (*string, error) {
    jo := map[string]any{
        "fileText":        "This text was be written to the file by this example program for simple JSON",
        "fileDate":        time.Now().UnixMilli(),
        "fileDateISO8601": time.Now().UTC().Format("2006-01-02T15:04:05-0700"),
        "someArray":       []int{1, 2, 3, 4},
        "nestedMap": map[string]any{
            "foo":       "bar",
            "someArray": []float32{1.2, 2.4, 3.6, 4.8},
        },
    }

    fmt.Printf("Simple JSON created: %v\n", jo)

    // create the temp dir
    tmpDir, err := makeTempDir()
    if err != nil {
        return nil, err
    }

    // marshall the StructJSON object to a byte array
    jsonBytes, err := json.MarshalIndent(jo, "", "\t")
    if err != nil {
        return nil, err
    }

    // write the byte array to the file
    fileName := filepath.Join(*tmpDir, "simple.json")
    err = os.WriteFile(fileName, jsonBytes, 0o600)
    if err != nil {
        return nil, err
    }
    fmt.Printf("Temp JSON file %s was created and successfully written to\n", fileName)
    return &fileName, nil
}
```

Call the create simple JSON function

```go
fileName, err := createSimpleJSON()
if err != nil {
    fmt.Printf("Unable to continue, cannot write a simple JSON file: %s", err)
    os.Exit(-1)
}
```

Output

```bash
Simple JSON created: map[fileDate:1705782118924 fileDateISO8601:2024-01-20T20:21:58+0000 fileText:This text was be written to the file by this example program for simple JSON nestedMap:map[foo:bar someArray:[1.2 2.4 3.6 4.8]] someArray:[1 2 3 4]]
Temp dir /tmp/stl-go1095867815 was created
Temp JSON file /tmp/stl-go1095867815simple.json was created and successfully written to
{
        "fileDate": 1705782118924,
        "fileDateISO8601": "2024-01-20T20:21:58+0000",
        "fileText": "This text was be written to the file by this example program for simple JSON",
        "nestedMap": {
                "foo": "bar",
                "someArray": [
                        1.2,
                        2.4,
                        3.6,
                        4.8
                ]
        },
        "someArray": [
                1,
                2,
                3,
                4
        ]
}
```

### Read the simple JSON from a file and process it as a generic map

```go
func readSimpleJSONFile(fileName *string) (map[string]any, error) {
    if fileName != nil {
        jsonBytes, err := os.ReadFile(*fileName)
        if err != nil {
            return nil, err
        }

        // unmarshal the file into a basic JSON Object
        var jo map[string]any
        if err1 := json.Unmarshal(jsonBytes, &jo); err1 != nil {
            return nil, err1
        }

        // print it back out as a generic JSON
        jsonOutBytes, err := json.MarshalIndent(jo, "", "\t")
        if err != nil {
            return nil, err
        }
        fmt.Println(string(jsonOutBytes))

        return jo, nil
    }
    return nil, fmt.Errorf("file name is nil, cannot continue")
}
```

Call the read simple JSON file function and interact with the map that is returned

```go
simpleJSON, err := readSimpleJSONFile(fileName)
if err != nil {
    fmt.Printf("Unable to continue, cannot read the simple json file %s: %s", *fileName, err)
    os.Exit(-1)
}

// you can now interact directly with the simple JSON
for key, value := range simpleJSON {
    fmt.Printf("Simple JSON key %s has a value of %v\n", key, value)
}

// you can also interact with specific keys in the map
if value, ok := simpleJSON["fileDateISO8601"]; ok {
    fmt.Printf("Simple JSON value of key \"fileDateISO8601\" %s\n", value)
}
```

Output

```go
Simple JSON key fileText has a value of This text was be written to the file by this example program for simple JSON
Simple JSON key nestedMap has a value of map[foo:bar someArray:[1.2 2.4 3.6 4.8]]
Simple JSON key someArray has a value of [1 2 3 4]
Simple JSON key fileDate has a value of 1.705782118924e+12
Simple JSON key fileDateISO8601 has a value of 2024-01-20T20:21:58+0000
Simple JSON value of key "fileDateISO8601" 2024-01-20T20:21:58+0000
```

## JSON utilizing a struct for marshalling / unmarshalling

One of the very nice thing Go does is by [tagging](https://www.practical-go-lessons.com/post/how-to-add-and-read-go-struct-tags-cbt2mue6togs70jopvi0) a struct object you can use it to marshall / unmarshall JSON to struts easily and to access the values without recursion.

### The tagged struct definition

Notice the 'json:"tag"' appended to the structure

```go
type StructJSON struct {
    FileDate        *int64          `json:"fileDate,omitempty"`
    FileDateISO8601 *string         `json:"fileDateISO8601,omitempty"`
    FileText        *string         `json:"fileText,omitempty"`
    SomeArray       *[]int          `json:"someArray,omitempty"`
    NestedMap       *map[string]any `json:"nestedMap,omitempty"`
}
```

Receiver function that will write the StructJSON object to a file

```go
func (jo StructJSON) persist() (*string, error) {
    // create the temp dir
    tmpDir, err := makeTempDir()
    if err != nil {
        return nil, err
    }

    // marshall the StructJSON object to a byte array
    jsonBytes, err := json.MarshalIndent(jo, "", "\t")
    if err != nil {
        return nil, err
    }

    // write the byte array to the file
    fileName := filepath.Join(*tmpDir, "structBased.json")
    err = os.WriteFile(fileName, jsonBytes, 0o600)
    if err != nil {
        return nil, err
    }
    fmt.Printf("Temp JSON file %s was created and successfully written to\n", fileName)
    return &fileName, nil
}
```

### Create a struct based JSON Object

Create a StuctJSON object and write it as JSON to disk using the struct's receiver function

```go
func createStructJSON() (*string, error) {
    fileText := "This text was be written to the file by this example program for struct based JSON"

    // because we're using pointers in our struct we need to create the variables first
    now := time.Now()
    millis := now.UnixMilli()
    iso8601 := now.UTC().Format("2006-01-02T15:04:05-0700")
    someArray := []int{1, 2, 3, 4}
    nesteMap := map[string]any{
        "foo":       "bar",
        "someArray": []float32{1.2, 2.4, 3.6, 4.8},
    }

    // we use the addresses when creating the object
    jo := StructJSON{
        FileText:        &fileText,
        FileDate:        &millis,
        FileDateISO8601: &iso8601,
        SomeArray:       &someArray,
        NestedMap:       &nesteMap,
    }

    // write the file out
    fileName, err := jo.persist()
    if err != nil {
        return nil, err
    }

    return fileName, nil
}
```

Call the create struct JSON function

```go
fileName, err := createStructJSON()
if err != nil {
    fmt.Printf("Unable to continue, cannot write a struct based JSON file: %s", err)
    os.Exit(-1)
}
```

Output

```bash
Temp dir /tmp/stl-go2316460679 was created
Temp JSON file /tmp/stl-go2316460679/structBased.json was created and successfully written to
```

### Read a file containing JSON and unmarshall it to a StructJSON object

```go
func readStructJSONFile(fileName *string) (*StructJSON, error) {
    if fileName != nil {
        jsonBytes, err := os.ReadFile(*fileName)
        if err != nil {
            return nil, err
        }

        // output the string we read in from the file
        fmt.Printf("Data read from %s is:\n%s\n", *fileName, string(jsonBytes))

        // unmarshal the json as the struct
        var structJSON StructJSON
        if err := json.Unmarshal(jsonBytes, &structJSON); err != nil {
            return nil, err
        }
        return &structJSON, nil
    }
    return nil, fmt.Errorf("file name is nil, cannot continue")
}
```

Call the read file function

```go
structJSON, err := readStructJSONFile(fileName)
if err != nil {
    fmt.Printf("Unable to continue, cannot read the struct based json file %s: %s", *fileName, err)
    os.Exit(-1)
}

// you can now interact directly with the struct
fmt.Printf("%s text %s\n, milliseconds %d, which is easier to use but harder to read than ISO8601 %s",
    *fileName, *structJSON.FileText, *structJSON.FileDate, *structJSON.FileDateISO8601)
```

Output

```bash
Data read from /tmp/stl-go3920536142/structBased.json is:
{
        "fileDate": 1705785741063,
        "fileDateISO8601": "2024-01-20T21:22:21+0000",
        "fileText": "This text was be written to the file by this example program for struct based JSON",
        "someArray": [
                1,
                2,
                3,
                4
        ],
        "nestedMap": {
                "foo": "bar",
                "someArray": [
                        1.2,
                        2.4,
                        3.6,
                        4.8
                ]
        }
}
/tmp/stl-go3920536142/structBased.json text This text was be written to the file by this example program for struct based JSON
, milliseconds 1705785741063, which is easier to use but harder to read than ISO8601 2024-01-20T21:22:21+0000
```

## Example

You can see this live in action in our [json_example.go](json_example.go) example.  To run this example:

```bash
go run examples/json/json_example.go
```

Output

```bash
Simple JSON created: map[fileDate:1705785741032 fileDateISO8601:2024-01-20T21:22:21+0000 fileText:This text was be written to the file by this example program for simple JSON nestedMap:map[foo:bar someArray:[1.2 2.4 3.6 4.8]] someArray:[1 2 3 4]]
Temp dir /tmp/stl-go3976499485/ was created
Temp JSON file /tmp/stl-go3976499485/simple.json was created and successfully written to
{
        "fileDate": 1705785741032,
        "fileDateISO8601": "2024-01-20T21:22:21+0000",
        "fileText": "This text was be written to the file by this example program for simple JSON",
        "nestedMap": {
                "foo": "bar",
                "someArray": [
                        1.2,
                        2.4,
                        3.6,
                        4.8
                ]
        },
        "someArray": [
                1,
                2,
                3,
                4
        ]
}
Simple JSON key fileDate has a value of 1.705785741032e+12
Simple JSON key fileDateISO8601 has a value of 2024-01-20T21:22:21+0000
Simple JSON key fileText has a value of This text was be written to the file by this example program for simple JSON
Simple JSON key nestedMap has a value of map[foo:bar someArray:[1.2 2.4 3.6 4.8]]
Simple JSON key someArray has a value of [1 2 3 4]
Simple JSON value of key "fileDateISO8601" 2024-01-20T21:22:21+0000
Temp dir C:\Users\root\AppData\Local\Temp\stl-go3920536142 was created
Temp JSON file /tmp/stl-go3920536142/structBased.json was created and successfully written to
Data read from /tmp/stl-go3920536142/structBased.json is:
{
        "fileDate": 1705785741063,
        "fileDateISO8601": "2024-01-20T21:22:21+0000",
        "fileText": "This text was be written to the file by this example program for struct based JSON",
        "someArray": [
                1,
                2,
                3,
                4
        ],
        "nestedMap": {
                "foo": "bar",
                "someArray": [
                        1.2,
                        2.4,
                        3.6,
                        4.8
                ]
        }
}
/tmp/stl-go3920536142/structBased.json text This text was be written to the file by this example program for struct based JSON
, milliseconds 1705785741063, which is easier to use but harder to read than ISO8601 2024-01-20T21:22:21+0000
```
