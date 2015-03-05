Feature: git prune-branches: don't remove the current empty feature branch if there are open changes

  As a developer pruning branches
  I don't want my current empty branch deleted if I have open changes
  So that I can prune my branches without losing current work.


  Background:
    Given I have a feature branch named "stale_feature_1"
    And I have a feature branch named "stale_feature_2"
    And the following commits exist in my repository
      | BRANCH | LOCATION         | MESSAGE     |
      | main   | local and remote | main commit |
    And I am on the "stale_feature_2" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git prune-branches`


  Scenario: result
    Then it runs the Git commands
      | BRANCH          | COMMAND                          |
      | stale_feature_2 | git fetch --prune                |
      | stale_feature_2 | git stash -u                     |
      | stale_feature_2 | git checkout main                |
      | main            | git push origin :stale_feature_1 |
      | main            | git branch -d stale_feature_1    |
      | main            | git checkout stale_feature_2     |
      | stale_feature_2 | git stash pop                    |
    And I end up on the "stale_feature_2" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | REPOSITORY | BRANCHES              |
      | local      | main, stale_feature_2 |
      | remote     | main, stale_feature_2 |
      | coworker   | main                  |


  Scenario: undoing the prune
    When I run `git prune-branches --undo`
    Then it runs the Git commands
      | BRANCH          | COMMAND                                      |
      | stale_feature_2 | git stash -u                                 |
      | stale_feature_2 | git checkout main                            |
      | main            | git branch stale_feature_1 [SHA:Initial commit] |
      | main            | git push -u origin stale_feature_1           |
      | main            | git checkout stale_feature_2                 |
      | stale_feature_2 | git stash pop                                |
    And I end up on the "stale_feature_2" branch
    Then the existing branches are
      | REPOSITORY | BRANCHES                               |
      | local      | main, stale_feature_1, stale_feature_2 |
      | remote     | main, stale_feature_1, stale_feature_2 |
      | coworker   | main                                   |
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
