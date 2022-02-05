Feature: rename a parent branch

  Background:
    Given my repo has a feature branch "parent"
    And my repo has a feature branch "child" as a child of "parent"
    And my repo contains the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | child  | local, remote | child commit  |
      | parent | local, remote | parent commit |
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
    And my repo now has the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | child  | local, remote | child commit  |
      | new    | local, remote | parent commit |
    And Git Town now knows this branch hierarchy
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
    And my repo is left with my initial commits
    And my repo now has its initial branches and branch hierarchy
