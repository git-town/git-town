Feature: git prune-branches: don't remove the current empty feature branch if there are open changes

  As a developer pruning branches
  I don't want my current empty branch deleted if I have open changes
  So that I can prune my branches without losing current work.


  Background:
    Given I have a feature branch named "feature" behind main
    And I have a feature branch named "stale_feature" behind main
    And I am on the "feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git prune-branches`


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND                        |
      | feature | git fetch --prune              |
      | feature | git stash -u                   |
      | feature | git checkout main              |
      | main    | git push origin :stale_feature |
      | main    | git branch -d stale_feature    |
      | main    | git checkout feature           |
      | feature | git stash pop                  |
    And I end up on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | remote     | main, feature |
      | coworker   | main          |


  Scenario: undoing the prune
    When I run `git prune-branches --undo`
    Then it runs the Git commands
      | BRANCH  | COMMAND                                              |
      | feature | git stash -u                                         |
      | feature | git checkout main                                    |
      | main    | git branch stale_feature [SHA:behind feature commit] |
      | main    | git push -u origin stale_feature                     |
      | main    | git checkout feature                                 |
      | feature | git stash pop                                        |
    And I end up on the "feature" branch
    Then the existing branches are
      | REPOSITORY | BRANCHES                     |
      | local      | main, feature, stale_feature |
      | remote     | main, feature, stale_feature |
      | coworker   | main                         |
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
