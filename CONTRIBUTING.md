```shell
$ make help

Available targets:

build:             Build the project and save it in the ./build folder, use the CRDA_VERSION var for setting the version
build/all:         Build the entire project (binary and image)
build/image:       Build the image using the the value from the CRDA_VERSION var
build/image/push:  Build and push the image using the the value from the CRDA_VERSION var
download/openapi:  Download the backend's openapi.yaml specification file and save at the project's root
generate/openapi:  Generate code from an ./openapi.yaml spec file (do not use in CI)
help:              This help screen
lint:              Lint the code (will download golintci to the ./bin folder)
lint/all:          Lint the entire project (code, ci, dockerfile)
lint/ci:           Lint the ci (will download actionlint to the ./bin folder)
lint/dockerfile:   Lint the Dockerfile (using Hadolint image, do not use inside a container)
test:              Run all unit tests
test/cov:          Run all unit tests and print coverage report, use the COVERAGE_THRESHOLD var for setting threshold
test/mut:          Run mutation tests (will download gremlins to the ./bin folder)

```
