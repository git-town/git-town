Feature: disable interactive mode via local Git config

  Background:
    Given a local Git repo
    And Git setting "git-town.interactive" is "false"
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | local     |
    And the current branch is "existing"
    And Git Town is not configured
    When I run "git-town append new"

  @debug @this
  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      no main branch configured and interactivity disabled via CLI.
      
      To configure:
      git config git-town.main-branch <branch>
      """

  Scenario: undo
    When I run "git-town undo --interactive=false"
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      no main branch configured and interactivity disabled via CLI.
      
      To configure:
      git config git-town.main-branch <branch>
      """
