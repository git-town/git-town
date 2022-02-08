Feature: rename the current branch

  Background:
    Given my repo has a feature branch "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, remote | main commit |
      | old    | local, remote | old commit  |
    And I am on the "old" branch
    When I run "git-town rename-branch new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | old    | git fetch --prune --tags |
      |        | git branch new old       |
      |        | git checkout new         |
      | new    | git push -u origin new   |
      |        | git push origin :old     |
      |        | git branch -D old        |
    And I am now on the "new" branch
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, remote | main commit |
      | new    | local, remote | old commit  |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                               |
      | new    | git branch old {{ sha 'old commit' }} |
      |        | git push -u origin old                |
      |        | git push origin :new                  |
      |        | git checkout old                      |
      | old    | git branch -D new                     |
    And I am now on the "old" branch
    And my repo now has its initial branches and branch hierarchy
