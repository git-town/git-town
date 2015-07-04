Feature: git prune-branches: remove stale feature branches when run on the main branch

  (see ../../feature_branch/behind_main/without_open_changes.feature)


  Background:
    Given I have a feature branch named "stale_feature" behind main
    And I am on the "main" branch
    And I have an uncommitted file
    When I run `git prune-branches`


  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND                        |
      | main   | git fetch --prune              |
      |        | git push origin :stale_feature |
      |        | git branch -d stale_feature    |
    And I end up on the "main" branch
    And I still have my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES |
      | local      | main     |
      | remote     | main     |


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
    And I still have my uncommitted file
