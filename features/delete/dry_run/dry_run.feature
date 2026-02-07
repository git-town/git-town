Feature: dry-run deleting the current feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | current | feature | main   | local, origin |
      | other   | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | current | local, origin | current commit |
      | other   | local, origin | other commit   |
    And the current branch is "current"
    And the current branch is "current" and the previous branch is "other"
    When I run "git-town delete --dry-run"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | current | git fetch --prune --tags |
      |         | git push origin :current |
      |         | git checkout other       |
      | other   | git branch -D current    |
    And the initial branches and lineage exist now
    And the initial commits exist now
  #
  # Cannot test undo because dry-run now doesn't create a runstate.
