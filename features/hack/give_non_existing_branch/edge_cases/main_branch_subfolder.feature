Feature: in a subfolder on the main branch

  Background:
    Given the commits
      | BRANCH | LOCATION | MESSAGE       | FILE NAME        |
      | main   | local    | folder commit | new_folder/file1 |
    And the current branch is "main"
    And an uncommitted file
    When I run "git-town hack new" in the "new_folder" folder

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git add -A               |
      |        | git stash                |
      |        | git rebase origin/main   |
      |        | git push                 |
      |        | git branch new main      |
      |        | git checkout new         |
      | new    | git stash pop            |
    And the current branch is now "new"
    And the uncommitted file still exists
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | folder commit |
      | new    | local         | folder commit |
    And this branch lineage exists now
      | BRANCH | PARENT |
      | new    | main   |

  Scenario: undo
    When I run "git town undo"
    Then it runs the commands
      | BRANCH | COMMAND           |
      | new    | git add -A        |
      |        | git stash         |
      |        | git checkout main |
      | main   | git branch -D new |
      |        | git stash pop     |
    And the current branch is now "main"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | folder commit |
    And no lineage exists now
