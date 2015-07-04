Feature: git kill: killing the given branch with child branches

  (see ../../current_branch/on_feature_branch/with_child_branches.feature)


  Background:
    Given I have a feature branch named "feature-1"
    And I have a feature branch named "feature-2" as a child of "feature-1"
    And I have a feature branch named "feature-3" as a child of "feature-2"
    And the following commits exist in my repository
      | BRANCH    | LOCATION         | MESSAGE          | FILE NAME      |
      | feature-1 | local and remote | feature 1 commit | feature_1_file |
      | feature-2 | local and remote | feature 2 commit | feature_2_file |
      | feature-3 | local and remote | feature 3 commit | feature_3_file |
    And I am on the "feature-3" branch
    And I have an uncommitted file
    When I run `git kill feature-2`


  Scenario: result
    Then it runs the Git commands
      | BRANCH    | COMMAND                    |
      | feature-3 | git fetch --prune          |
      |           | git push origin :feature-2 |
      |           | git branch -D feature-2    |
    And I end up on the "feature-3" branch
    And I still have my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES                   |
      | local      | main, feature-1, feature-3 |
      | remote     | main, feature-1, feature-3 |
    And I have the following commits
      | BRANCH    | LOCATION         | MESSAGE          | FILE NAME      |
      | feature-1 | local and remote | feature 1 commit | feature_1_file |
      | feature-3 | local and remote | feature 3 commit | feature_3_file |
    And Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT    |
      | feature-1 | main      |
      | feature-3 | feature-1 |


  Scenario: undoing the kill
    When I run `git kill --undo`
    Then it runs the Git commands
      | BRANCH    | COMMAND                                            |
      | feature-3 | git branch feature-2 <%= sha 'feature 2 commit' %> |
      |           | git push -u origin feature-2                       |
    And I end up on the "feature-3" branch
    And I again have my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES                              |
      | local      | main, feature-1, feature-2, feature-3 |
      | remote     | main, feature-1, feature-2, feature-3 |
    And I am left with my original commits
    And Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT    |
      | feature-1 | main      |
      | feature-2 | feature-1 |
      | feature-3 | feature-2 |
