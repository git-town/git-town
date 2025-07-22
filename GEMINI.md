# Git Town

This project contains a Go application.

## Development Guidelines

General development guidelines are in the
[developer documentation](docs/DEVELOPMENT.md).

## Automated testing

Git Town has an extensive suite of automated tests that assists in developing
correctly functioning and bug-free code.

### Unit Tests

Unit tests are located in the respective package directories alongside the
source code. To run all unit tests for the project, use the following command:

```bash
make unit
```

### Linters

Git Town leans heavily on linters. Please execute them after making changes to
verify their correctness. To run them, use the following command:

```
make lint
```

### End-to-End Tests

End-to-end tests are located in the "features" directory. To run all end-to-end
tests:

```
make cuke
```
