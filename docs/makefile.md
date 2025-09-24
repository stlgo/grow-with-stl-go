# Makefile

## What is a makefile?

Makefile is a defined set of tasks that is executed using [make](https://www.gnu.org/software/make/).  This allows anyone to build the application without knowing in depth each applicable command to do it.  In the context of this application the makefile defines the following:

- [build](#build): causes make to execute the steps for a frontend and a backend build
- [lint](#lint): causes make to execute the steps for a frontend and backend code lint
- [test](#test): causes make to execute the steps for backend unit tests
- [coverage](#coverage): causes make to execute steps for unit test coverage
- [verify](#verify): causes make to execute the steps to execute build, lint and coverage
- [dist](#dist): causes make to execute the steps to build a binary distribution

This information is sourced from the [Makefile](Makefile) inside the project

## Build

Make build will cause a general binary build command for the backend and frontend

### Example Make Build

To run a build command run: make build

```bash
$ make build
Executing frontend build steps...

up to date, audited 112 packages in 741ms

25 packages are looking for funding
  run `npm fund` for details

found 0 vulnerabilities
Frontend build completed successfully
Executing backend build steps...
Backend build completed successfully
```

### Makefile Build Definitions

#### Build Makefile Commands

```makefile
.PHONY: build
build: frontend-build
build: backend-build
```

#### Frontend Build

```makefile
.PHONY: frontend-build
frontend-build:
 @echo "Executing frontend build steps..."
 @cd $(WEB_DIR) && npm install && cd ..
 @echo "Frontend build completed successfully"
```

#### Backend Build

```makefile
.PHONY: backend-build
backend-build:
 @echo "Executing backend build steps..."
 @mkdir -p $(BUILD_DIR)
 @go build -o $(MAIN)$(EXTENSION) $(GO_FLAGS) cmd/main.go
 @echo "Backend build completed successfully"
```

## Lint

Make lint will cause linters to be run on both the front end and back end

### Example Make Lint

To run a lint command run: make lint

```bash
$ make lint
Checking that go.mod is up to date...
go.mod check completed successfully
Running frontend linting step...

/go-learning-series/web/grow-with-stlgo-admin/js/seeds.js
   56:35  warning  Missing radix parameter                 radix
   66:58  warning  Identifier name 'c' is too short (< 2)  id-length
   72:33  warning  Missing radix parameter                 radix
   72:57  warning  Missing radix parameter                 radix
   75:33  warning  Missing radix parameter                 radix
   75:57  warning  Missing radix parameter                 radix
  152:21  warning  Identifier name 'p' is too short (< 2)  id-length
  153:17  warning  Identifier name 'p' is too short (< 2)  id-length

/go-learning-series/web/grow-with-stlgo-admin/js/admin.js
  196:42  warning  A function with a name starting with an uppercase letter should only be used as a constructor  new-cap
  212:48  warning  Identifier name 'e' is too short (< 2)                                                         id-length
  407:26  warning  A function with a name starting with an uppercase letter should only be used as a constructor  new-cap

/go-learning-series/web/grow-with-stlgo-admin/js/main.js
  54:23  warning  'admin' is assigned a value but never used  no-unused-vars
  55:23  warning  'seeds' is assigned a value but never used  no-unused-vars

/go-learning-series/web/grow-with-stlgo-admin/js/main.js
  53:23  warning  'seeds' is assigned a value but never used  no-unused-vars

✖ 14 problems (0 errors, 14 warnings)

Frontend linting completed successfully
Running backend linting step...
0 issues.
Backend linting completed successfully
```

### Makefile Lint Definitions

#### Lint Makefile Commands

```makefile
.PHONY: lint
lint: tidy-lint
lint: frontend-lint
lint: backend-lint
```

#### Tidy Lint

```makefile
.PHONY: tidy-lint
tidy-lint:
 @echo "Checking that go.mod is up to date..."
 @go mod tidy
 @echo "go.mod check completed successfully"
```

#### Frontend Lint

```makefile
.PHONY: frontend-lint
frontend-lint:
 @echo "Running frontend linting step..."
 @cd $(WEB_DIR) && npx eslint --fix . && cd ..
 @echo "Frontend linting completed successfully"
```

#### Backend Lint

```makefile
.PHONY: backend-lint
backend-lint:
 @echo "Running backend linting step..."
 @$(LINTER) run --config $(LINTER_CONFIG)
 @echo "Backend linting completed successfully"
```

## Test

Make test will cause unit tests to be performed

### Example Make Test

To run a lint command run: make test

```bash
$ make test
Performing backend unit test step...
?       stl-go/grow-with-stl-go/cmd     [no test files]
?       stl-go/grow-with-stl-go/examples/channels       [no test files]
?       stl-go/grow-with-stl-go/examples/clients/REST   [no test files]
?       stl-go/grow-with-stl-go/examples/clients/WebSocket      [no test files]
?       stl-go/grow-with-stl-go/examples/cryptography   [no test files]
?       stl-go/grow-with-stl-go/examples/files  [no test files]
?       stl-go/grow-with-stl-go/examples/goroutines     [no test files]
?       stl-go/grow-with-stl-go/examples/hello-world    [no test files]
?       stl-go/grow-with-stl-go/examples/json   [no test files]
?       stl-go/grow-with-stl-go/examples/logging        [no test files]
?       stl-go/grow-with-stl-go/examples/maps_and_slices        [no test files]
?       stl-go/grow-with-stl-go/examples/primitives     [no test files]
?       stl-go/grow-with-stl-go/examples/structures     [no test files]
?       stl-go/grow-with-stl-go/examples/timed_tasks    [no test files]
ok      stl-go/grow-with-stl-go/pkg/admin       0.369s [no tests to run]
ok      stl-go/grow-with-stl-go/pkg/audit       0.378s
ok      stl-go/grow-with-stl-go/pkg/commands    0.542s
ok      stl-go/grow-with-stl-go/pkg/configs     0.745s
ok      stl-go/grow-with-stl-go/pkg/cryptography        7.941s
ok      stl-go/grow-with-stl-go/pkg/log 0.665s
ok      stl-go/grow-with-stl-go/pkg/seeds       0.615s
ok      stl-go/grow-with-stl-go/pkg/utils       0.554s
ok      stl-go/grow-with-stl-go/pkg/weather     0.416s
?       stl-go/grow-with-stl-go/pkg/webservice  [no test files]
Backend unit tests completed successfully
```

### Makefile Test Definitions

#### Test Makefile Commands

```makefile
.PHONY: unit-test
test: backend-unit-test
```

#### Backend Test

```makefile
.PHONY: backend-unit-test
backend-unit-test:
 @echo "Performing backend unit test step..."
 @go test -run $(TESTS) $(PKG) $(TESTFLAGS) $(COVER_FLAGS)
 @echo "Backend unit tests completed successfully"
```

## Coverage

Make coverage will cause the unit tests to be run and show output of the test coverage of the system

### Example Make Coverage

To run a lint command run: make coverage

```bash
$ make coverage
Performing backend unit test step...
        stl-go/grow-with-stl-go/cmd             coverage: 0.0% of statements
        stl-go/grow-with-stl-go/examples/channels               coverage: 0.0% of statements
        stl-go/grow-with-stl-go/examples/clients/REST           coverage: 0.0% of statements
        stl-go/grow-with-stl-go/examples/clients/WebSocket              coverage: 0.0% of statements
        stl-go/grow-with-stl-go/examples/cryptography           coverage: 0.0% of statements
        stl-go/grow-with-stl-go/examples/files          coverage: 0.0% of statements
        stl-go/grow-with-stl-go/examples/goroutines             coverage: 0.0% of statements
        stl-go/grow-with-stl-go/examples/hello-world            coverage: 0.0% of statements
        stl-go/grow-with-stl-go/examples/json           coverage: 0.0% of statements
        stl-go/grow-with-stl-go/examples/logging                coverage: 0.0% of statements
        stl-go/grow-with-stl-go/examples/maps_and_slices                coverage: 0.0% of statements
        stl-go/grow-with-stl-go/examples/primitives             coverage: 0.0% of statements
        stl-go/grow-with-stl-go/examples/structures             coverage: 0.0% of statements
        stl-go/grow-with-stl-go/examples/timed_tasks            coverage: 0.0% of statements
ok      stl-go/grow-with-stl-go/pkg/admin       0.598s  coverage: 0.0% of statements [no tests to run]
ok      stl-go/grow-with-stl-go/pkg/audit       0.547s  coverage: 0.0% of statements
ok      stl-go/grow-with-stl-go/pkg/commands    0.885s  coverage: 23.1% of statements
ok      stl-go/grow-with-stl-go/pkg/configs     0.382s  coverage: 31.4% of statements
ok      stl-go/grow-with-stl-go/pkg/cryptography        10.290s coverage: 44.6% of statements
ok      stl-go/grow-with-stl-go/pkg/log 0.688s  coverage: 34.1% of statements
ok      stl-go/grow-with-stl-go/pkg/seeds       0.702s  coverage: 15.8% of statements
ok      stl-go/grow-with-stl-go/pkg/utils       0.664s  coverage: 25.0% of statements
ok      stl-go/grow-with-stl-go/pkg/weather     0.275s  coverage: 0.0% of statements
        stl-go/grow-with-stl-go/pkg/webservice          coverage: 0.0% of statements
Backend unit tests completed successfully
Generating backend coverage report...
Backend coverage report completed successfully
```

### Makefile Coverage Definitions

#### Coverage Makefile Commands

```makefile
.PHONY: coverage
coverage: backend-coverage
```

#### Backend Coverage

```makefile
.PHONY: backend-coverage
backend-coverage: TESTFLAGS = -covermode=atomic -coverprofile=fullcover.out
backend-coverage: backend-unit-test
 @echo "Generating backend coverage report..."
 @grep -vE "$(COVER_EXCLUDE)" fullcover.out > $(COVER_PROFILE)
 @echo "Backend coverage report completed successfully"
```

## Verify

Make verify will cause a build, test, coverage and lint to run to verify that the system is in good shape

### Example Make Verify

To run a lint command run: make verify

```bash
$ make verify
Executing frontend build steps...

up to date, audited 112 packages in 782ms

25 packages are looking for funding
  run `npm fund` for details

found 0 vulnerabilities
Frontend build completed successfully
Executing backend build steps...
Backend build completed successfully
Performing backend unit test step...
        stl-go/grow-with-stl-go/cmd             coverage: 0.0% of statements
        stl-go/grow-with-stl-go/examples/channels               coverage: 0.0% of statements
        stl-go/grow-with-stl-go/examples/clients/REST           coverage: 0.0% of statements
        stl-go/grow-with-stl-go/examples/clients/WebSocket              coverage: 0.0% of statements
        stl-go/grow-with-stl-go/examples/cryptography           coverage: 0.0% of statements
        stl-go/grow-with-stl-go/examples/files          coverage: 0.0% of statements
        stl-go/grow-with-stl-go/examples/goroutines             coverage: 0.0% of statements
        stl-go/grow-with-stl-go/examples/hello-world            coverage: 0.0% of statements
        stl-go/grow-with-stl-go/examples/json           coverage: 0.0% of statements
        stl-go/grow-with-stl-go/examples/logging                coverage: 0.0% of statements
        stl-go/grow-with-stl-go/examples/maps_and_slices                coverage: 0.0% of statements
        stl-go/grow-with-stl-go/examples/primitives             coverage: 0.0% of statements
        stl-go/grow-with-stl-go/examples/structures             coverage: 0.0% of statements
        stl-go/grow-with-stl-go/examples/timed_tasks            coverage: 0.0% of statements
ok      stl-go/grow-with-stl-go/pkg/admin       0.510s  coverage: 0.0% of statements [no tests to run]
ok      stl-go/grow-with-stl-go/pkg/audit       0.534s  coverage: 0.0% of statements
ok      stl-go/grow-with-stl-go/pkg/commands    0.752s  coverage: 23.1% of statements
ok      stl-go/grow-with-stl-go/pkg/configs     0.615s  coverage: 30.9% of statements
ok      stl-go/grow-with-stl-go/pkg/cryptography        5.698s  coverage: 44.6% of statements
ok      stl-go/grow-with-stl-go/pkg/log 0.563s  coverage: 34.1% of statements
ok      stl-go/grow-with-stl-go/pkg/seeds       0.560s  coverage: 15.8% of statements
ok      stl-go/grow-with-stl-go/pkg/utils       0.542s  coverage: 25.0% of statements
ok      stl-go/grow-with-stl-go/pkg/weather     0.432s  coverage: 0.0% of statements
        stl-go/grow-with-stl-go/pkg/webservice          coverage: 0.0% of statements
Backend unit tests completed successfully
Generating backend coverage report...
Backend coverage report completed successfully
Checking that go.mod is up to date...
go.mod check completed successfully
Running frontend linting step...

go-learning-series/web/grow-with-stlgo-admin/js/seeds.js
   56:35  warning  Missing radix parameter                 radix
   66:58  warning  Identifier name 'c' is too short (< 2)  id-length
   72:33  warning  Missing radix parameter                 radix
   72:57  warning  Missing radix parameter                 radix
   75:33  warning  Missing radix parameter                 radix
   75:57  warning  Missing radix parameter                 radix
  152:21  warning  Identifier name 'p' is too short (< 2)  id-length
  153:17  warning  Identifier name 'p' is too short (< 2)  id-length

/go-learning-series/web/grow-with-stlgo-admin/js/admin.js
  196:42  warning  A function with a name starting with an uppercase letter should only be used as a constructor  new-cap
  212:48  warning  Identifier name 'e' is too short (< 2)                                                         id-length
  407:26  warning  A function with a name starting with an uppercase letter should only be used as a constructor  new-cap

/go-learning-series/web/grow-with-stlgo-admin/js/main.js
  54:23  warning  'admin' is assigned a value but never used  no-unused-vars
  55:23  warning  'seeds' is assigned a value but never used  no-unused-vars

/go-learning-series/web/grow-with-stlgo-admin/js/main.js
  53:23  warning  'seeds' is assigned a value but never used  no-unused-vars

✖ 14 problems (0 errors, 14 warnings)

Frontend linting completed successfully
Running backend linting step...
0 issues.
Backend linting completed successfully
```

### Makefile Verify Definitions

#### Verify Makefile Commands

```makefile
.PHONY: verify
verify: build
verify: coverage
verify: lint
```

This will call items already covered in this document

- [build](#build)
- [coverage](#coverage)
- [lint](#lint)

## Dist

Make dist will create a gziped tar archive to be created as a binary distribution of the application

### Example Make Dist

To run a build command run: make dist

```bash
$ make dist
Executing frontend build steps...

up to date, audited 112 packages in 747ms

25 packages are looking for funding
  run `npm fund` for details

found 0 vulnerabilities
Frontend build completed successfully
Executing backend build steps...
Backend build completed successfully
Executing distribution build steps...
Distribution build completed successfully
```

### Makefile Dist Definitions

#### Dist Makefile Commands

```makefile
.PHONY: dist
dist: build
dist: build-distribution
```

This will call items already covered in this document

- [build](#build)

#### Build Distribution

```makefile
.PHONY: build-distribution
build-distribution:
 @echo "Executing distribution build steps..."
 @mkdir -p $(DIST_DIR)
 @mkdir -p $(DIST_DIR)/bin
 @cp $(SCRIPT_DIR)/gwstlg.sh $(DIST_DIR)/bin/
 @chmod 755 $(DIST_DIR)/bin/gwstlg.sh
 @cp $(SCRIPT_DIR)/.gwstlg.service $(DIST_DIR)/bin/
 @cp $(MAIN)$(EXTENSION) $(DIST_DIR)/bin/
 @cp -R $(WEB_DIR) $(DIST_DIR)
 @cd $(TMP_DIR) && tar cf - gwstlg-$(COMPILED_VERSION) | gzip -9 > gwstlg-$(COMPILED_VERSION).tar.gz
 @echo "Distribution build completed successfully"
```
