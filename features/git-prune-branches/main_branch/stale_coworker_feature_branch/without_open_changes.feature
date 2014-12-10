Feature: git-prune-branches: on the main branch with a stale coworker feature branch without open changes

  Background:
    Given my coworker has a feature branch named "stale_feature" behind main
    And I am on the "main" branch
    When I run `git prune-branches`


  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND                        |
      | main   | git fetch --prune              |
      | main   | git push origin :stale_feature |
    And I end up on the "main" branch
    And the existing branches are
      | REPOSITORY | BRANCHES            |
      | local      | main                |
      | remote     | main                |
      | coworker   | main, stale_feature |
