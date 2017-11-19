# Measuring Test Coverage

- compile the binary with test coverage

  ```
  go test -c -coverpkg ./...
  ```

  This creates a file `./git-town.test`

- run the binary, outputting coverage data

  ```
  ./git-town.test -test.coverprofile=coverage.cov config
  ```

  This creates a file `./coverage.cov`


- see the coverage

  ```
  go tool cover -html=coverage.cov
  ```
