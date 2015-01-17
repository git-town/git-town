Feature: git kill: killing the current feature branch with a deleted tracking branch (with open changes)

  As a user killing the current feature branch whose remote branch has been deleted
  I want the command to succeed anyways
  So that killing branches is robust and reliable.


  Background:
    Given I have feature branches named "active-feature" and "orphaned-feature"
    And the following commits exist in my repository
      | BRANCH           | LOCATION         | MESSAGE         | FILE NAME        |
      | active-feature   | local and remote | active commit   | good_file        |
      | orphaned-feature | local and remote | orphaned commit | unfortunate_file |
    And the "orphaned-feature" branch gets deleted on the remote
    And I am on the "orphaned-feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git kill`


  Scenario: result
    Then it runs the Git commands
      | BRANCH           | COMMAND                                 |
      | orphaned-feature | git fetch --prune                       |
      | orphaned-feature | git add -A                              |
      | orphaned-feature | git commit -m 'WIP on orphaned-feature' |
      | orphaned-feature | git checkout main                       |
      | main             | git branch -D orphaned-feature          |
    And I end up on the "main" branch
    And I don't have any uncommitted files
    And the existing branches are
      | REPOSITORY | BRANCHES             |
      | local      | main, active-feature |
      | remote     | main, active-feature |
    And I have the following commits
      | BRANCH         | LOCATION         | MESSAGE       | FILE NAME |
      | active-feature | local and remote | active commit | good_file |


  Scenario: undoing the kill
    When I run `git kill --undo`
    Then it runs the Git commands
      | BRANCH           | COMMAND                                                   |
      | main             | git branch orphaned-feature [SHA:WIP on orphaned-feature] |
      | main             | git checkout orphaned-feature                             |
      | orphaned-feature | git reset [SHA:orphaned commit]                           |
    And I end up on the "orphaned-feature" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | REPOSITORY | BRANCHES                               |
      | local      | main, orphaned-feature, active-feature |
      | remote     | main, active-feature                   |
    And I have the following commits
      | BRANCH           | LOCATION         | MESSAGE         | FILE NAME        |
      | active-feature   | local and remote | active commit   | good_file        |
      | orphaned-feature | local            | orphaned commit | unfortunate_file |
