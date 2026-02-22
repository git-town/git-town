@skipWindows
Feature: no TTY, missing main branch

  Scenario: main branch
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE   | PARENT | LOCATIONS |
      | feature | (none) |        | local     |
    And Git Town is not configured
    And the current branch is "feature"
    When I run "git-town diff-parent" in a non-TTY shell
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      Error: no main branch configured and no interactive terminal available.

      To configure, run "git config git-town.main-branch <branch>".
      To set up interactively, run "git town init" in a shell with TTY.
      """
