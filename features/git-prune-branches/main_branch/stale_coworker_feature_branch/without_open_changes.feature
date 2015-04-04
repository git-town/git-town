Feature: git prune-branches: remove stale coworker branches when run on the main branch (without open changes)

  (see ./with_open_changes.feature)


  Background:
    Given my coworker has a feature branch named "stale_feature" behind main
    And I am on the "main" branch
    When I run `git prune-branches`


  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND                        |
      | main   | git fetch --prune              |
      |        | git push origin :stale_feature |
    And I end up on the "main" branch
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
      |        | git push -u origin stale_feature                     |
    And I end up on the "main" branch
    Then the existing branches are
      | REPOSITORY | BRANCHES            |
      | local      | main, stale_feature |
      | remote     | main, stale_feature |
      | coworker   | main, stale_feature |
