Feature: git kill: killing the current feature branch with a deleted tracking branch (without open changes)

  (see ./deleted_tracking_branch_with_open_changes.feature)


  Background:
    Given I have feature branches named "good-feature" and "orphaned-feature"
    And the following commits exist in my repository
      | BRANCH           | LOCATION         | MESSAGE         | FILE NAME        |
      | good-feature     | local and remote | good commit     | good_file        |
      | orphaned-feature | local and remote | orphaned commit | unfortunate_file |
    And the "orphaned-feature" branch gets deleted on the remote
    And I am on the "orphaned-feature" branch
    When I run `git kill`


  Scenario: result
    Then it runs the Git commands
      | BRANCH           | COMMAND                        |
      | orphaned-feature | git fetch --prune              |
      | orphaned-feature | git checkout main              |
      | main             | git branch -D orphaned-feature |
    And I end up on the "main" branch
    And the existing branches are
      | REPOSITORY | BRANCHES           |
      | local      | main, good-feature |
      | remote     | main, good-feature |
    And I have the following commits
      | BRANCH       | LOCATION         | MESSAGE     | FILE NAME |
      | good-feature | local and remote | good commit | good_file |


  Scenario: undoing the kill
    When I run `git kill --undo`
    Then it runs the Git commands
      | BRANCH | COMMAND                                           |
      | main   | git branch orphaned-feature [SHA:orphaned commit] |
      | main   | git checkout orphaned-feature                     |
    And I end up on the "orphaned-feature" branch
    And the existing branches are
      | REPOSITORY | BRANCHES                             |
      | local      | main, orphaned-feature, good-feature |
      | remote     | main, good-feature                   |
    And I have the following commits
      | BRANCH           | LOCATION         | MESSAGE         | FILE NAME        |
      | good-feature     | local and remote | good commit     | good_file        |
      | orphaned-feature | local            | orphaned commit | unfortunate_file |
