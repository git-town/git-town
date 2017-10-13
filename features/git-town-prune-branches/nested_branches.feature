Feature: git town-prune-branches: delete branches that were shipped or removed on another machine

  As a developer pruning branches with child branches
  I want that the pruned branches get removed from the branch hierarchy metadata
  So that my workspace is in a consistent state after pruning.

  Rules:
  - pruned branches are completely removed from the branch hierarchy


  Background:
    Given my repository has the feature branches "feature" and "feature-child"
    And the following commits exist in my repository
      | BRANCH        | LOCATION         | MESSAGE              |
      | feature       | local and remote | feature commit       |
      | feature-child | local and remote | feature-child commit |
    And Git Town is aware of this branch hierarchy
      | BRANCH        | PARENT  |
      | feature       | main    |
      | feature-child | feature |
    And the "feature" branch gets deleted on the remote
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run `git-town prune-branches`


  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND               |
      | main   | git fetch --prune     |
      |        | git branch -D feature |
    And I end up on the "main" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES            |
      | local      | main, feature-child |
      | remote     | main, feature-child |
    And Git Town is now aware of this branch hierarchy
      | BRANCH        | PARENT |
      | feature-child | main   |


  Scenario: undo
    When I run `git-town prune-branches --undo`
    Then Git Town runs the commands
      | BRANCH | COMMAND                                        |
      | main   | git branch feature <%= sha 'feature commit' %> |
    And I end up on the "main" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES                     |
      | local      | main, feature, feature-child |
      | remote     | main, feature-child          |
    And Git Town is now aware of this branch hierarchy
      | BRANCH        | PARENT  |
      | feature-child | feature |
      | feature       | main    |
