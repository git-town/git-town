@skipWindows
Feature: interactivity disabled, no main branch

  Background:
    Given a local Git repo
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS |
      | branch-1 | feature | main     | local     |
      | branch-2 | feature | branch-1 | local     |
    And Git Town is not configured
    And the current branch is "branch-2"
    When I run "git-town swap --interactive=false"

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
