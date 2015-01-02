Feature: git kill: killing the current feature branch with a deleted tracking branch (with open changes)

  As a user killing the current feature branch whose remote branch has been deleted
  I want the command to succeed anyways
  So that killing branches is robust and reliable.


  Background:
    Given I have feature branches named "good-feature" and "orphaned-feature"
    And the following commits exist in my repository
      | BRANCH           | LOCATION         | MESSAGE         | FILE NAME        |
      | good-feature     | local and remote | good commit     | good_file        |
      | orphaned-feature | local and remote | orphaned commit | unfortunate_file |
    And the "orphaned-feature" branch gets deleted on the remote
    And I am on the "orphaned-feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git kill`


  Scenario: result
    Then it runs the Git commands
      | BRANCH       | COMMAND                             |
      | dead-feature | git fetch --prune                   |
      | dead-feature | git add -A                          |
      | dead-feature | git commit -m 'WIP on dead-feature' |
      | dead-feature | git checkout main                   |
      | main         | git branch -D dead-feature          |
    And I end up on the "main" branch
    And I don't have any uncommitted files
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | remote     | main, feature |
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE     | FILE NAME |
      | feature | local and remote | good commit | good_file |


  Scenario: undoing the kill
    When I run `git kill --undo`
    Then it runs the Git commands
      | BRANCH       | COMMAND                                           |
      | main         | git branch dead-feature [SHA:WIP on dead-feature] |
      | main         | git checkout dead-feature                         |
      | dead-feature | git reset [SHA:dead-end commit]                   |
    Then I end up on the "orphaned-feature" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | REPOSITORY | BRANCHES                             |
      | local      | main, orphaned-feature, good-feature |
      | remote     | main, good-feature                   |
    And I have the following commits
      | BRANCH           | LOCATION         | MESSAGE         | FILE NAME        |
      | feature          | local and remote | good commit     | good_file        |
      | orphaned-feature | local            | orphaned commit | unfortunate_file |
