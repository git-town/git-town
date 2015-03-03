Feature: git prune-branches: don't remove used feature branches when called on the main branch (with open changes)

  As a developer pruning branches
  I want my feature branches with commits to not be deleted
  So that I can keep my branches organized without losing work.


  Background:
    Given I have a feature branch named "my-feature" ahead of main
    And my coworker has a feature branch named "co-feature" ahead of main
    And I am on the "main" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git prune-branches`


  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND           |
      | main   | git fetch --prune |
    And I end up on the "main" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | REPOSITORY | BRANCHES                     |
      | local      | main, my-feature             |
      | remote     | main, my-feature, co-feature |
      | coworker   | main, co-feature             |


  Scenario: undoing the operation
    When I run `git prune-branches --undo`
    Then I get the error "Cannot undo"
    And it runs no Git commands
    And I am still on the "main" branch
    Then the existing branches are
      | REPOSITORY | BRANCHES                     |
      | local      | main, my-feature             |
      | remote     | main, my-feature, co-feature |
      | coworker   | main, co-feature             |
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
