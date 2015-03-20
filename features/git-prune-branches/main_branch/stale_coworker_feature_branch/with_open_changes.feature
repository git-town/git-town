Feature: git prune-branches: remove stale coworker branches when run on the main branch (without open changes)

  As a developer pruning branches
  I want my coworker's merged branches to be deleted from the remote repository
  So that all remaining branches are relevant and my team can focus on their current work.


  Background:
    Given my coworker has a feature branch named "stale_feature" behind main
    And I am on the "main" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git prune-branches`


  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND                        |
      | main   | git fetch --prune              |
      | main   | git push origin :stale_feature |
    And I end up on the "main" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | REPOSITORY | BRANCHES            |
      | local      | main                |
      | remote     | main                |
      | coworker   | main, stale_feature |


  Scenario: undoing the operation
    When I run `git prune-branches --undo`
    Then it runs the Git commands
      | BRANCH | COMMAND                                              |
      | main   | git branch stale_feature <%= sha 'Initial commit' %> |
      | main   | git push -u origin stale_feature                     |
    And I end up on the "main" branch
    Then the existing branches are
      | REPOSITORY | BRANCHES            |
      | local      | main, stale_feature |
      | remote     | main, stale_feature |
      | coworker   | main, stale_feature |
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
