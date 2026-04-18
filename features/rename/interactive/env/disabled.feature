@skipWindows
Feature: interactivity disabled, no main branch

  Background:
    Given a local Git repo
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | local     |
    And Git Town is not configured
    And the current branch is "existing"
    When I run "git-town rename new" with these environment variables
      | GIT_TOWN_INTERACTIVE | false |

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      no main branch configured and interactivity disabled via environment variable.
      
      To configure:
      git config git-town.main-branch <branch>
      """

  Scenario: undo
    When I run "git-town undo" with these environment variables
      | GIT_TOWN_INTERACTIVE | false |
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      no main branch configured and interactivity disabled via environment variable.
      
      To configure:
      git config git-town.main-branch <branch>
      """
