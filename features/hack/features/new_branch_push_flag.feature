Feature: auto-push the new branch

  Background:
    Given the "new-branch-push-flag" setting is "true"
    And the commits
      | BRANCH | LOCATION | MESSAGE       |
      | main   | origin   | origin commit |
    And the current branch is "main"
    When I run "git-town hack new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git rebase origin/main   |
      |        | git branch new main      |
      |        | git checkout new         |
      | new    | git push -u origin new   |
    And the current branch is now "new"
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | origin commit |
      | new    | local, origin | origin commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | new    | main   |

  Scenario: undo
    When I run "git town undo"
    Then it runs the commands
      | BRANCH | COMMAND              |
      | new    | git push origin :new |
      |        | git checkout main    |
      | main   | git branch -d new    |
    And the current branch is now "main"
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | origin commit |
    And Git Town is now aware of no branch hierarchy
