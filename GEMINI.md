# Git Town

This project contains a Go application.

## Development Guidelines

General development guidelines are in the
[developer documentation](docs/DEVELOPMENT.md).

When making changes, you can change any file in the current folder and its
subfolders. Never change files outside this folder. Never create new Git
branches, and never make any Git commits. I will review the changes you make and
then commit them on my own.

## Automated testing

Git Town has an extensive suite of automated tests that assists in developing
correctly functioning and bug-free code.

### Unit Tests

To run all unit tests for the project, use the following command:

```bash
make unit
```

### Linters

Please execute the linters after making changes to verify the correctness of
your changes. To run them, use the following command:

```
make lint
```

### End-to-End Tests

End-to-end tests are located in the "features" directory. They take a while to
execute, so only run them to verify that everything still works after you are
done making changes. To run all end-to-end tests:

```
make cuke
```
