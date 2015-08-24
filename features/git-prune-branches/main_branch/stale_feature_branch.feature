Feature: git prune-branches: remove stale feature branches when run on the main branch

  (see ../../feature_branch/behind_main/without_open_changes.feature)


  Background:
    Given I have a feature branch named "stale-feature" behind main
    And I am on the "main" branch
    And I have an uncommitted file
    When I run `git prune-branches`


  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                        |
      | main   | git fetch --prune              |
      |        | git push origin :stale-feature |
      |        | git branch -d stale-feature    |
    And I end up on the "main" branch
    And I still have my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES |
      | local      | main     |
      | remote     | main     |


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
    And I still have my uncommitted file
