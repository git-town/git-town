Feature: auto-push new branches

  Background:
    Given the new-branch-push-flag configuration is true
    And my repo has a feature branch "old"
    And my repo contains the commits
      | BRANCH | LOCATION      | MESSAGE        |
      | old    | local, remote | feature commit |
    And I am on the "old" branch
    When I run "git-town prepend new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | old    | git fetch --prune --tags |
      |        | git checkout main        |
      | main   | git rebase origin/main   |
      |        | git branch new main      |
      |        | git checkout new         |
      | new    | git push -u origin new   |
    And I am now on the "new" branch
    And my repo now has the commits
      | BRANCH | LOCATION      | MESSAGE        |
      | old    | local, remote | feature commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | new    | main   |
      | old    | new    |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND              |
      | new    | git push origin :new |
      |        | git checkout main    |
      | main   | git branch -d new    |
      |        | git checkout old     |
    And I am now on the "old" branch
    And my repo is left with my original commits
    And Git Town now has the initial branch hierarchy
