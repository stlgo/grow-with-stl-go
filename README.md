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

3. Run your program by executing "go run hello-world"

    ```bash
    $ go run hello-world.go
    Hello world from the St. Louis Go learning Series
    ```

4. You could, if you so chose, to create an executable of your hello world program by executing "go build main.go" this will create a main (or main.exe if you're on windows) binary

    ```bash
    go build main.go
    ```

5. You can now run your executable without the need to invoke go commands

    ```bash
    $ main
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
7. [Goroutines](examples/goroutines/README.md)
8. [Channels](exmaples/channels/README.go)

### Whelp, you came this far, may as well go a bit further

Sometimes seeing things in context helps.  The next set of examples is included in the "Grow with stl-go" application example.  It is a fully functioning webservice with a UI and REST / WebSocket apis

1. [Logging](examples/logging/README.md)
2. [Linting](docs/linting.md)

## The "Grow with stl-go" example application

More details about the [Grow with stl-go](docs/grow-with-stl-go-application.md) can be found on the [docs](docs/grow-with-stl-go-application.md).

### Prerequisites

## Special considerations for programming

Programming is both an art and a science and that concept is sometimes lost in translation.  If you've read Isaac Asimov's [The Relativity of Wrong](https://www.sas.upenn.edu/~dbalmer/eportfolio/Nature%20of%20Science_Asimov.pdf) you'll understand what he put into words is that we may not always get everything 100% perfectly right every time but so long as we continue to move towards a more correct and better way to do things we'll improve.  To that end Tim Peters wrote a small blurb called the [Zen of Python](https://en.wikipedia.org/wiki/Zen_of_Python), and while it does have the word python in the title, the concepts can be applied to this or any project you may be involved in and should be considered in how you approach problem solving:

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
