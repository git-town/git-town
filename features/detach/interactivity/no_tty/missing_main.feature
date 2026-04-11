@skipWindows
Feature: no TTY, missing parent branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE   | PARENT | LOCATIONS     |
      | branch-1 | (none) |        | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  |
      | branch-1 | local, origin | commit 1 |
    And Git Town is not configured
    And the current branch is "branch-1"
    When I run "git-town detach" in a non-TTY shell

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                  |
      | branch-1 | git fetch --prune --tags |
    And Git Town prints the error:
      """
      no main branch configured and no interactive terminal available.

      To configure, run "git config git-town.main-branch <branch>".
      To set up interactively, run "git town init" in a shell with TTY.
      """
