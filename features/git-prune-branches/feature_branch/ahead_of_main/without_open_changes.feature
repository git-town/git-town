Feature: git prune-branches: keep used feature branches when run on a feature branch (without open changes)

  (see ./with_open_changes.feature)


  Background:
    Given I have a feature branch named "feature" ahead of main
    And I have a feature branch named "stale_feature" behind main
    And I am on the "feature" branch
    When I run `git prune-branches`


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND                        |
      | feature | git fetch --prune              |
      | feature | git checkout main              |
      | main    | git push origin :stale_feature |
      | main    | git branch -d stale_feature    |
      | main    | git checkout feature           |
    Then I end up on the "feature" branch
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | remote     | main, feature |
      | coworker   | main          |
