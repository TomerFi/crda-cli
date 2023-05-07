CRDA_VERSION ?= staging# set version for prod ie CRDA_VERSION=1.2.3 (do not use v)
COMMIT_HASH = $(strip $(shell git rev-parse --short HEAD))
TIMESTAMP =  $(strip $(shell date +%s))
BUILD_DATE = $(strip $(shell date -u +"%Y-%m-%dT%H:%M:%SZ"))
VENDOR_NAME = Red Hat, Inc.

COVERAGE_THRESHOLD ?= 60

IMAGE_BUILDER ?= podman
IMAGE_NAME ?= quay.io/ecosystem-appeng/crda-cli
FULL_IMAGE_NAME = $(strip $(IMAGE_NAME):$(CRDA_VERSION))

# if this is modified, modify the FROM instruction at the final stage in Dockerfile
BASE_IMAGE_NAME = registry.access.redhat.com/ubi9/go-toolset:1.18.10-4

# get os and architecture and save as OS_ARCH
OS_ARCH = $(shell go env GOOS)-$(shell go env GOARCH)

default: help

.PHONY: help
## This help screen
help: make2help

LOCALBIN = $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

LOCALBUILD = $(shell pwd)/build
$(LOCALBUILD):
	mkdir -p $(LOCALBUILD)

OPENAPI_FILENAME = openapi.yaml
OPENAPI_SPEC = $(shell pwd)/${OPENAPI_FILENAME}

.PHONY: test
## Run all unit tests
test:
	go test -v ./...

.PHONY: test/cov
## Run all unit tests and print coverage report, use the COVERAGE_THRESHOLD var for setting threshold
test/cov: test/cov/report go-test-coverage

.PHONY: test/cov/report
test/cov/report:
	go test -failfast -coverprofile=cov.out -v ./...
	go tool cover -func=cov.out
	go tool cover -html=cov.out -o cov.html

.PHONY: test/mut
## Run mutation tests (will download gremlins to the ./bin folder)
test/mut: gremlins

LDFLAGS=-ldflags="\
-X 'github.com/rhecosystemappeng/crda-cli/pkg/utils.version=${CRDA_VERSION}' \
-X 'github.com/rhecosystemappeng/crda-cli/pkg/utils.commitHash=${COMMIT_HASH}' \
-X 'github.com/rhecosystemappeng/crda-cli/pkg/utils.timestamp=${TIMESTAMP}' \
-X 'github.com/rhecosystemappeng/crda-cli/pkg/utils.vendorInfo=${VENDOR_NAME}' \
"

.PHONY: build/all
## Build the entire project (binary and image)
build/all: build build/image

.PHONY: build
## Build the project and save it in the ./build folder, use the CRDA_VERSION var for setting the version
build:
	go build ${LDFLAGS} -o ${LOCALBUILD}/crda-${CRDA_VERSION}-${OS_ARCH} main.go

.PHONY: build/image
## Build the image using the the value from the CRDA_VERSION var
build/image:
	digest=$(${IMAGE_BUILDER} image inspect --format '{{ index .Digest }}' ${BASE_IMAGE_NAME})
	${IMAGE_BUILDER} build \
	--build-arg BASE_IMAGE_NAME=${BASE_IMAGE_NAME} \
	--build-arg BASE_IMAGE_DIGEST=${digest} \
	--build-arg BUILD_DATE=${BUILD_DATE} \
	--build-arg COMMIT_HASH=${COMMIT_HASH} \
	--build-arg CRDA_VERSION=${CRDA_VERSION} \
	--tag ${FULL_IMAGE_NAME} .

.PHONY: build/image/push
## Build and push the image using the the value from the CRDA_VERSION var
build/image/push: build/image
	${IMAGE_BUILDER} push ${FULL_IMAGE_NAME}

.PHONY: lint/all
## Lint the entire project (code, ci, dockerfile)
lint/all: lint lint/ci lint/dockerfile

.PHONY: lint
## Lint the code (will download golintci to the ./bin folder)
lint: fmt golintci

.PHONY: lint/ci
## Lint the ci (will download actionlint to the ./bin folder)
lint/ci: actionlint

.PHONY: lint/dockerfile
## Lint the Dockerfile (using Hadolint image, do not use inside a container)
lint/dockerfile:
	${IMAGE_BUILDER} run --rm -i docker.io/hadolint/hadolint:latest < Dockerfile

.PHONY: generate/openapi
## Generate code from an ./openapi.yaml spec file (do not use in CI)
generate/openapi: oapi_codegen

.PHONY: download/openapi
## Download the backend's openapi.yaml specification file and save at the project's root
download/openapi: remove/${OPENAPI_FILENAME} ${OPENAPI_SPEC}

.PHONY: remove/${OPENAPI_FILENAME}
remove/${OPENAPI_FILENAME}:
	rm -f ${OPENAPI_FILENAME}

.PHONY: fmt
fmt:
	go fmt ./...

GOLINTCI_BIN = ${LOCALBIN}/golangci-lint

.PHONY: golintci
golintci: ${GOLINTCI_BIN}
	${GOLINTCI_BIN} run

${GOLINTCI_BIN}:
	GOBIN=${LOCALBIN} go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

ACTIONLINT_BIN = ${LOCALBIN}/actionlint

.PHONY: actionlint
actionlint: ${ACTIONLINT_BIN}
	${ACTIONLINT_BIN} --verbose

# recommendation: manually install shellcheck and verify it's on your PATH, it will be picked up by actionlint
${ACTIONLINT_BIN}:
	GOBIN=${LOCALBIN} go install github.com/rhysd/actionlint/cmd/actionlint@latest

GREMLINS_BIN = ${LOCALBIN}/gremlins

.PHONY: gremlins
gremlins: ${GREMLINS_BIN}
	${GREMLINS_BIN} unleash

${GREMLINS_BIN}:
	GOBIN=${LOCALBIN} go install github.com/go-gremlins/gremlins/cmd/gremlins@latest

OAPI_CODEGEN_BIN = ${LOCALBIN}/oapi-codegen

.PHONY: oapi_codegen
oapi_codegen: ${OAPI_CODEGEN_BIN} ${OPENAPI_SPEC}
	${OAPI_CODEGEN_BIN} -generate types -package api -o pkg/backend/api/types_generated.go ${OPENAPI_SPEC}

${OAPI_CODEGEN_BIN}:
	GOBIN=${LOCALBIN} go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

${OPENAPI_SPEC}:
	wget https://raw.githubusercontent.com/RHEcosystemAppEng/crda-backend/main/src/main/resources/META-INF/${OPENAPI_FILENAME}

MAKE2HELP_BIN = ${LOCALBIN}/make2help

.PHONY: make2help
make2help: ${MAKE2HELP_BIN}
	@printf "Available targets:\n\n"
	@${MAKE2HELP_BIN}  $(MAKEFILE_LIST)
	@printf "\n"

${MAKE2HELP_BIN}:
	GOBIN=${LOCALBIN} go install github.com/Songmu/make2help/cmd/make2help@latest

GO_TEST_COVERAGE_BIN = ${LOCALBIN}/go-test-coverage

go-test-coverage: ${GO_TEST_COVERAGE_BIN}
	${GO_TEST_COVERAGE_BIN} -p cov.out -k 0 -t ${COVERAGE_THRESHOLD}

${GO_TEST_COVERAGE_BIN}:
	GOBIN=${LOCALBIN} go install github.com/vladopajic/go-test-coverage/v2@latest
