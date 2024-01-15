# Code Linting

## What is code linting?

Linting is the automated checking of your source code for programmatic and stylistic errors. This is done by using a lint tool (otherwise known as linter). A lint tool is a basic static code analyzer.  So let's face it, we all make mistakes.  Some mistakes are small whitespace or spelling errors that are the source of some minor embarrassment and some are large logic errors that cost millions or worse.  We have tools available to us to help us avoid **some** (not all obviously) errors.

For go, and the examples in the project we use these (and maybe a few more):

- [go-fmt](https://pkg.go.dev/fmt)
- [golangci-lint](https://github.com/golangci/golangci-lint)
- [go-critic](https://github.com/go-critic/go-critic)
- [go-unit-tests](https://go.dev/doc/tutorial/add-a-test)

## pre-commit

[pre-commit](https://pre-commit.com/) can be used as a git hook to test your code prior to committal.  This allows us to execute our linters, [cyclomatic complexity tests](https://en.wikipedia.org/wiki/Cyclomatic_complexity) style guides, automated test and whatever we can think of automatically prior to committing code to the repository.  Working with pre-commit is fairly easy but it does require python to be available on the system

### Installing pre-commit

1. Issue the pip command to install: "pip install pre-commit"

   ```bash
    $ pip install pre-commit
    Successfully installed cfgv-3.4.0 distlib-0.3.8 filelock-3.13.1 identify-2.5.33 nodeenv-1.8.0 platformdirs-4.1.0 pre-commit-3.6.0 pyyaml-6.0.1 setuptools-69.0.3 virtualenv-20.25.0
    ```

    **NOTE:**
    Python will need do be on your path for this to work, if it isn't you may need to adjust your environment variables and restart your shell / IDE

2. Test that pre-commit is avaialble on the command line by issuing: "pre-commit --version"

   ```bash
    $ pre-commit --version
    pre-commit 3.6.0
    ```

3. Create [(or use the example in this project)](../../.pre-commit-config.yaml) a .pre-commit-config.yaml file at the root of your project
4. Test the pre-commit by issuing a "git add ." followed by "pre-commit"

   ```bash
   $ git add .
   $ pre-commit
    trim trailing whitespace.................................................Passed
    fix end of files.........................................................Passed
    check yaml...............................................................Passed
    check for added large files..............................................Passed
    go-cyclo.............................................(no files to check)Skipped
    validate toml........................................(no files to check)Skipped
    Check files aren't using go's testing package........(no files to check)Skipped
    go fmt...............................................(no files to check)Skipped
    go imports...........................................(no files to check)Skipped
    golangci-lint........................................(no files to check)Skipped
    go-critic............................................(no files to check)Skipped
    go-unit-tests........................................(no files to check)Skipped
    go-build.............................................(no files to check)Skipped
    go-mod-tidy..........................................(no files to check)Skipped
    PS D:\documents\bandgeekphotos.org\stl-go\go-learning-series>
   ```

    **NOTE:**
    There may be some installations that occur on the first run

5. The pre-commit git hook will run the pre-commit for you prior to your commits.  This is very helpful in keeping the repo clean, and most good projects will run your code through a linter prior to committal anyway.  Install the git hook script by issuing: "pre-commit install"

   ```bash
   $ re-commit install
   pre-commit installed at .git\hooks\pre-commit
   ```

6. Test your pre-commit hook by creating a git commit:

   ```bash
   $ git add .
   $ git commit -am "pre-commit added to the project"
    trim trailing whitespace.................................................Passed
    fix end of files.........................................................Passed
    check yaml...............................................................Passed
    check for added large files..............................................Passed
    go-cyclo.............................................(no files to check)Skipped
    validate toml........................................(no files to check)Skipped
    Check files aren't using go's testing package........(no files to check)Skipped
    go fmt...............................................(no files to check)Skipped
    go imports...........................................(no files to check)Skipped
    golangci-lint........................................(no files to check)Skipped
    go-critic............................................(no files to check)Skipped
    go-unit-tests........................................(no files to check)Skipped
    go-build.............................................(no files to check)Skipped
    go-mod-tidy..........................................(no files to check)Skipped
    [master 2eed9f1] pre-commit added to the project
    5 files changed, 394 insertions(+), 2 deletions(-)
    create mode 100755 .golangci.yaml
    create mode 100755 .pre-commit-config.yaml
    create mode 100755 examples/linting/README.md
   ```

   **NOTE:**
   You may need to adjust the .git\hooks\pre-commit file if there are any pathing or newline (unix -vs- windows line breaks) irregularities on windows
