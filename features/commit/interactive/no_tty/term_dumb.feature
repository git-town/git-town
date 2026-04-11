Feature: TERM=dumb, no main branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
    And Git Town is not configured
    And the current branch is "branch-2"
    And an uncommitted file "changes" with content "my changes"
    And I ran "git add changes"
    When I run "git-town commit --up -m commit-1b" with these environment variables
      | TERM | dumb |

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      no main branch configured and only a dumb terminal available.

      To configure:
      git config git-town.main-branch <branch>
      """
