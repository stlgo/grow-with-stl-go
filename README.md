# Go Learning Series

## What is go and what can you do with it?

Go, also known as Golang, is a statically typed compiled programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson. It is syntactically similar to C, but also has memory safety, garbage collection, structural typing, and CSP-style concurrency.  Docker, Kubernetes, Terraform, etcd, ist.io and many other tools are written in go.

## Learning Go

Learning anything is a uniquely individual process.  There are some hypotheses that state there are somewhere between [4](https://vark-learn.com/introduction-to-vark/the-vark-modalities/) and [8](https://www.viewsonic.com/library/education/the-8-learning-styles/) different ways people learn.  Though it's often hard to understand how to best learn a new thing arming yourself with as many arrows in your quiver is the best way to get there.  Thankfully there is a wealth of online tools to help with that.  Here are some:

### For the impatient

<https://go.dev/tour/welcome/1>

### Self paced example driven learning

<https://go.dev/doc/effective_go>

<https://gobyexample.com/>

### Code Review Comments

Code Review Comments along with Effective Go [(above)](https://go.dev/doc/effective_go) are the base guidelines for "idiomatic" go.  The golint tool uses them for its warnings:\
<https://go.dev/wiki/CodeReviewComments>

## Installing go

Go can be downloaded and installed from here: <https://go.dev/dl/>

## VSCode (Optional but recommended)

[VSCode](https://code.visualstudio.com/download) is a popular free [IDE](https://en.wikipedia.org/wiki/Integrated_development_environment) that supports several languages including go.

### VSCode Extensions

#### You will need several plugins to enable go development in VSCode

- [Go for Visual Studio Code](https://marketplace.visualstudio.com/items?itemName=golang.Go)
- [Go Critic](https://marketplace.visualstudio.com/items?itemName=neverik.go-critic)
- [vscode-go-syntax](https://marketplace.visualstudio.com/items?itemName=dunstontc.vscode-go-syntax)

#### There are several other extensions that will come into use with this learning project

- [markdownlint](https://marketplace.visualstudio.com/items?itemName=DavidAnson.vscode-markdownlint)
- [Spelling Checker for Visual Studio Code](https://marketplace.visualstudio.com/items?itemName=streetsidesoftware.code-spell-checker)
- [VS Code Makefile Tools](https://marketplace.visualstudio.com/items?itemName=ms-vscode.makefile-tools)

## Write your very first go program

1. Initialize a go.mod with the command "go mod init stl-go/go-learning-series"

    ```bash
    $ go mod init stl-go/go-learning-series
    go: creating new go.mod: module stl-go/go-learning-series
    ```

2. Edit / create a main.go file:

    ```go
    package main

    import "fmt"

    func main() {
        fmt.Println("Hello world from the St. Louis Go learning Series")
    }
    ```

3. Run your program by executing "go run main.go"

    ```bash
    $ go run main.go
    Hello world from the St. Louis Go learning Series
    ```

4. You could, if you so chose, to create an executable of your hello world program by executing "go build main.go" this will create a main (or main.exe if you're on windows) binary

    ```bash
    go build main.go
    ```

5. You can now run your executable without the need to invoke go commands

    ```bash
    $ ./main
    Hello world from the St. Louis Go learning Series
    ```

6. You can build the project by running "go build" this will create a go-learning-series (or go-learning-series.exe if you're on windows) binary

    ```bash
    go build
    ```

7. You can now run your executable without the need to invoke go commands

    ```bash
    $ ./go-learning-series
    Hello world from the St. Louis Go learning Series
    ```

Congratulations you've now run your first Go program!

## Follow us further down the rabbit hole

We have some examples and some documentation on how to go from "hello world" to a fully functional webservice that can handle both [REST](https://en.wikipedia.org/wiki/REST) and [WebSockets](https://en.wikipedia.org/wiki/WebSocket), we won't say in 45 minutes or less, but with some time and elbow grease it's doable.

### The basics

1. [Primitive Types](examples/primitives/README.md)
2. [Structures](examples/structures/README.md)
3. [Maps and Slices](examples/maps_and_slices/README.md)
4. [Logging](examples/logging/README.md)
5. [Files](examples/files/README.md)
6. [JSON](examples/json/README.md)
7. [Go routines](examples/goroutines/README.md)
8. [Channels](exmaples/channels/README.go)
9. [Timed Tasks](docs/timed_tasks.md)
10. [Logging](examples/logging/README.md)
11. [Cryptography](docs/cryptography.md)

### Getting things going

1. Configurations
2. Database access
3. Static web hosting
4. [Virtual hosting](docs/vhosting.md)
5. [WebService clients](examples/clients/README.md)

### Whelp, you came this far, may as well go a bit further

Sometimes seeing things in context helps.  The next set of examples is included in the "Grow with stl-go" application example.  It is a fully functioning webservice with a UI and REST / WebSocket apis

1. Makefile
2. [Linting](docs/linting.md)

## The "Grow with stl-go" example application

More details about the [Grow with stl-go](docs/grow-with-stl-go-application.md) can be found on the [docs](docs/grow-with-stl-go-application.md).

### Prerequisites

1. [Go](https://go.dev/dl/)
2. [Make](https://www.gnu.org/software/make/)*
3. [npm](https://nodejs.org/en)
4. (optional for linting) [python](https://www.python.org/downloads/)

*For windows you can use [cygwin](https://www.cygwin.com/), note that you'll need [MinGW](https://www.mingw-w64.org/) on the path prior to cygwin to compile the Sqlite module.

## Working with the sample app

### Clone the repository

```bash
git clone https://github.com/stlgo/grow-with-stl-go.git
```

Execute a make inside the project directory

```bash
$ make
Executing frontend build steps...
npm WARN deprecated @fortawesome/fontawesome-free-solid@5.0.13: This package is deprecated. See https://git.io/fNCzJ for information about upgrading.

added 232 packages, and audited 233 packages in 29s

36 packages are looking for funding
  run `npm fund` for details

2 moderate severity vulnerabilities

Some issues need review, and may require choosing
a different dependency.

Run `npm audit` for details.
Frontend build completed successfully
Executing backend build steps...
Backend build completed successfully
```

Start the application

```bash
$ bin/grow-with-stl-go --loglevel 6
[stl-go] 2024/02/28 13:17:25 stl-go/grow-with-stl-go/pkg/configs/configs.go:118: [ERROR] open etc/grow-with-stl-go.json: The system cannot find the path specified.
[stl-go] 2024/02/28 13:17:25 stl-go/grow-with-stl-go/pkg/configs/configs.go:119: [INFO] No configuration found building a default configuration
[stl-go] 2024/02/28 13:17:25 stl-go/grow-with-stl-go/pkg/configs/configs.go:126: [ERROR] unexpected end of JSON input
[stl-go] 2024/02/28 13:17:25 stl-go/grow-with-stl-go/pkg/configs/apiuser.go:54: [WARN] Password generated for user 'admin', password 5363839e28768ba89b5fef2c503c4885e21f7fe1bde18963351a965f3e3706b8 - DO NOT USE THIS FOR PRODUCTION
[stl-go] 2024/02/28 13:17:25 stl-go/grow-with-stl-go/pkg/configs/apiuser.go:54: [WARN] Password generated for user 'user', password d5f2bf6124cc4a35294b831266a1cf126fc82fa03c778c68378461650c5d59c5 - DO NOT USE THIS FOR PRODUCTION
[stl-go] 2024/02/28 13:17:25 stl-go/grow-with-stl-go/pkg/configs/webservice.go:35: [INFO] No webservice config found, generating ssl keys, host and port info
[stl-go] 2024/02/28 13:17:25 stl-go/grow-with-stl-go/pkg/cryptography/ssl.go:40: [WARN] Generating private key C:\temp\grow-with-stl-go\etc\key.pem.  DO NOT USE THIS FOR PRODUCTION
[stl-go] 2024/02/28 13:17:25 stl-go/grow-with-stl-go/pkg/cryptography/ssl.go:47: [WARN] Generating public key C:\temp\grow-with-stl-go\etc\cert.pem.  DO NOT USE THIS FOR PRODUCTION
[stl-go] 2024/02/28 13:17:25 stl-go/grow-with-stl-go/pkg/configs/SQLite.go:57: [DEBUG] No embedded database found in the config file, generating a default configuration
[stl-go] 2024/02/28 13:17:25 stl-go/grow-with-stl-go/pkg/configs/SQLite.go:88: [DEBUG] Using encrypted aud database in C:\temp\grow-with-stl-go\etc\grow-with-stl-go.db
[stl-go] 2024/02/28 13:17:25 stl-go/grow-with-stl-go/pkg/configs/configs.go:282: [DEBUG] Rewriting etc/grow-with-stl-go.json to ensure data is enciphered on disk
[stl-go] 2024/02/28 13:17:25 stl-go/grow-with-stl-go/pkg/configs/SQLite.go:132: [TRACE] Audit table WebSocket was created if it didn't already exist
[stl-go] 2024/02/28 13:17:26 stl-go/grow-with-stl-go/pkg/configs/SQLite.go:132: [TRACE] Audit table REST was created if it didn't already exist
[stl-go] 2024/02/28 13:17:26 stl-go/grow-with-stl-go/pkg/configs/SQLite.go:132: [TRACE] Audit table user was created if it didn't already exist
[stl-go] 2024/02/28 13:17:26 stl-go/grow-with-stl-go/pkg/webservice/websocket.go:74: [DEBUG] Regestering 'seeds' as a WebSocket Endpoint
[stl-go] 2024/02/28 13:17:26 stl-go/grow-with-stl-go/pkg/webservice/webservice.go:53: [DEBUG] Regestering seeds as a REST Endpoint
[stl-go] 2024/02/28 13:17:26 stl-go/grow-with-stl-go/pkg/configs/SQLite.go:132: [TRACE] Audit table seeds was created if it didn't already exist
[stl-go] 2024/02/28 13:17:26 stl-go/grow-with-stl-go/pkg/configs/SQLite.go:156: [TRACE] Inserting default chive for table seeds if it doesn't already exist
[stl-go] 2024/02/28 13:17:26 stl-go/grow-with-stl-go/pkg/configs/SQLite.go:156: [TRACE] Inserting default bell pepper for table seeds if it doesn't already exist
[stl-go] 2024/02/28 13:17:26 stl-go/grow-with-stl-go/pkg/configs/SQLite.go:156: [TRACE] Inserting default serrano for table seeds if it doesn't already exist
[stl-go] 2024/02/28 13:17:26 stl-go/grow-with-stl-go/pkg/configs/SQLite.go:156: [TRACE] Inserting default galahad for table seeds if it doesn't already exist
[stl-go] 2024/02/28 13:17:26 stl-go/grow-with-stl-go/pkg/configs/SQLite.go:156: [TRACE] Inserting default plum regal for table seeds if it doesn't already exist
[stl-go] 2024/02/28 13:17:26 stl-go/grow-with-stl-go/pkg/configs/SQLite.go:156: [TRACE] Inserting default oreagno for table seeds if it doesn't already exist
[stl-go] 2024/02/28 13:17:26 stl-go/grow-with-stl-go/pkg/configs/SQLite.go:156: [TRACE] Inserting default ailsa craig for table seeds if it doesn't already exist
[stl-go] 2024/02/28 13:17:26 stl-go/grow-with-stl-go/pkg/configs/SQLite.go:156: [TRACE] Inserting default walla walla for table seeds if it doesn't already exist
[stl-go] 2024/02/28 13:17:26 stl-go/grow-with-stl-go/pkg/configs/SQLite.go:156: [TRACE] Inserting default carbon for table seeds if it doesn't already exist
[stl-go] 2024/02/28 13:17:26 stl-go/grow-with-stl-go/pkg/configs/SQLite.go:156: [TRACE] Inserting default basil for table seeds if it doesn't already exist
[stl-go] 2024/02/28 13:17:26 stl-go/grow-with-stl-go/pkg/configs/SQLite.go:156: [TRACE] Inserting default poblano for table seeds if it doesn't already exist
[stl-go] 2024/02/28 13:17:26 stl-go/grow-with-stl-go/pkg/configs/SQLite.go:156: [TRACE] Inserting default san marzano for table seeds if it doesn't already exist
[stl-go] 2024/02/28 13:17:26 stl-go/grow-with-stl-go/pkg/configs/SQLite.go:156: [TRACE] Inserting default dill for table seeds if it doesn't already exist
[stl-go] 2024/02/28 13:17:26 stl-go/grow-with-stl-go/pkg/configs/SQLite.go:156: [TRACE] Inserting default patterson for table seeds if it doesn't already exist
[stl-go] 2024/02/28 13:17:26 stl-go/grow-with-stl-go/pkg/configs/SQLite.go:156: [TRACE] Inserting default red wing for table seeds if it doesn't already exist
[stl-go] 2024/02/28 13:17:26 stl-go/grow-with-stl-go/pkg/configs/SQLite.go:156: [TRACE] Inserting default jalapeno for table seeds if it doesn't already exist
[stl-go] 2024/02/28 13:17:26 stl-go/grow-with-stl-go/pkg/configs/configs.go:282: [DEBUG] Rewriting etc/grow-with-stl-go.json to ensure data is enciphered on disk
[stl-go] 2024/02/28 13:17:26 stl-go/grow-with-stl-go/pkg/seeds/seeds.go:160: [DEBUG] 4 categories in inventory cache
[stl-go] 2024/02/28 13:17:26 stl-go/grow-with-stl-go/pkg/webservice/websocket.go:74: [DEBUG] Regestering 'admin' as a WebSocket Endpoint
[stl-go] 2024/02/28 13:17:26 stl-go/grow-with-stl-go/pkg/webservice/server.go:57: [DEBUG] Attempting to serve static content from web
[stl-go] 2024/02/28 13:17:26 stl-go/grow-with-stl-go/pkg/webservice/server.go:63: [INFO] Attempting to start webservice on localhost:10443
```

Login to the web page with the generated ids / passwords in the output above at <localhost:10443> or flex the REST and WebSocket apis via the procedures in our [documentation](docs/grow-with-stl-go-application.md).

## Special considerations for programming

Programming is both an art and a science and that concept is sometimes lost in translation.  If you've read Isaac Asimov's [The Relativity of Wrong](https://www.sas.upenn.edu/~dbalmer/eportfolio/Nature%20of%20Science_Asimov.pdf) you'll understand that we may not always get everything 100% perfectly right every time but so long as we continue to move towards a more correct and better way to do things we'll improve.  To that end Tim Peters wrote a small blurb called the [Zen of Python](https://en.wikipedia.org/wiki/Zen_of_Python), and while it does have the word python in the title, the concepts can be applied to this or any project you may be involved in and should be considered in how you approach problem solving:

```text
    Beautiful is better than ugly.
    Explicit is better than implicit.
    Simple is better than complex.
    Complex is better than complicated.
    Flat is better than nested.
    Sparse is better than dense.
    Readability counts.
    Special cases aren't special enough to break the rules.
    Although practicality beats purity.
    Errors should never pass silently.
    Unless explicitly silenced.
    In the face of ambiguity, refuse the temptation to guess.
    There should be one-- and preferably only one --obvious way to do it.
    Although that way may not be obvious at first unless you're Dutch.
    Now is better than never.
    Although never is often better than right now.
    If the implementation is hard to explain, it's a bad idea.
    If the implementation is easy to explain, it may be a good idea.
    Namespaces are one honking great idea – let's do more of those!
```
