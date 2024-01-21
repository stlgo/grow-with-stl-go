# Files

Files reads and writes are one of the fundamental things we do in programming.  Go makes this SUPER easy.

Examples from our [files.go](files_example.go) file:

## Temp Dir

If you don't have a defined directory to write things we can create an os default temp directory for our use

### Create Temp Dir function

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

Call the create temp dir function

```go
tmpDir, err := makeTempDir()
if err != nil {
    return nil, err
}
```

Output

```bash
Temp dir /temp/stl-go2895443955 was created
```

## Simple Files

Just your basic every day average non compressed regular ASCII files.

### Write a simple file

```go
func writeSimpleFile() (*string, error) {
    // create the temp dir
    tmpDir, err := makeTempDir()
    if err != nil {
        return nil, err
    }

    // write a basic text file
    if tmpDir != nil {
        fileName := filepath.Join(*tmpDir, "simpleFile.txt")
        txt := fmt.Sprintf("This text was be written to the file '%s' by this example program on %s", fileName, time.Now().Format("Mon Jan 2 15:04:05 MST 2006"))
        err := os.WriteFile(fileName, []byte(txt), 0o600)
        if err != nil {
            return nil, err
        }
        fmt.Printf("Temp file %s was created and successfully written to\n", fileName)
        return &fileName, nil
    }
    return nil, fmt.Errorf("directory is nil, cannot continue")
}
```

Call the write file function

```go
fileName, err := writeSimpleFile()
if err != nil {
    fmt.Printf("Unable to continue, cannot write a simple file: %s", err)
    os.Exit(-1)
}
```

Output

```bash
Temp dir /tmp/stl-go3632940039 was created
Temp file /tmp/stl-go3632940039/simpleFile.txt was created and successfully written to
```

Look at the file on disk

```bash
$ more stl-go2895443955/simpleFile.txt
This text was be written to the file '/tmp/stl-go3632940039/simpleFile.txt' by this example program on Sat Jan 20 13:40:08 CST 2024
```

### Read a simple file

```go
func readSimpleFile(fileName *string) (*string, error) {
    if fileName != nil {
        bytes, err := os.ReadFile(*fileName)
        if err != nil {
            return nil, err
        }
        txt := string(bytes)
        return &txt, nil
    }
    return nil, fmt.Errorf("file name is nil, cannot continue")
}
```

Call the read file function

```go
fileText, err := readSimpleFile(fileName)
if err != nil {
    fmt.Printf("Unable to continue, cannot read the file %s: %s", *fileName, err)
    os.Exit(-1)
}

if fileText != nil {
    fmt.Printf("Text from %s is as follows: \n%s\n", *fileName, *fileText)
}
```

Output

```bash
Text from /tmp/stl-go2895443955/simpleFile.txt is as follows:
This text was be written to the file '/tmp/stl-go2895443955/simpleFile.txt' by this example program on Sat Jan 20 13:40:08 CST 2024
```

## Files are good, compressed files are better

So we tend to use and store a lot of information, compression helps in both transport and storage by consuming less space.  In this case we're looking at [gzip](https://www.gzip.org/) to compress the files, it does a fairly good job of about 10:1 compression ratios on most text based documents.

### Write a compressed file

```go
func writeGzipFile() (*string, error) {
    // create the temp dir
    tmpDir, err := makeTempDir()
    if err != nil {
        return nil, err
    }

    // write a basic text file
    if tmpDir != nil {
        fileName := filepath.Join(*tmpDir, "gzipFile.txt.gz")
        txt := fmt.Sprintf("This text was be written to the file '%s' by this example program on %s", fileName, time.Now().Format("Mon Jan 2 15:04:05 MST 2006"))

        // create the gzip file
        fi, err := os.Create(fileName)
        if err != nil {
            return nil, err
        }
        defer fi.Close()

        // create the gzip file writer
        gzw := gzip.NewWriter(fi)
        defer gzw.Close()

        // create the buffered writer
        bfw := bufio.NewWriter(gzw)
        defer bfw.Flush()

        numBytes, err := bfw.WriteString(txt)
        if err != nil {
            return nil, err
        }

        fmt.Printf("%d bytes were written to %s\n", numBytes, fileName)
        return &fileName, nil
    }
    return nil, fmt.Errorf("the tmp directory is nil, cannot continue")
}
```

Call the write compressed file function

```go
fileName, err := writeGzipFile()
if err != nil {
    fmt.Printf("Unable to continue, cannot write a simple file: %s", err)
    os.Exit(-1)
}
```

Output

```bash
Temp dir /tmp/stl-go2395066881 was created
160 bytes were written to /tmp/stl-go2395066881/gzipFile.txt.gz
```

Look at the file on disk, to do this you'll need to use [zcat](https://linux.die.net/man/1/zcat) or some other form of compressed file viewer

```bash
$ zcat gzipFile.txt.gz
This text was be written to the file '/tmp/stl-go2395066881/gzipFile.txt.gz' by this example program on Sun Jan 21 09:29:42 CST 2024
```

### Read a compressed file

```go
func readGzipFile(fileName *string) (*string, error) {
    if fileName != nil {
        // crack open the file
        f, err := os.Open(*fileName)
        if err != nil {
            return nil, err
        }
        defer f.Close()

        // create a gzip file reader on the open file handler
        gzr, err := gzip.NewReader(f)
        if err != nil {
            return nil, err
        }
        defer gzr.Close()

        bytes, err := io.ReadAll(gzr)
        if err != nil {
            return nil, err
        }

        txt := string(bytes)
        return &txt, nil
    }
    return nil, fmt.Errorf("file name is nil, cannot continue")
}
```

Call the read compressed file function

```go
fileText, err := readGzipFile(fileName)
if err != nil {
    fmt.Printf("Unable to continue, cannot read the file %s: %s", *fileName, err)
    os.Exit(-1)
}

if fileText != nil {
    fmt.Printf("Text from %s is as follows: \n%s\n", *fileName, *fileText)
}
```

Output

```bash
Text from /tmp/stl-go2395066881/gzipFile.txt.gz is as follows:
This text was be written to the file '/tmp/stl-go2395066881/gzipFile.txt.gz' by this example program on Sun Jan 21 09:29:42 CST 2024
```

## Example

You can see this live in action in our [files.go](files_example.go)  example.  To run this example:

```bash
go run examples/files/files_example.go
```

Output

```bash
Temp dir /tmp/stl-go3154821507 was created
Temp file /tmp/stl-go3154821507/simpleFile.txt was created and successfully written to
Text from /tmp/stl-go3154821507/simpleFile.txt is as follows:
This text was be written to the file '/tmp/stl-go3154821507/simpleFile.txt' by this example program on Sat Jan 20 19:43:00 CST 2024
Temp dir /tmp/stl-go2788945162 was created
160 bytes were written to /tmp/stl-go2788945162/gzipFile.txt.gz
Text from /tmp/stl-go2788945162/gzipFile.txt.gz is as follows:
This text was be written to the file '/tmp/stl-go2788945162/gzipFile.txt.gz' by this example program on Sat Jan 20 19:43:00 CST 2024
```
