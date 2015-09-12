# Testing

## Running Tests

```bash
# running the different test types
rake         # runs all tests
rake lint    # runs the linters
rake test    # runs the feature tests

# running individual scenarios/features
cucumber <filename>[:<lineno>]
cucumber -n '<scenario/feature name>'

# running tests in parallel
bin/cuke [cucumber parameters]
# set the environment variable PARALLEL_TEST_PROCESSORS to override the
# auto-detected number of processors

# auto-fixing formatting issues
rake format
```

The `rake [parameters]` commands above can also be run as `bundle exec rake [parameters]`
if you encounter Ruby versioning issues.

Git Town's [CI server](https://circleci.com/gh/Originate/git-town)
automatically tests all commits and pull requests,
and notifies you via email and through status badges in pull requests
about problems.


## Debugging

To see the output of the Git commands run in tests, you can set the
`DEBUG_COMMANDS` environment variable while running your specs:

```bash
$ DEBUG_COMMANDS=true cucumber <filename>[:<lineno>]
```

Alternatively, you can also add a `@debug-commands` flag to the respective
Cucumber spec:

  ```cucumber
  @debug-commands
  Scenario: foo bar baz
    Given ...
  ```

For even more detailed output, you can use the `DEBUG` variable or tag
in a similar fashion.
If set, Git Town prints every shell command executed during the tests
(includes setup, inspection of the Git status, and the Git commands),
and the respective console output.
