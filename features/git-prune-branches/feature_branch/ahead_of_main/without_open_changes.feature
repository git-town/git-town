Feature: git prune-branches: keep used feature branches when run on a feature branch (without open changes)

  (see ./with_open_changes.feature)


  Background:
    Given I have a feature branch named "feature" ahead of main
    And I am on the "feature" branch
    When I run `git prune-branches`


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND              |
      | feature | git fetch --prune    |
      |         | git checkout main    |
      | main    | git checkout feature |
    And I end up on the "feature" branch
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | remote     | main, feature |
      | coworker   | main          |


  Scenario: undoing the operation
    When I run `git prune-branches --undo`
    Then it runs the Git commands
      | BRANCH  | COMMAND              |
      | feature | git checkout main    |
      | main    | git checkout feature |
    And I end up on the "feature" branch
    Then the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | remote     | main, feature |
      | coworker   | main          |

