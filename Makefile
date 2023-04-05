VERSION ?= 0.0.0# set version for prod ie VERSION=1.2.3 (do not use v)
COMMIT_HASH ?= $(strip $(shell git rev-parse --short HEAD))
TIMESTAMP ?=  $(strip $(shell date +%s))
VENDOR_NAME ?= development# set to Red Hat for prod

# get os and architecture and save as OS_ARCH
ifeq ($(OS),Windows_NT)
	ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
		OS_ARCH := win32-amd64
	endif
	ifeq ($(PROCESSOR_ARCHITECTURE),x86)
		OS_ARCH := win32-ia32
	endif
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		OS_ARCH := linux
	endif
	ifeq ($(UNAME_S),Darwin)
		OS_ARCH := darwin
	endif
	UNAME_P := $(shell uname -p)
	ifeq ($(UNAME_P),x86_64)
		OS_ARCH := ${OS_ARCH}-amd64
	endif
	ifneq ($(filter %86,$(UNAME_P)),)
		OS_ARCH := ${OS_ARCH}-ia32
	endif
	ifneq ($(filter arm%,$(UNAME_P)),)
		OS_ARCH := ${OS_ARCH}-arm
	endif
endif

default: help

.PHONY: help
## This help screen
help:
	@printf "Available targets:\n\n"
	@awk '/^[a-zA-Z\-_0-9%:\\]+/ { \
	  helpMessage = match(lastLine, /^## (.*)/); \
	  if (helpMessage) { \
	    helpCommand = $$1; \
	    helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
      gsub("\\\\", "", helpCommand); \
      gsub(":+$$", "", helpCommand); \
	    printf "  \x1b[32;01m%-35s\x1b[0m %s\n", helpCommand, helpMessage; \
	  } \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST) | sort -u
	@printf "\n"

LOCALBIN = $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

LOCALBUILD = $(shell pwd)/build
$(LOCALBUILD):
	mkdir -p $(LOCALBUILD)

.PHONY: test
## Run all unit tests
test:
	go test -v ./...

.PHONY: test/cov
## Run all unit tests and print coverage report
test/cov:
	go test -coverprofile=cov.out -v ./...
	go tool cover -func=cov.out
	go tool cover -html=cov.out -o cov.html

.PHONY: test/mut
## Run mutation tests (will download gremlins to the ./bin folder)
test/mut: gremlins

LDFLAGS=-ldflags="\
-X 'github.com/rhecosystemappeng/crda-cli/pkg/utils.version=${VERSION}' \
-X 'github.com/rhecosystemappeng/crda-cli/pkg/utils.commitHash=${COMMIT_HASH}' \
-X 'github.com/rhecosystemappeng/crda-cli/pkg/utils.timestamp=${TIMESTAMP}' \
-X 'github.com/rhecosystemappeng/crda-cli/pkg/utils.vendorInfo=${VENDOR_NAME}' \
"

.PHONY: build
## Build the project and save it in the ./build folder, use the VERSION and VENDOR_NAME vars for prod build
build:
	go build ${LDFLAGS} -o ${LOCALBUILD}/crda-${VERSION}-${OS_ARCH} main.go

.PHONY: build/full
## Lint, test, and build the project in the ./build folder
build/full: lint test/cov build

.PHONY: lint
## Run linters against the project (will download golintci to the ./bin folder)
lint: fmt golintci

.PHONY: fmt
fmt:
	go fmt ./...

GOLINTCI_BIN = ${LOCALBIN}/golangci-lint

.PHONY: golintci
golintci: ${GOLINTCI_BIN}
	${GOLINTCI_BIN} run

${GOLINTCI_BIN}:
	GOBIN=${LOCALBIN} go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

GREMLINS_BIN = ${LOCALBIN}/gremlins

.PHONY: gremlins
gremlins: ${GREMLINS_BIN}
	${GREMLINS_BIN} unleash

${GREMLINS_BIN}:
	GOBIN=${LOCALBIN} go install github.com/go-gremlins/gremlins/cmd/gremlins@latest
