# Git Town

## Development Guidelines

- you can change any file in the current folder and its subfolders
- never change files outside the Git repository
- never create new Git branches
- never make Git commits
- I will review the changes you make and then commit them on my own.

## Automated testing

To run all unit tests for the project, use this command:

```bash
make unit
```

### Linters

Please execute the linters after making changes to verify the correctness of
your changes. To run them, use the following command:

```bash
make lint
```

### End-to-End Tests

End-to-end tests are defined in the "features" directory. They take a while to
execute, so only run them to verify that everything still works after you are
done making changes. To run all end-to-end tests:

```bash
make cuke
```
