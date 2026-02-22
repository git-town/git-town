@skipWindows
Feature: no TTY, missing main branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE   | PARENT | LOCATIONS     |
      | alpha | (none) |        | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME  | FILE CONTENT  |
      | alpha  | local, origin | alpha commit | alpha-file | alpha content |
    And the current branch is "alpha"
    When I run "git-town merge" in a non-TTY shell

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | alpha  | git fetch --prune --tags |
    And Git Town prints the error:
      """
      Error: cannot determine parent branch for "alpha": no interactive terminal available

      To configure, run:
      git checkout alpha && git-town set-parent <parent-branch>
      """
