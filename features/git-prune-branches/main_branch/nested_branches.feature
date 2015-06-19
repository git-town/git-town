Feature: git prune-branches: remove pruned branches from the branch hierarchy

  As a developer pruning branches with child branches
  I want that the pruned branches get removed from the branch hierarchy metadata
  So that my workspace is in a consistent state after pruning.


  Background:
    Given the following commits exist in my repository
      | BRANCH | LOCATION         | MESSAGE     | FILE NAME |
      | main   | local and remote | main commit | main_file |
    And I have a stale feature branch named "feature-1-stale" with its tip at "Initial commit"
    And I have a feature branch named "feature-2-active" ahead of main
    And Git Town is aware of this branch hierarchy
      | BRANCH           | PARENT          |
      | feature-1-stale  | main            |
      | feature-2-active | feature-1-stale |
    And I am on the "main" branch
    When I run `git prune-branches`


  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND                          |
      | main   | git fetch --prune                |
      |        | git push origin :feature-1-stale |
      |        | git branch -d feature-1-stale    |
    And I end up on the "main" branch
    And the existing branches are
      | REPOSITORY | BRANCHES               |
      | local      | main, feature-2-active |
      | remote     | main, feature-2-active |
    And Git Town is now aware of this branch hierarchy
      | BRANCH           | PARENT          |
      | feature-2-active | feature-1-stale |


  Scenario: undoing the prune
    When I run `git prune-branches --undo`
    Then it runs the Git commands
      | BRANCH | COMMAND                                                |
      | main   | git branch feature-1-stale <%= sha 'Initial commit' %> |
      |        | git push -u origin feature-1-stale                     |
    And I end up on the "main" branch
    Then the existing branches are
      | REPOSITORY | BRANCHES                                |
      | local      | main, feature-1-stale, feature-2-active |
      | remote     | main, feature-1-stale, feature-2-active |
    And Git Town is now aware of this branch hierarchy
      | BRANCH           | PARENT          |
      | feature-1-stale  | main            |
      | feature-2-active | feature-1-stale |
