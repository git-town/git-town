Feature: prepend a branch to a feature branch

  Background:
    Given my repo has a feature branch "old"
    And my repo contains the commits
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, remote | old commit |
    And I am on the "old" branch
    And my workspace has an uncommitted file
    When I run "git-town prepend parent"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | old    | git fetch --prune --tags |
      |        | git add -A               |
      |        | git stash                |
      |        | git checkout main        |
      | main   | git rebase origin/main   |
      |        | git branch parent main   |
      |        | git checkout parent      |
      | parent | git stash pop            |
    And I am now on the "parent" branch
    And my workspace still contains my uncommitted file
    And my repo now has the commits
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, remote | old commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | old    | parent |
      | parent | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND              |
      | parent | git add -A           |
      |        | git stash            |
      |        | git checkout main    |
      | main   | git branch -d parent |
      |        | git checkout old     |
      | old    | git stash pop        |
    And I am now on the "old" branch
    And my workspace still contains my uncommitted file
    And my repo is left with my initial commits
    And Git Town now has the initial branch hierarchy
