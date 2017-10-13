Feature: git town-kill: killing the current feature branch with child branches

  As a user killing the current feature branch that has child branches
  I want that the current branch is cleanly removed from the branch hierarchy metadata
  So that killing branches is robust and reliable.


  Background:
    Given my repository has a feature branch named "feature-1"
    And my repository has a feature branch named "feature-2" as a child of "feature-1"
    And it has a feature branch named "feature-3" as a child of "feature-2"
    And the following commits exist in my repository
      | BRANCH    | LOCATION         | MESSAGE          |
      | feature-1 | local and remote | feature 1 commit |
      | feature-2 | local and remote | feature 2 commit |
      | feature-3 | local and remote | feature 3 commit |
    And I am on the "feature-2" branch
    And my workspace has an uncommitted file
    When I run `git-town kill`


  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                          |
      | feature-2 | git fetch --prune                |
      |           | git push origin :feature-2       |
      |           | git add -A                       |
      |           | git commit -m "WIP on feature-2" |
      |           | git checkout feature-1           |
      | feature-1 | git branch -D feature-2          |
    And I end up on the "feature-1" branch
    And I don't have any uncommitted files
    And the existing branches are
      | REPOSITORY | BRANCHES                   |
      | local      | main, feature-1, feature-3 |
      | remote     | main, feature-1, feature-3 |
    And my repository has the following commits
      | BRANCH    | LOCATION         | MESSAGE          |
      | feature-1 | local and remote | feature 1 commit |
      | feature-3 | local and remote | feature 3 commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT    |
      | feature-1 | main      |
      | feature-3 | feature-1 |


  Scenario: undoing the kill
    When I run `git-town kill --undo`
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                            |
      | feature-1 | git branch feature-2 <%= sha 'WIP on feature-2' %> |
      |           | git checkout feature-2                             |
      | feature-2 | git reset <%= sha 'feature 2 commit' %>            |
      |           | git push -u origin feature-2                       |
    And I end up on the "feature-2" branch
    And my workspace has the uncommitted file again
    And the existing branches are
      | REPOSITORY | BRANCHES                              |
      | local      | main, feature-1, feature-2, feature-3 |
      | remote     | main, feature-1, feature-2, feature-3 |
    And my repository is left with my original commits
    And Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT    |
      | feature-1 | main      |
      | feature-2 | feature-1 |
      | feature-3 | feature-2 |
