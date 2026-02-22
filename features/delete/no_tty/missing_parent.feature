@skipWindows
Feature: no TTY, missing parent branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE   | PARENT | LOCATIONS     |
      | current | (none) |        | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | current | local, origin | current commit |
    And the current branch is "current"
    When I run "git-town delete" in a non-TTY shell

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | current | git fetch --prune --tags |
      |         | git push origin :current |
      |         | git checkout main        |
      | main    | git branch -D current    |
