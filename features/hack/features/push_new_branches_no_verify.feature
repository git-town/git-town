Feature: auto-push the new branch without running Git push hooks

  Background:
    Given Git Town setting "push-new-branches" is "true"
    And the commits
      | BRANCH | LOCATION | MESSAGE       |
      | main   | origin   | origin commit |
    And the current branch is "main"

  Scenario: set to "false"
    Given Git Town setting "push-hook" is "false"
    When I run "git-town hack new"
    Then it runs the commands
      | BRANCH | COMMAND                            |
      | main   | git fetch --prune --tags           |
      |        | git rebase origin/main             |
      |        | git branch new main                |
      |        | git checkout new                   |
      | new    | git push --no-verify -u origin new |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | origin commit |
      | new    | local, origin | origin commit |
    And this branch lineage exists now
      | BRANCH | PARENT |
      | new    | main   |

  Scenario: set to "true"
    Given Git Town setting "push-hook" is "true"
    When I run "git-town hack new"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git rebase origin/main   |
      |        | git branch new main      |
      |        | git checkout new         |
      | new    | git push -u origin new   |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | origin commit |
      | new    | local, origin | origin commit |
    And this branch lineage exists now
      | BRANCH | PARENT |
      | new    | main   |
