Feature: git prune-branches: remove stale coworker branches when run on the main branch

  As a developer pruning branches
  I want my coworker's merged branches to be deleted from the remote repository
  So that all remaining branches are relevant and my team can focus on their current work.


  Background:
    Given my coworker has a feature branch named "stale-feature" behind main
    And I am on the "main" branch
    And I have an uncommitted file
    When I run `git prune-branches`


  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                        |
      | main   | git fetch --prune              |
      |        | git push origin :stale-feature |
    And I end up on the "main" branch
    And I still have my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES            |
      | local      | main                |
      | remote     | main                |
      | coworker   | main, stale-feature |


  Scenario: undoing the operation
    When I run `git prune-branches --undo`
    Then it runs the commands
      | BRANCH | COMMAND                                              |
      | main   | git branch stale-feature <%= sha 'Initial commit' %> |
      |        | git push -u origin stale-feature                     |
    And I end up on the "main" branch
    Then the existing branches are
      | REPOSITORY | BRANCHES            |
      | local      | main, stale-feature |
      | remote     | main, stale-feature |
      | coworker   | main, stale-feature |
    And I still have my uncommitted file
