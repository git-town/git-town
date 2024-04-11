Feature: on a forked repo

  Background:
    Given an upstream repo
    And the commits
      | BRANCH | LOCATION | MESSAGE         |
      | main   | upstream | upstream commit |
    And the current branch is "main"
    When I run "git-town hack new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git rebase origin/main   |
      |        | git fetch upstream main  |
      |        | git rebase upstream/main |
      |        | git push                 |
      |        | git branch new main      |
      |        | git checkout new         |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH | LOCATION                | MESSAGE         |
      | main   | local, origin, upstream | upstream commit |
      | new    | local                   | upstream commit |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND           |
      | new    | git checkout main |
      | main   | git branch -D new |
    And the current branch is now "main"
    And these commits exist now
      | BRANCH | LOCATION                | MESSAGE         |
      | main   | local, origin, upstream | upstream commit |
    And no lineage exists now
