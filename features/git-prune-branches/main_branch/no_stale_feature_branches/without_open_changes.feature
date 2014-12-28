Feature: git prune-branches: don't remove used feature branches when called on the main branch (without open changes)

  (see ./with_open_changes.feature)


  Background:
    Given I have a feature branch named "my-feature" ahead of main
    And my coworker has a feature branch named "co-feature" ahead of main
    And I am on the "main" branch
    When I run `git prune-branches`


  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND           |
      | main   | git fetch --prune |
    And I end up on the "main" branch
    And the existing branches are
      | REPOSITORY | BRANCHES                     |
      | local      | main, my-feature             |
      | remote     | main, my-feature, co-feature |
      | coworker   | main, co-feature             |
