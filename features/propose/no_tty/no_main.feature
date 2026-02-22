@skipWindows
Feature: no TTY, no main branch

  Background:
    Given a local Git repo
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | local     |
    And Git Town is not configured
    And the current branch is "existing"
    When I run "git-town propose" in a non-TTY shell

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      no main branch configured and no interactive terminal available.

      To configure, run "git config git-town.main-branch <branch>".
      To set up interactively, run "git town init" in a shell with TTY.
      """

  Scenario: undo
    When I run "git-town undo" in a non-TTY shell
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      no main branch configured and no interactive terminal available.

      To configure, run "git config git-town.main-branch <branch>".
      To set up interactively, run "git town init" in a shell with TTY.
      """
