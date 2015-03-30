Feature: git kill: killing the current feature branch with a deleted tracking branch (without open changes)

  (see ./with_open_changes.feature)


  Background:
    Given I have feature branches named "other-feature" and "current-feature"
    And the following commits exist in my repository
      | BRANCH          | LOCATION         | MESSAGE                | FILE NAME        |
      | other-feature   | local and remote | other feature commit   | good_file        |
      | current-feature | local and remote | current feature commit | unfortunate_file |
    And the "current-feature" branch gets deleted on the remote
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
      | BRANCH        | LOCATION         | MESSAGE              | FILE NAME |
      | other-feature | local and remote | other feature commit | good_file |


  Scenario: undoing the kill
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
      | BRANCH          | LOCATION         | MESSAGE                | FILE NAME        |
      | other-feature   | local and remote | other feature commit   | good_file        |
      | current-feature | local            | current feature commit | unfortunate_file |
