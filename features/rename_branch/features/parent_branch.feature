Feature: rename a parent branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    Given the current branch is "parent"
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | child  | local, origin | child commit  |
      | parent | local, origin | parent commit |
    When I run "git-town rename-branch parent new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | parent | git fetch --prune --tags |
      |        | git branch new parent    |
      |        | git checkout new         |
      | new    | git push -u origin new   |
      |        | git push origin :parent  |
      |        | git branch -D parent     |
    And the current branch is now "new"
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
      |        | git push origin :new                        |
      |        | git checkout parent                         |
      | parent | git branch -D new                           |
    And the current branch is now "parent"
    And the initial commits exist
    And the initial branches and lineage exist
