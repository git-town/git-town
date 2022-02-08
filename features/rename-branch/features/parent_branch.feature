Feature: rename a parent branch

  Background:
    Given my repo has a feature branch "parent"
    And my repo has a feature branch "child" as a child of "parent"
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | child  | local, origin | child commit  |
      | parent | local, origin | parent commit |
    And I am on the "parent" branch
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
    And I am now on the "new" branch
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE       |
      | child  | local, origin | child commit  |
      | new    | local, origin | parent commit |
    And Git Town is now aware of this branch hierarchy
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
    And I am now on the "parent" branch
    And now the initial commits exist
    And my repo now has its initial branches and branch hierarchy
