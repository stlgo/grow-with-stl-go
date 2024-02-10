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

## SHA256 Sum

An [SHA256 sum](https://csrc.nist.gov/pubs/fips/180-4/upd1/final) is a way to make sure that file you're working with is complete and / or tampered with.  While it is not every day that you will do summing of files, it is however best practice to do so when transferring files between systems.

Take for example a sum calculated by the sha256 command:

```bash
$ sha256sum gzipFile.txt.gz
dae819de5d7d04a956cdd4bdcd5fe05405001eea5390ed53e9fb919050e6e3bb *gzipFile.txt.gz
```

Compare this to the sha sum as calculated by the program and written to the .sha256 file:

```bash
$ more gzipFile.txt.gz.sha256
dae819de5d7d04a956cdd4bdcd5fe05405001eea5390ed53e9fb919050e6e3bb
```

If the strings were different you could tell that the file is not complete and / or was tampered with and armed with that knowledge you can take the appropriate action for your situation.

### Create the SHA256 sum of a file

```go
func getSHA256Sum(fileName *string) (*string, error) {
    if fileName != nil {
        f, err := os.Open(*fileName)
        if err != nil {
            return nil, err
        }
        defer f.Close()

        hash := sha256.New()
        if _, err := io.Copy(hash, f); err != nil {
            return nil, err
        }

        sum := fmt.Sprintf("%x", hash.Sum(nil))
        return &sum, nil
    }
    return nil, errors.New("nil filename cannot calculate sha256 sum")
}
```

### Write the SHA256 to a .sha256 file next to the original

```go
func writeSHA256Sum(fileName *string) error {
    if fileName != nil {
        shaSum, err := getSHA256Sum(fileName)
        if err != nil || shaSum == nil {
            fmt.Printf("Unable to continue, cannot sha sum of file: %s", err)
            os.Exit(-1)
        }

        fmt.Printf("SHA256 sum of %s is %s\n", *fileName, *shaSum)
        sha256SumFileName := fmt.Sprintf("%s.sha256", *fileName)
        err = os.WriteFile(sha256SumFileName, []byte(*shaSum), 0o600)
        if err != nil {
            return err
        }
        fmt.Printf("SHA256 Sum of temp file %s was created and successfully written to %s\n", *fileName, sha256SumFileName)
        return nil
    }
    return errors.New("nil filename cannot write sha256 sum")
}
```

### Compare the SHA256 sum on disk to that of the file you're looking at

```go
func compareSHA256Sum(fileName *string) error {
    if fileName != nil {
        shaSum, err := getSHA256Sum(fileName)
        if err != nil || shaSum == nil {
            fmt.Printf("Unable to continue, cannot sha sum of file: %s", err)
            os.Exit(-1)
        }

        sumFileName := fmt.Sprintf("%s.sha256", *fileName)
        bytes, err := os.ReadFile(sumFileName)
        if err != nil {
            return err
        }
        shaSumFromFile := string(bytes)
        if *shaSum != shaSumFromFile {
            return fmt.Errorf("file %s has a different sha256 hash %s from the stored hash %s", *fileName, *shaSum, shaSumFromFile)
        }

        fmt.Printf("File %s hash is the same as the one stored in %s\n", *fileName, sumFileName)
        return nil
    }
    return errors.New("nil filename cannot compare sha256 sums")
}
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
$ more /tmp/stl-go2895443955/simpleFile.txt
This text was be written to the file '/tmp/stl-go3632940039/simpleFile.txt' by this example program on Sat Jan 20 13:40:08 CST 2024
```

Create the sha256 sum of the file

```go
if err = writeSHA256Sum(fileName); err != nil {
    fmt.Printf("Unable to continue, cannot write a the sum of a simple file: %s", err)
    os.Exit(-1)
}
```

Look at the sum file on disk

```bash
$ more /tmp/stl-go2895443955/simpleFile.txt.sha256
9fe0fda743a6fafeb22033bf0ad409714ee384a91a8da6c688670c7a3af38e68
```

### Read a simple file

```go
func readSimpleFile(fileName *string) (*string, error) {
    if fileName != nil {
        // sha256 hash compare first
        if err := compareSHA256Sum(fileName); err != nil {
            return nil, err
        }

        // read the file
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
File /tmp/stl-go2895443955/simpleFile.txt hash is the same as the one stored in /tmp/stl-go2895443955/simpleFile.txt.sha256
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

Create the sha256 sum of the file

```go
if err = writeSHA256Sum(fileName); err != nil {
    fmt.Printf("Unable to continue, cannot write a the sum of a simple file: %s", err)
    os.Exit(-1)
}
```

Look at the sum file on disk

```bash
$ more /tmp/stl-go2395066881/simpleFile.txt.gz.sha256
dae819de5d7d04a956cdd4bdcd5fe05405001eea5390ed53e9fb919050e6e3bb
```

### Read a compressed file

```go
func readGzipFile(fileName *string) (*string, error) {
    if fileName != nil {
        // sha256 hash compare first
        if err := compareSHA256Sum(fileName); err != nil {
            return nil, err
        }

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
File /tmp/stl-go2395066881/gzipFile.txt.gz hash is the same as the one stored in /tmp/stl-go2395066881/gzipFile.txt.gz.sha256
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
Temp dir /tmp/stl-go2302759081 was created
Temp file /tmp/stl-go2302759081/simpleFile.txt was created and successfully written to
SHA256 sum of /tmp/stl-go2302759081/simpleFile.txt is 9fe0fda743a6fafeb22033bf0ad409714ee384a91a8da6c688670c7a3af38e68
SHA256 Sum of temp file /tmp/stl-go2302759081/\simpleFile.txt was created and successfully written to /tmp/stl-go2302759081/simpleFile.txt.sha256
File /tmp/stl-go2302759081/simpleFile.txt hash is the same as the one stored in /tmp/stl-go2302759081/simpleFile.txt.sha256
Text from /tmp/stl-go2302759081/simpleFile.txt is as follows:
This text was be written to the file '/tmp/stl-go2302759081/simpleFile.txt' by this example program on Sat Feb 10 16:35:20 CST 2024
Temp dir /tmp/stl-go2280262460 was created
160 bytes were written to /tmp/stl-go2280262460/gzipFile.txt.gz
File /tmp/stl-go2280262460/gzipFile.txt.gz hash is the same as the one stored in /tmp/stl-go2280262460/gzipFile.txt.gz.sha256
Text from /tmp/stl-go2280262460/gzipFile.txt.gz is as follows:
This text was be written to the file '/tmp/stl-go2280262460/gzipFile.txt.gz' by this example program on Sat Feb 10 16:35:20 CST 2024
```
