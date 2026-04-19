@skipWindows
Feature: interactivity disabled, no main branch

  Background:
    Given a Git repo with origin
    And the committed configuration file:
      """
      interactive = false
      """
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | local     |
    And Git Town is not configured
    And the current branch is "existing"
    When I run "git-town prepend new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                  |
      | existing | git fetch --prune --tags |
    And Git Town prints the error:
      """
      no main branch configured and interactivity disabled via config file.
      
      To configure:
      git config git-town.main-branch <branch>
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      no main branch configured and interactivity disabled via config file.
      
      To configure:
      git config git-town.main-branch <branch>
      """
