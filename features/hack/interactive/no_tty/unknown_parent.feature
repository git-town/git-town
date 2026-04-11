@skipWindows
Feature: no TTY, unknown parent

  Scenario:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE   | PARENT | LOCATIONS     |
      | branch | (none) |        | local, origin |
    And the current branch is "branch"
    When I run "git-town hack new" in a non-TTY shell
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | branch | git fetch --prune --tags |
      |        | git checkout main        |
      | main   | git checkout -b new      |
