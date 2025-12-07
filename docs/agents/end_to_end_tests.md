# End-to-End Tests

End-to-end tests are defined in the "features" directory. They are written in
Cucumber. Run them after you are done making changes. To run all end-to-end
tests:

```bash
make cuke
```

If an end-to-end test fails, you can debug it this way:

- Add the `@this` tag to a specific scenario and then run `make cukethis` to
  execute only the tagged scenario
- Add the `@debug @this` tags to a Cucumber scenario to see the CLI output of
  the Git Town command that the test executes
