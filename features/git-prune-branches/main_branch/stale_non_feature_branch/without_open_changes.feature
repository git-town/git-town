Feature: git prune-branches: don't remove stale non-feature branches when called from the main branch (without open changes)

  (see ./with_open_changes.feature)


  Background:
    Given I have a non-feature branch "production" behind main
    And I am on the "main" branch
    When I run `git prune-branches`


  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND           |
      | main   | git fetch --prune |
    Then I end up on the "main" branch
    And the existing branches are
      | REPOSITORY | BRANCHES         |
      | local      | main, production |
      | remote     | main, production |
      | coworker   | main             |
