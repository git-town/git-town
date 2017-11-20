# Measuring Test Coverage

Measuring test coverage is a bit trickier than normal
because the end-to-end tests of Git Town are in Ruby.
We therefore measure test coverage using these steps:

1. Compile a test binary called `git-town.test`.
   This is done via `bin/build-test`
   The test binary creates a coverage file called `coverage.cov`
   each time it is run.

2. Run that test binary in our tests,
   generating coverage data for that particular test.
   For example, instead of `git-town hack foo`,
   our tests now run:

   ```
   ./git-town.test -test.coverprofile=coverage.cov hack foo
   ```

3. After each test,
   save the coverage file.

4. When all tests are run, merge all coverage files together
   into the coverage file for the entire code base.

5. Travis-CI uploads that coverage file to Coveralls.io.


To see the coverage, look at Coveralls.
You can also see it locally by running `bin/cuke` and then:

```
go tool cover -html=coverage.cov
```
