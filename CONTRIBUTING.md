```shell
$ make help

Available targets:

  build                               Build the project and save it in the ./build folder, use the CRDA_VERSION and VENDOR_NAME vars for prod build
  build/full                          Lint, test, and build the project in the ./build folder
  build/image                         Build the image with using the version as a tag
  build/image/push                    Build and push the image with using the version as a tag
  help                                This help screen
  lint                                Run linters against the project (will download golintci to the ./bin folder)
  test                                Run all unit tests
  test/cov                            Run all unit tests and print coverage report
  test/mut                            Run mutation tests (will download gremlins to the ./bin folder)
```
