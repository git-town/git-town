Feature: rename a prototype branch

  Background:
    Given a Git repo with origin
    And the branch
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | prototype | prototype | main   | local, origin |
    And the current branch is "prototype"
    And the commits
      | BRANCH    | LOCATION      | MESSAGE             |
      | prototype | local, origin | experimental commit |
    When I run "git-town rename-branch prototype new"

  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                    |
      | prototype | git fetch --prune --tags   |
      |           | git branch new prototype   |
      |           | git checkout new           |
      | new       | git push -u origin new     |
      |           | git push origin :prototype |
      |           | git branch -D prototype    |
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
      | BRANCH    | COMMAND                                              |
      | new       | git branch prototype {{ sha 'experimental commit' }} |
      |           | git push -u origin prototype                         |
      |           | git push origin :new                                 |
      |           | git checkout prototype                               |
      | prototype | git branch -D new                                    |
    And the current branch is now "prototype"
    And the prototype branches are now "prototype"
    And the initial commits exist
    And the initial branches and lineage exist
