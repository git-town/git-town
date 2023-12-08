Feature: inside an uncommitted subfolder on the current feature branch

  Background:
    Given the current branch is a feature branch "existing"
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
    And an uncommitted file in folder "new_folder"
    When I run "git-town hack new" in the "new_folder" folder

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
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | new    | local         | main commit |
    And this branch lineage exists now
      | BRANCH   | PARENT |
      | existing | main   |
      | new      | main   |

  Scenario: undo
    When I run "git town undo"
    Then it runs the commands
      | BRANCH   | COMMAND               |
      | new      | git add -A            |
      |          | git stash             |
      |          | git checkout existing |
      | existing | git branch -D new     |
      |          | git stash pop         |
    And the current branch is now "existing"
    And the initial commits exist
    And the initial lineage exists
