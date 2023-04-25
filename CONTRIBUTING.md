```shell
$ make help

Available targets:

  build                               Build the project and save it in the ./build folder, use the CRDA_VERSION var for setting the version
  build/all                           Build the entire project (binary and image)
  build/image                         Build the image using the the value from the CRDA_VERSION var
  build/image/push                    Build and push the image using the the value from the CRDA_VERSION var
  help                                This help screen
  lint                                Lint the code (will download golintci to the ./bin folder)
  lint/actions                        Lint the ci (will download actionlint to the ./bin folder)
  lint/all                            Lint the entire project (code, ci, dockerfile)
  lint/dockerfile                     Lint the Dockerfile (using Hadolint image, do not use inside a container)
  test                                Run all unit tests
  test/cov                            Run all unit tests and print coverage report
  test/mut                            Run mutation tests (will download gremlins to the ./bin folder)
```
