Feature: a parent branch of a local branch was shipped

  Background:
    Given my repo has a feature branch "parent"
    And my repo has a feature branch "child" as a child of "parent"
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | parent | local, origin | parent commit |
      | child  | local, origin | child commit  |
    And origin deletes the "parent" branch
    And I am on the "main" branch
    When I run "git-town prune-branches"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git branch -D parent     |
    And I am now on the "main" branch
    And the existing branches are
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, child |
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | child  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                     |
      | main   | git branch parent {{ sha 'parent commit' }} |
    And I am now on the "main" branch
    And my repo now has its initial branches and branch hierarchy
