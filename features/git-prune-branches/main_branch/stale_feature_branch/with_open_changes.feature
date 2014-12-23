Feature: git prune-branches: remove stale feature branches when run on the main branch (with open changes)

  As a developer pruning branches on the main branch
  I want all stale branches to be deleted
  So that all remaining branches are relevant, and I can focus on my current work.


  Background:
    Given I have a feature branch named "stale_feature" behind main
    And I am on the "main" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git prune-branches`


  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND                        |
      | main   | git fetch --prune              |
      | main   | git push origin :stale_feature |
      | main   | git branch -d stale_feature    |
    And I end up on the "main" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | REPOSITORY | BRANCHES |
      | local      | main     |
      | remote     | main     |
      | coworker   | main     |
