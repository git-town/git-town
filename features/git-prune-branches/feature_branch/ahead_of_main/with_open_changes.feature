Feature: git prune-branches: keep used feature branches when run on a feature branch (without open changes)

  As a developer pruning branches
  I want my feature branches with commits to not be deleted
  So that I can keep my repository clean without losing work.


  Background:
    Given I have a feature branch named "feature" ahead of main
    And I am on the "feature" branch
    And I have an uncommitted file
    When I run `git prune-branches`


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND              |
      | feature | git fetch --prune    |
      |         | git stash -u         |
      |         | git checkout main    |
      | main    | git checkout feature |
      | feature | git stash pop        |
    And I end up on the "feature" branch
    And I still have my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | remote     | main, feature |
      | coworker   | main          |


  Scenario: undoing the operation
    When I run `git prune-branches --undo`
    Then it runs the Git commands
      | BRANCH  | COMMAND              |
      | feature | git stash -u         |
      |         | git checkout main    |
      | main    | git checkout feature |
      | feature | git stash pop        |
    And I end up on the "feature" branch
    And I still have my uncommitted file
    Then the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | remote     | main, feature |
      | coworker   | main          |
