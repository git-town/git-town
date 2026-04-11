Feature: TERM=dumb, no main branch

  Background:
    Given a local Git repo
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS |
      | branch-1 | feature | main     | local     |
      | branch-2 | feature | branch-1 | local     |
    And Git Town is not configured
    And the current branch is "branch-2"
    When I run "git-town swap" with these environment variables
      | TERM | dumb |

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      no main branch configured and only a dumb terminal available.

      To configure:
      git config git-town.main-branch <branch>
      """

  Scenario: undo
    When I run "git-town undo" with these environment variables
      | TERM | dumb |
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      no main branch configured and only a dumb terminal available.

      To configure:
      git config git-town.main-branch <branch>
      """
