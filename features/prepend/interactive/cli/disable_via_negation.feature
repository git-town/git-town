@skipWindows
Feature: interactivity disabled, no main branch

  Background:
    Given a local Git repo
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | local     |
    And Git Town is not configured
    And the current branch is "existing"
    When I run "git-town prepend new --non-interactive"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      no main branch configured and interactivity disabled via CLI.

      To configure:
      git config git-town.main-branch <branch>
      """

  Scenario: undo
    When I run "git-town undo --non-interactive"
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      no main branch configured and interactivity disabled via CLI.

      To configure:
      git config git-town.main-branch <branch>
      """
