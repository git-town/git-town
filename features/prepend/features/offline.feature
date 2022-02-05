Feature: offline mode

  Background:
    Given Git Town is in offline mode
    And my repo has a feature branch "old"
    And my repo contains the commits
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, remote | old commit |
    And I am on the "old" branch
    When I run "git-town prepend new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                |
      | old    | git checkout main      |
      | main   | git rebase origin/main |
      |        | git branch new main    |
      |        | git checkout new       |
    And I am now on the "new" branch
    And my repo now has the commits
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, remote | old commit |
    And Git Town now knows this branch hierarchy
      | BRANCH | PARENT |
      | new    | main   |
      | old    | new    |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND           |
      | new    | git checkout main |
      | main   | git branch -d new |
      |        | git checkout old  |
    And I am now on the "old" branch
    And my repo is left with my initial commits
    And Git Town now knows the initial branch hierarchy
