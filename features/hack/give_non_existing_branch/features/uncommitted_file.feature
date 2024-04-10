@smoke
Feature: on the main branch

  Background:
    Given the current branch is a feature branch "existing"
    And the commits
      | BRANCH   | LOCATION | MESSAGE         |
      | main     | origin   | main commit     |
      | existing | local    | existing commit |
    And an uncommitted file
    When I run "git-town hack new"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                  |
      | existing | git fetch --prune --tags |
      |          | git add -A               |
      |          | git stash                |
      |          | git checkout main        |
      | main     | git rebase origin/main   |
      |          | git branch new main      |
      |          | git checkout new         |
      | new      | git stash pop            |
    And the current branch is now "new"
    And the uncommitted file still exists
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         |
      | main     | local, origin | main commit     |
      | existing | local         | existing commit |
      | new      | local         | main commit     |
    And this lineage exists now
      | BRANCH   | PARENT |
      | existing | main   |
      | new      | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND                                     |
      | new      | git add -A                                  |
      |          | git stash                                   |
      |          | git checkout main                           |
      | main     | git reset --hard {{ sha 'initial commit' }} |
      |          | git checkout existing                       |
      | existing | git branch -D new                           |
      |          | git stash pop                               |
    And the current branch is now "existing"
    And the initial commits exist
    And the initial branches and lineage exist
