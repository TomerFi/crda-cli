```shell
$ make help

Available targets:

  build                               Build the project and save it in the ./build folder, use the VERSION and VENDOR_NAME vars for prod build
  build/full                          Lint, test, and build the project in the ./build folder
  help                                This help screen
  lint                                Run linters against the project (will download golintci to the ./bin folder)
  test                                Run all unit tests
  test/cov                            Run all unit tests and print coverage report
  test/mut                            Run mutation tests (will download gremlins to the ./bin folder)
```
