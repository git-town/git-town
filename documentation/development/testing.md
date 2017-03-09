# Testing

## Running Tests

```bash
# running the different test types
rake           # runs all tests
bin/lint       # runs the linters
bin/lint_bash  # runs the Bash linters
bin/lint_go    # runs the Go linters
rake test      # runs the feature tests

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


## Mocking

Certain tests require the Git remote to be set to a real value
on GitHub or Bitbucket.
This causes `git push` operations to go to GitHub during testing,
which is undesirable.
To circumvent this problem, Git Town allows to mock the Git remote
by setting the Git configuration value
`git-town.testing.remote-url` to the respective value.
To keep this behavior clean and secure,
this also requires an environment variable `GIT_TOWN_ENV` to be set to `test`.


## Auto-running tests

The Git Town code base works with
[Tertestrial](https://github.com/Originate/tertestrial-server).
