Feature: rename a contribution branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS     |
      | contribution | contribution | main   | local, origin |
    And the current branch is "contribution"
    And the commits
      | BRANCH       | LOCATION      | MESSAGE               |
      | contribution | local, origin | somebody elses commit |
    When I run "git-town rename-branch contribution new"

  Scenario: result
    Then it runs the commands
      | BRANCH       | COMMAND                            |
      | contribution | git fetch --prune --tags           |
      |              | git branch --move contribution new |
      |              | git checkout new                   |
      | new          | git push -u origin new             |
      |              | git push origin :contribution      |
    And the current branch is now "new"
    And the contribution branches are now "new"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE               |
      | new    | local, origin | somebody elses commit |
    And this lineage exists now
      | BRANCH | PARENT |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH       | COMMAND                                                   |
      | new          | git branch contribution {{ sha 'somebody elses commit' }} |
      |              | git push -u origin contribution                           |
      |              | git push origin :new                                      |
      |              | git checkout contribution                                 |
      | contribution | git branch -D new                                         |
    And the current branch is now "contribution"
    And the contribution branches are now "contribution"
    And the initial commits exist now
    And the initial branches and lineage exist now
