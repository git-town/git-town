Feature: auto-push the new branch to origin

  Background:
    Given the "new-branch-push-flag" setting is "true"
    And the commits
      | BRANCH | LOCATION | MESSAGE     |
      | main   | origin   | main commit |
    And I am on the "main" branch
    When I run "git-town append new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git rebase origin/main   |
      |        | git branch new main      |
      |        | git checkout new         |
      | new    | git push -u origin new   |
    And I am now on the "new" branch
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | new    | local, origin | main commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | new    | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND              |
      | new    | git push origin :new |
      |        | git checkout main    |
      | main   | git branch -d new    |
    And I am now on the "main" branch
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
    And Git Town is now aware of no branch hierarchy
