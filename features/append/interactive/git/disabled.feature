Feature: disable interactive mode via local Git config

  Background:
    Given a local Git repo
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | local     |
    And Git Town is not configured
    And Git setting "git-town.interactive" is "false"
    And the current branch is "existing"
    When I run "git-town append new"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      no main branch configured and interactivity disabled via Git metadata.

      To configure:
      git config git-town.main-branch <branch>
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      no main branch configured and interactivity disabled via Git metadata.

      To configure:
      git config git-town.main-branch <branch>
      """
