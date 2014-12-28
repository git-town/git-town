Feature: git prune-branches: remove stale feature branches when run on a feature branch (without open changes)

  As a developer pruning branches
  I want all merged branches to be deleted
  So that my remaining branches are relevant and I can focus on my current work.


  Background:
    Given I have a feature branch named "feature" behind main
    And I have a feature branch named "stale_feature" behind main
    And I am on the "feature" branch
    When I run `git prune-branches`


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND                        |
      | feature | git fetch --prune              |
      | feature | git checkout main              |
      | main    | git push origin :feature       |
      | main    | git branch -d feature          |
      | main    | git push origin :stale_feature |
      | main    | git branch -d stale_feature    |
    And I end up on the "main" branch
    And the existing branches are
      | REPOSITORY | BRANCHES |
      | local      | main     |
      | remote     | main     |
      | coworker   | main     |
