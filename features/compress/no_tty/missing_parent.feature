@skipWindows
Feature: no TTY, missing parent branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE   | PARENT | LOCATIONS     |
      | feature | (none) |        | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE |
      | feature | local, origin | commit  |
    And the current branch is "feature"
    When I run "git-town compress" in a non-TTY shell

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot determine parent branch for "feature": no interactive terminal available

      To configure, run:
      git checkout feature && git-town set-parent <parent-branch>
      """
