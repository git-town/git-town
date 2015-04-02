Feature: git kill: killing the current feature branch without a tracking branch (without open changes)

  (see ./with_open_changes.feature)


  Background:
    Given I have a feature branch named "other-feature"
    And I have a local feature branch named "current-feature"
    And the following commits exist in my repository
      | BRANCH          | LOCATION         | MESSAGE                | FILE NAME            |
      | other-feature   | local and remote | other feature commit   | other_feature_file   |
      | current-feature | local            | current feature commit | current_feature_file |
    And I am on the "current-feature" branch
    When I run `git kill`


  Scenario: result
    Then it runs the Git commands
      | BRANCH          | COMMAND                       |
      | current-feature | git fetch --prune             |
      |                 | git checkout main             |
      | main            | git branch -D current-feature |
    And I end up on the "main" branch
    And the existing branches are
      | REPOSITORY | BRANCHES            |
      | local      | main, other-feature |
      | remote     | main, other-feature |
    And I have the following commits
      | BRANCH        | LOCATION         | MESSAGE              | FILE NAME          |
      | other-feature | local and remote | other feature commit | other_feature_file |


  Scenario: Undoing a kill of a local feature branch
    When I run `git kill --undo`
    Then it runs the Git commands
      | BRANCH | COMMAND                                                        |
      | main   | git branch current-feature <%= sha 'current feature commit' %> |
      |        | git checkout current-feature                                   |
    And I end up on the "current-feature" branch
    And the existing branches are
      | REPOSITORY | BRANCHES                             |
      | local      | main, current-feature, other-feature |
      | remote     | main, other-feature                  |
    And I have the following commits
      | BRANCH          | LOCATION         | MESSAGE                | FILE NAME            |
      | current-feature | local            | current feature commit | current_feature_file |
      | other-feature   | local and remote | other feature commit   | other_feature_file   |
