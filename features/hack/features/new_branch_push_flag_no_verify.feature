Feature: auto-push the new branch without running Git push hooks

  Background:
    Given setting "new-branch-push-flag" is "true"
    And the commits
      | BRANCH | LOCATION | MESSAGE       |
      | main   | origin   | origin commit |
    And the current branch is "main"

  Scenario: set to "false"
    Given setting "push-verify" is "false"
    When I run "git-town hack new"
    Then it runs the commands
      | BRANCH | COMMAND                            |
      | main   | git fetch --prune --tags           |
      |        | git rebase origin/main             |
      |        | git branch new main                |
      |        | git checkout new                   |
      | new    | git push --no-verify -u origin new |
    And the current branch is now "new"
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | origin commit |
      | new    | local, origin | origin commit |
    And this branch hierarchy exists now
      | BRANCH | PARENT |
      | new    | main   |

  Scenario: set to "true"
    Given setting "push-verify" is "true"
    When I run "git-town hack new"
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
    And this branch hierarchy exists now
      | BRANCH | PARENT |
      | new    | main   |
