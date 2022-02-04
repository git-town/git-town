Feature: auto-push the new branch

  Background:
    Given the new-branch-push-flag configuration is true
    And my repo contains the commits
      | BRANCH | LOCATION | MESSAGE       |
      | main   | remote   | remote commit |
    And I am on the "main" branch
    When I run "git-town hack new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git rebase origin/main   |
      |        | git branch new main      |
      |        | git checkout new         |
      | new    | git push -u origin new   |
    And I am now on the "new" branch
    And my repo now has the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, remote | remote commit |
      | new    | local, remote | remote commit |
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
    And I am now on the "main" branch
    And my repo now has the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, remote | remote commit |
    And Git Town now has no branch hierarchy information
