# Copyright (c) 2019 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# Export path so that the JS linting tools can get access to npm & node
# this has to be done before the shell invocation
SHELL=/bin/bash

# Obtain the version and git commit info
VERSION_FILE=version
COMPILED_VERSION=`cat $(VERSION_FILE)`

TOOLBINDIR    := tools/bin
WEBDIR        := web
LINTER        := golangci-lint
LINTER_CONFIG := .golangci.yaml

# docker
DOCKER_MAKE_TARGET  := build

# docker image options
DOCKER_REGISTRY     ?= quay.io
DOCKER_FORCE_CLEAN  ?= true
DOCKER_IMAGE_NAME   ?= grow-with-stl-go
DOCKER_IMAGE_PREFIX ?= stl-go
DOCKER_IMAGE_TAG    ?= dev
DOCKER_IMAGE        ?= $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_PREFIX)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)
DOCKER_TARGET_STAGE ?= release
PUBLISH             ?= false

# test flags
COVERAGE_OUTPUT := coverage.out

TESTFLAGS     ?= -count=1

# go options
PKG                 ?= ./...
TESTS               ?= .
COVER_FLAGS         ?=
COVER_PROFILE       ?= cover.out
COVER_EXCLUDE       ?= (zz_generated)

# Override the value of the version variable in main.go
LD_FLAGS= "-X stl-go/grow-with-stl-go/pkg/commands.version=${COMPILED_VERSION}"
GO_FLAGS  := -ldflags=$(LD_FLAGS) -trimpath
BUILD_DIR := bin

# Find all main.go files under cmd, excluding grow-with-stl-go itself
MAIN      := $(BUILD_DIR)/grow-with-stl-go
EXTENSION :=

ifeq ($(OS),Windows_NT)
	EXTENSION=.exe
endif

DIRS = internal
RECURSIVE_DIRS = $(addprefix ./, $(addsuffix /..., $(DIRS)))

### Composite Make Commands ###

.PHONY: $(MAIN)
$(MAIN): build

.PHONY: build
build: frontend-build
build: backend-build

.PHONY: lint
lint: tidy-lint
lint: frontend-lint
lint: backend-lint

.PHONY: unit-test
test: backend-unit-test

.PHONY: coverage
coverage: backend-coverage

.PHONY: verify
verify: build
verify: coverage
verify: lint

### Backend (Go) Make Commands ###

.PHONY: backend-build
backend-build:
	@echo "Executing backend build steps..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(MAIN)$(EXTENSION) $(GO_FLAGS) cmd/main.go
	@echo "Backend build completed successfully"

.PHONY: backend-unit-test
backend-unit-test:
	@echo "Performing backend unit test step..."
	@go test -run $(TESTS) $(PKG) $(TESTFLAGS) $(COVER_FLAGS)
	@echo "Backend unit tests completed successfully"

.PHONY: backend-coverage
backend-coverage: TESTFLAGS = -covermode=atomic -coverprofile=fullcover.out
backend-coverage: backend-unit-test
	@echo "Generating backend coverage report..."
	@grep -vE "$(COVER_EXCLUDE)" fullcover.out > $(COVER_PROFILE)
	@echo "Backend coverage report completed successfully"

.PHONY: backend-lint
backend-lint:
	@echo "Running backend linting step..."
	@$(LINTER) run --config $(LINTER_CONFIG)
	@echo "Backend linting completed successfully"

### Frontend (Angular) Make Commands ###

.PHONY: frontend-build
frontend-build:
	@echo "Executing frontend build steps..."
	@cd $(WEBDIR)/grow-with-stlgo && npm install && cd ../..
	@cd $(WEBDIR)/grow-with-stlgo-admin && npm install && cd ../..
	@echo "Frontend build completed successfully"

.PHONY: frontend-lint
frontend-lint:
	@echo "Running frontend linting step..."
	@cd $(WEBDIR)/grow-with-stlgo && npx eslint --fix . && cd ..
	@cd $(WEBDIR)/grow-with-stlgo-admin && npx eslint --fix . && cd ..
	@echo "Frontend linting completed successfully"

### Misc. Linting Commands ###

.PHONY: tidy-lint
tidy-lint:
	@echo "Checking that go.mod is up to date..."
	@go mod tidy
	@echo "go.mod check completed successfully"

### Docker ###

.PHONY: images
images: docker-image

.PHONY: docker-image
docker-image:
ifeq ($(USE_PROXY), true)
	@docker build . --network=host \
		--build-arg http_proxy=$(PROXY) \
		--build-arg https_proxy=$(PROXY) \
		--build-arg HTTP_PROXY=$(PROXY) \
		--build-arg HTTPS_PROXY=$(PROXY) \
		--build-arg no_proxy=$(NO_PROXY) \
		--build-arg NO_PROXY=$(NO_PROXY) \
	    --build-arg MAKE_TARGET=$(DOCKER_MAKE_TARGET) \
	    --tag $(DOCKER_IMAGE) \
	    --target $(DOCKER_TARGET_STAGE) \
	    --force-rm=$(DOCKER_FORCE_CLEAN)
else
	@docker build . --network=host \
	    --build-arg MAKE_TARGET=$(DOCKER_MAKE_TARGET) \
	    --tag $(DOCKER_IMAGE) \
	    --target $(DOCKER_TARGET_STAGE) \
	    --force-rm=$(DOCKER_FORCE_CLEAN)
endif
ifeq ($(PUBLISH), true)
	@docker push $(DOCKER_IMAGE)
endif

.PHONY: print-docker-image-tag
print-docker-image-tag:
	@echo "$(DOCKER_IMAGE)"

.PHONY: docker-image-test-suite
docker-image-test-suite: DOCKER_MAKE_TARGET = "lint cover"
docker-image-test-suite: DOCKER_TARGET_STAGE = builder
docker-image-test-suite: docker-image

.PHONY: docker-image-unit-tests
docker-image-unit-tests: DOCKER_MAKE_TARGET = coverage
docker-image-unit-tests: DOCKER_TARGET_STAGE = builder
docker-image-unit-tests: docker-image

.PHONY: docker-image-lint
docker-image-lint: DOCKER_MAKE_TARGET = lint
docker-image-lint: DOCKER_TARGET_STAGE = builder
docker-image-lint: docker-image

.PHONY: clean
clean:
	@echo "Removing build directories..."
	rm -rf $(BUILD_DIR) $(COVERAGE_OUTPUT)
	@echo "Removal completed successfully"

# The golang-unit zuul job calls the env target, so create one
# Note: on windows if there is a WSL curl in c:\windows\system32
#       it will cause problems installing the lint tools.
#       The use of cygwin curl is working however
.PHONY: env
