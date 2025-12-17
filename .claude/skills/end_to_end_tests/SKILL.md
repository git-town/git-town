---
name: end-to-end-tests
description: after making changes, run end-to-end tests to ensure that the product still works
allowed-tools: make, go
---

# End-to-end tests

## Instructions

When you are done making changes, run all end-to-end tests through this command:

```bash
make cuke
```

If the exit code is 1, read the error messages and fix all errors, then re-run
the end-to-end tests until they run successfully and exit with code 0.

End-to-end tests are defined in the "features" folder. They are written in the
Gherkin language and execute using Cucumber for Go.

If an end-to-end test fails, you can re-run only this test to verify that your
fix works:

- Add the `@debug @this` tag to a specific scenario
- Run `make cukethis` to execute only the tagged scenario
- when the test passes, remove the `@debug @this` tags
