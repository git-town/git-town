@skipWindows
Feature: no TTY, missing parent branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE   | PARENT | LOCATIONS     |
      | branch-1 | (none) |        | local, origin |
      | branch-2 | (none) |        | local, origin |
    And the current branch is "branch-2"
    And I ran "git-town set-parent branch-1"
    And the current branch is "branch-1"
    When I run "git-town delete" in a non-TTY shell

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                   |
      | branch-1 | git fetch --prune --tags  |
      |          | git push origin :branch-1 |
      |          | git checkout branch-2     |
      | branch-2 | git branch -D branch-1    |
