# Logging

Logging is really just a way to keep track of things that happen in the program.  Now, there are ways, and there are better ways to do this but some of that squabble is better left to philosophers and poets.  Instead this is a somewhat practical guide as to how to log effectively.

## Best to begin at the beginning

Logging is just stuff you want a user to potentially know about.  Some of it is things we want to know out of curiosity, some of it is critical information needed for security audits.  How you handle it depends on the situation and the methods you're choosing to use.  Most of the logging that is done in these example programs will be to the [standard out](https://en.wikipedia.org/wiki/Standard_streams) of the console.  While possible to write them directly to file most of the time the practitioner of the program will redirect or tee the output to a file at runtime often because the programs are run in a [headless](https://en.wikipedia.org/wiki/Headless_software) fashion.

## Logging basics

Let's take for example we have a webservice where a user "Charlie" successfully logs in.  We want to display this to our app owners.

### A way, not a great way

We can simply print it to standard out:

```go
fmt.Printf("User %s has logged in\n", username)
```

Output:

```bash
User Charlie has logged in
```

The problem with this output is we know the who and the what, we don't know the when or the where.

### A slightly better, but still not great way to do it

We can use Go's built in [logger](https://pkg.go.dev/log)

```go
log.Printf("User %s has logged in\n", username)
```

Output:

```bash
2024/01/15 17:21:53 User Charlie has logged in
```

The problem with this output is we know the who, the what and the when, we don't the where.

### A pretty good way to do it

We can use the [logging wrapper class present in the "Grow with stl-go" project](../../pkg/log/log.go)

```go
log.Infof("User %s has logged in\n", username)
```

Output:

```bash
[stl-go] 2024/01/15 15:47:21 stl-go/go-learning-series/examples/logging/logging_example.go:65: [INFO] User Charlie has logged in
```

What you see in this example is what could be considered a complete log line.  It starts with the system generating the message, the timestamp of the message, the package and file that printed the message, the line number in that file that called the logger, what log level it was logged at and the message.  If you needed to find out why something got logged you can refer straight to the file and line and start walking the code tree based on that.

## Log Levels

Logging can be thought of as an increasing level of permission.  In our example we start with 1 (Fatal) and go to 6 (Trace), where messages that are sent at a level 1 will always be displayed, but messages at level 6 may not be.

### The Levels

1. Fatal - to be used when something happens and the program should no longer continue.  It will cause the program to exit with an os.Exit(-1).
2. Error - to be used when something happens in a way that isn't considered correct and should be flagged as such
3. Warn - to be used for when something isn't quite an error but it's more important than just an informational message
4. Info - to be used to display pertinent information that our app owners may care about
5. Debug - to be used for deeper level details needed when writing a program
6. Trace - to be used in the fine grain details we may care about when writing a program

### How Log Levels work

If for example you set your log level = 3 (warn) at the start of, or dynamically while running your program, log.Warn, log.Error and log.Fatal would be displayed; log.Trace, log.Debug and log.Info would not.

### Example

You can see this live in action in our [logging_test.go](../../pkg/log/log_test.go) example.  To run this example:

```bash
go test -v ./pkg/log/log_test.go
```

Output

```bash
=== RUN   TestLoggingTrace
=== RUN   TestLoggingTrace/TraceViewable
[stl-go] 2026/01/05 14:02:24 D:/documents/bandgeekphotos.org/stl-go/go-learning-series/pkg/log/log.go:88: [TRACE] log level set to TRACE
=== RUN   TestLoggingTrace/TracefViewable
=== RUN   TestLoggingTrace/TraceNotViewable
=== RUN   TestLoggingTrace/TracefNotViewable
--- PASS: TestLoggingTrace (0.00s)
    --- PASS: TestLoggingTrace/TraceViewable (0.00s)
    --- PASS: TestLoggingTrace/TracefViewable (0.00s)
    --- PASS: TestLoggingTrace/TraceNotViewable (0.00s)
    --- PASS: TestLoggingTrace/TracefNotViewable (0.00s)
=== RUN   TestLoggingDebug
=== RUN   TestLoggingDebug/DebugViewable
=== RUN   TestLoggingDebug/DebugfViewable
=== RUN   TestLoggingDebug/DebugNotViewable
=== RUN   TestLoggingDebug/DebugfNotViewable
--- PASS: TestLoggingDebug (0.00s)
    --- PASS: TestLoggingDebug/DebugViewable (0.00s)
    --- PASS: TestLoggingDebug/DebugfViewable (0.00s)
    --- PASS: TestLoggingDebug/DebugNotViewable (0.00s)
    --- PASS: TestLoggingDebug/DebugfNotViewable (0.00s)
=== RUN   TestLoggingInfo
=== RUN   TestLoggingInfo/InfoViewable
=== RUN   TestLoggingInfo/InfofViewable
=== RUN   TestLoggingInfo/InfoNotViewable
=== RUN   TestLoggingInfo/InfofNotViewable
--- PASS: TestLoggingInfo (0.00s)
    --- PASS: TestLoggingInfo/InfoViewable (0.00s)
    --- PASS: TestLoggingInfo/InfofViewable (0.00s)
    --- PASS: TestLoggingInfo/InfoNotViewable (0.00s)
    --- PASS: TestLoggingInfo/InfofNotViewable (0.00s)
=== RUN   TestLoggingWarn
=== RUN   TestLoggingWarn/WarnViewable
=== RUN   TestLoggingWarn/WarnfViewable
=== RUN   TestLoggingWarn/WarnNotViewable
=== RUN   TestLoggingWarn/WarnfNotViewable
--- PASS: TestLoggingWarn (0.00s)
    --- PASS: TestLoggingWarn/WarnViewable (0.00s)
    --- PASS: TestLoggingWarn/WarnfViewable (0.00s)
    --- PASS: TestLoggingWarn/WarnNotViewable (0.00s)
    --- PASS: TestLoggingWarn/WarnfNotViewable (0.00s)
=== RUN   TestLoggingError
=== RUN   TestLoggingError/ErrorViewable
=== RUN   TestLoggingError/ErrorfViewable
=== RUN   TestLoggingError/ErrorNotViewable
=== RUN   TestLoggingError/ErrorfNotViewable
--- PASS: TestLoggingError (0.00s)
    --- PASS: TestLoggingError/ErrorViewable (0.00s)
    --- PASS: TestLoggingError/ErrorfViewable (0.00s)
    --- PASS: TestLoggingError/ErrorNotViewable (0.00s)
    --- PASS: TestLoggingError/ErrorfNotViewable (0.00s)
PASS
ok      command-line-arguments  0.448s
```
