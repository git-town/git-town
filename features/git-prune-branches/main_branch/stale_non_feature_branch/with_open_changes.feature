Feature: git prune-branches: don't remove stale non-feature branches when called from the main branch (with open changes)

  As a developer pruning branches
  I want non-feature branches to not be deleted
  So that I can keep my repository clean without messing up my deployment infrastructure.


  Background:
    Given I have a non-feature branch "production" behind main
    And I am on the "main" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git prune-branches`


  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND           |
      | main   | git fetch --prune |
    Then I end up on the "main" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | REPOSITORY | BRANCHES         |
      | local      | main, production |
      | remote     | main, production |
      | coworker   | main             |
