Feature: on the main branch with an upstream repo

  Background:
    Given a Git repo clone
    And an upstream repo
    And the commits
      | BRANCH | LOCATION | MESSAGE         |
      | main   | upstream | upstream commit |
    And the current branch is "main"
    And Git Town setting "sync-upstream" is "true"
    And I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git rebase origin/main   |
      |        | git fetch upstream main  |
      |        | git rebase upstream/main |
      |        | git push                 |
      |        | git push --tags          |
    And all branches are now synchronized
    And the current branch is still "main"
    And these commits exist now
      | BRANCH | LOCATION                | MESSAGE         |
      | main   | local, origin, upstream | upstream commit |

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "main"
    And these commits exist now
      | BRANCH | LOCATION                | MESSAGE         |
      | main   | local, origin, upstream | upstream commit |
    And the initial branches and lineage exist
