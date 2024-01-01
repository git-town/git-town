Feature: on the main branch with an upstream repo

  Background:
    Given an upstream repo
    And the commits
      | BRANCH | LOCATION | MESSAGE         |
      | main   | local    | local commit    |
      |        | origin   | origin commit   |
      |        | upstream | upstream commit |
    And the current branch is "main"
    And Git Town setting "sync-upstream" is "false"
    And I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git rebase origin/main   |
      |        | git push                 |
      |        | git push --tags          |
    And all branches are now synchronized
    And the current branch is still "main"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE         |
      | main   | local, origin | origin commit   |
      |        |               | local commit    |
      |        | upstream      | upstream commit |

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "main"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE         |
      | main   | local, origin | origin commit   |
      |        |               | local commit    |
      |        | upstream      | upstream commit |
    And the initial branches and lineage exist
