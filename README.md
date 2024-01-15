# Go Learning Series

## What is go and what can you do with it?

Go, also known as Golang, is a statically typed compiled programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson. It is syntactically similar to C, but also has memory safety, garbage collection, structural typing, and CSP-style concurrency.  Docker, Kubernetes, Terraform, etcd, ist.io and many other tools are written in go.

## Learning Go

Learning anything is a uniquely individual process.  There are some hypotheses that state somewhere between [4](https://vark-learn.com/introduction-to-vark/the-vark-modalities/) and [8](https://www.viewsonic.com/library/education/the-8-learning-styles/) different ways people learn.  Though it's often hard to understand how to best learn a new thing arming yourself with as many arrows in your quiver is the best way to get there.  Thankfully there is a wealth of online tools to help with that.  Here are some:

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

#### There are several other extensions that will come into use during this course

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

Congratulations you've now run your first Go program!
