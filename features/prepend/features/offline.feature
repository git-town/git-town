Feature: offline mode

  Background:
    Given offline mode is enabled
    And a feature branch "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |
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
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |
    And Git Town is now aware of this branch hierarchy
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
    And now the initial commits exist
    And Git Town is now aware of the initial branch hierarchy
