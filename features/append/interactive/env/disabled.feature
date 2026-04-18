Feature: disable interactive mode via environment variable

  Background:
    Given a local Git repo
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | local     |
    And the current branch is "existing"
    And Git Town is not configured
    When I run "git-town append new" with these environment variables
      | GIT_TOWN_INTERACTIVE | false |

  @this
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
