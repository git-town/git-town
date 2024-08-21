Feature: rename a prototype branch

  Background:
    Given a Git repo with origin
    And the branch
      | NAME       | TYPE      | PARENT | LOCATIONS     |
      | experiment | prototype | main   | local, origin |
    And the current branch is "experiment"
    And the commits
      | BRANCH     | LOCATION      | MESSAGE             |
      | experiment | local, origin | experimental commit |
    When I run "git-town rename-branch experiment new"

  Scenario: result
    Then it runs the commands
      | BRANCH     | COMMAND                     |
      | experiment | git fetch --prune --tags    |
      |            | git branch new experiment   |
      |            | git checkout new            |
      | new        | git push -u origin new      |
      |            | git push origin :experiment |
      |            | git branch -D experiment    |
    And the current branch is now "new"
    And the prototype branches are now "new"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE             |
      | new    | local, origin | experimental commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | new    | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH     | COMMAND                                               |
      | new        | git branch experiment {{ sha 'experimental commit' }} |
      |            | git push -u origin experiment                         |
      |            | git push origin :new                                  |
      |            | git checkout experiment                               |
      | experiment | git branch -D new                                     |
    And the current branch is now "experiment"
    And the prototype branches are now "experiment"
    And the initial commits exist
    And the initial branches and lineage exist
