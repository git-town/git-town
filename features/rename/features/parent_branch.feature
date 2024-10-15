Feature: rename a parent branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And the current branch is "parent"
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | child  | local, origin | child commit  |
      | parent | local, origin | parent commit |
    When I run "git-town rename parent new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                      |
      | parent | git fetch --prune --tags     |
      |        | git branch --move parent new |
      |        | git checkout new             |
      | new    | git push -u origin new       |
      |        | git push origin :parent      |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | child  | local, origin | child commit  |
      | new    | local, origin | parent commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | child  | new    |
      | new    | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                     |
      | new    | git branch parent {{ sha 'parent commit' }} |
      |        | git push -u origin parent                   |
      |        | git checkout parent                         |
      | parent | git branch -D new                           |
      |        | git push origin :new                        |
    And the current branch is now "parent"
    And the initial commits exist now
    And the initial branches and lineage exist now
