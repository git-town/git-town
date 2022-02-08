Feature: offline mode

  Background:
    Given Git Town is in offline mode
    And my repo has a feature branch "old"
    And my repo contains the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | old    | local, origin | old commit  |
    And I am on the "old" branch
    When I run "git-town rename-branch new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND            |
      | old    | git branch new old |
      |        | git checkout new   |
      | new    | git branch -D old  |
    And I am now on the "new" branch
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | new    | local         | old commit  |
      | old    | origin        | old commit  |
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | new    | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                               |
      | new    | git branch old {{ sha 'old commit' }} |
      |        | git checkout old                      |
      | old    | git branch -D new                     |
    And I am now on the "old" branch
    And now the initial commits exist
    And my repo now has its initial branches and branch hierarchy
