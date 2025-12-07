---
name: end-to-end-tests
description: ensure the product works
allowed-tools: make, go
---

# End-to-end tests

## Instructions

Run end-to-end tests when you are done making changes. Fix all printed errors
and re-run the end-to-end tests until it finishes with exit code 0.

End-to-end tests are defined in the "features" folder. They are written in the
Gherkin language and execute using Cucumber for Go.

To run all end-to-end tests:

```bash
make cuke
```

If an end-to-end test fails, you can re-run only this test to verify that your
fix works:

- Add the `@debug @this` tag to a specific scenario and then run `make cukethis`
  to execute only the tagged scenario
