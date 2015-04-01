Feature: git prune-branches: keep used feature branches when run on a feature branch (without open changes)

  As a developer pruning branches
  I want my feature branches with commits to not be deleted
  So that I can keep my repository clean without losing work.


  Background:
    Given I have a feature branch named "feature" ahead of main
    And I have a feature branch named "stale_feature" behind main
    And I am on the "feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git prune-branches`


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND                        |
      | feature | git fetch --prune              |
      |         | git stash -u                   |
      |         | git checkout main              |
      | main    | git push origin :stale_feature |
      |         | git branch -d stale_feature    |
      |         | git checkout feature           |
      | feature | git stash pop                  |
    And I end up on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | remote     | main, feature |
      | coworker   | main          |


  Scenario: undoing the operation
    When I run `git prune-branches --undo`
    Then it runs the Git commands
      | BRANCH  | COMMAND                                              |
      | feature | git stash -u                                         |
      | feature | git checkout main                                    |
      | main    | git branch stale_feature <%= sha 'Initial commit' %> |
      | main    | git push -u origin stale_feature                     |
      | main    | git checkout feature                                 |
      | feature | git stash pop                                        |
    And I end up on the "feature" branch
    Then the existing branches are
      | REPOSITORY | BRANCHES                     |
      | local      | main, feature, stale_feature |
      | remote     | main, feature, stale_feature |
      | coworker   | main                         |
