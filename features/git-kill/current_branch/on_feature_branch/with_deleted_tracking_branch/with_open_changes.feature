Feature: git kill: killing the current feature branch with a deleted tracking branch (with open changes)

  As a user killing the current feature branch whose remote branch has been deleted
  I want the command to succeed anyways
  So that killing branches is robust and reliable.


  Background:
    Given I have feature branches named "current-feature" and "other-feature"
    And the following commits exist in my repository
      | BRANCH          | LOCATION         | MESSAGE                | FILE NAME            |
      | current-feature | local and remote | current feature commit | current_feature_file |
      | other-feature   | local and remote | other feature commit   | other_feature_file   |
    And the "current-feature" branch gets deleted on the remote
    And I am on the "current-feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git kill`


  Scenario: result
    Then it runs the Git commands
      | BRANCH          | COMMAND                                |
      | current-feature | git fetch --prune                      |
      | current-feature | git add -A                             |
      | current-feature | git commit -m 'WIP on current-feature' |
      | current-feature | git checkout main                      |
      | main            | git branch -D current-feature          |
    And I end up on the "main" branch
    And I don't have any uncommitted files
    And the existing branches are
      | REPOSITORY | BRANCHES            |
      | local      | main, other-feature |
      | remote     | main, other-feature |
    And I have the following commits
      | BRANCH        | LOCATION         | MESSAGE              | FILE NAME          |
      | other-feature | local and remote | other feature commit | other_feature_file |


  Scenario: undoing the kill
    When I run `git kill --undo`
    Then it runs the Git commands
      | BRANCH          | COMMAND                                                        |
      | main            | git branch current-feature <%= sha 'WIP on current-feature' %> |
      | main            | git checkout current-feature                                   |
      | current-feature | git reset <%= sha 'current feature commit' %>                  |
    And I end up on the "current-feature" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | REPOSITORY | BRANCHES                             |
      | local      | main, current-feature, other-feature |
      | remote     | main, other-feature                  |
    And I have the following commits
      | BRANCH          | LOCATION         | MESSAGE                | FILE NAME            |
      | other-feature   | local and remote | other feature commit   | other_feature_file   |
      | current-feature | local            | current feature commit | current_feature_file |
