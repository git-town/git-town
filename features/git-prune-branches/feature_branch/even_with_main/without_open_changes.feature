Feature: git prune-branches: don't remove the current empty feature branch if there are open changes

  As a developer pruning branches
  I don't want my current empty branch deleted if I have open changes
  So that I can prune my branches without losing current work.


  Background:
    Given the following commits exist in my repository
      | BRANCH | LOCATION         | MESSAGE     | FILE NAME |
      | main   | local and remote | main commit | main_file |
    And I have a stale feature branch named "stale_feature_1" with its tip at "main commit"
    And I have a stale feature branch named "stale_feature_2" with its tip at "main commit"
    And I am on the "stale_feature_1" branch
    When I run `git prune-branches`


  Scenario: result
    Then it runs the Git commands
      | BRANCH          | COMMAND                          |
      | stale_feature_1 | git fetch --prune                |
      |                 | git checkout main                |
      | main            | git push origin :stale_feature_1 |
      |                 | git branch -d stale_feature_1    |
      |                 | git push origin :stale_feature_2 |
      |                 | git branch -d stale_feature_2    |
    And I end up on the "main" branch
    And the existing branches are
      | REPOSITORY | BRANCHES              |
      | local      | main|
      | remote     | main|
      | coworker   | main                  |


  Scenario: undoing the prune
    When I run `git prune-branches --undo`
    Then it runs the Git commands
      | BRANCH          | COMMAND                                             |
      | main            | git branch stale_feature_2 <%= sha 'main commit' %> |
      |                 | git push -u origin stale_feature_2                  |
      |                 | git branch stale_feature_1 <%= sha 'main commit' %> |
      |                 | git push -u origin stale_feature_1                  |
      |                 | git checkout stale_feature_1                        |
    And I end up on the "stale_feature_1" branch
    Then the existing branches are
      | REPOSITORY | BRANCHES                               |
      | local      | main, stale_feature_1, stale_feature_2 |
      | remote     | main, stale_feature_1, stale_feature_2 |
      | coworker   | main                                   |
