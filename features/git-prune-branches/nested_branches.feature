Feature: git prune-branches: delete branches that were shipped or removed on another machine

  As a developer pruning branches with child branches
  I want that the pruned branches get removed from the branch hierarchy metadata
  So that my workspace is in a consistent state after pruning.

  Rules:
  - the branch hierarchy metadata of pruned branches is removed


  Background:
    Given I have feature branches named "feature-1" and "feature-1-child"
    And the following commits exist in my repository
      | BRANCH          | LOCATION         | MESSAGE                |
      | feature-1       | local and remote | feature-1 commit       |
      | feature-1-child | local and remote | feature-1-child commit |
    And Git Town is aware of this branch hierarchy
      | BRANCH          | PARENT    |
      | feature-1       | main      |
      | feature-1-child | feature-1 |
    And the "feature-1" branch gets deleted on the remote
    And I am on the "main" branch
    And I have an uncommitted file
    When I run `git prune-branches`


  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                 |
      | main   | git fetch --prune       |
      |        | git branch -D feature-1 |
    And I end up on the "main" branch
    And I still have my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES              |
      | local      | main, feature-1-child |
      | remote     | main, feature-1-child |
    And Git Town is now aware of this branch hierarchy
      | BRANCH          | PARENT |
      | feature-1-child | main   |


  Scenario: undo
    When I run `git prune-branches --undo`
    Then it runs the commands
      | BRANCH | COMMAND                                            |
      | main   | git branch feature-1 <%= sha 'feature-1 commit' %> |
    And I end up on the "main" branch
    And I still have my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES                         |
      | local      | main, feature-1, feature-1-child |
      | remote     | main, feature-1-child            |
    And Git Town is now aware of this branch hierarchy
      | BRANCH          | PARENT    |
      | feature-1-child | feature-1 |
      | feature-1       | main      |
