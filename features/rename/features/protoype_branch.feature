Feature: rename a prototype branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | prototype | prototype | main   | local, origin |
    And the current branch is "prototype"
    And the commits
      | BRANCH    | LOCATION      | MESSAGE             |
      | prototype | local, origin | experimental commit |
    When I run "git-town rename prototype new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                         |
      | prototype | git fetch --prune --tags        |
      |           | git branch --move prototype new |
      |           | git checkout new                |
      | new       | git push -u origin new          |
      |           | git push origin :prototype      |
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
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                              |
      | new       | git branch prototype {{ sha 'experimental commit' }} |
      |           | git push -u origin prototype                         |
      |           | git checkout prototype                               |
      | prototype | git branch -D new                                    |
      |           | git push origin :new                                 |
    And the current branch is now "prototype"
    And the prototype branches are now "prototype"
    And the initial commits exist now
    And the initial branches and lineage exist now
