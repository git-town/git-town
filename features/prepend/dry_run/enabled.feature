Feature: dry-run prepending a branch to a feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |
    And the current branch is "old"
    When I run "git-town prepend parent --dry-run"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      | old    | git fetch --prune --tags    |
      |        | git checkout -b parent main |
    And the initial branches and lineage exist now
    And the initial commits exist now
  #
  # Cannot test undo because dry-run now doesn't create a runstate.
