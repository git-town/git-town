Feature: git prune-branches: don't remove stale perennial branches when called from the main branch

  As a developer pruning branches
  I want perennial branches to not be deleted
  So that I can keep my repository clean without messing up my deployment infrastructure.


  Background:
    Given I have a perennial branch "production" behind main
    And I am on the "main" branch
    And I have an uncommitted file
    When I run `git prune-branches`


  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND           |
      | main   | git fetch --prune |
    And I end up on the "main" branch
    And I still have my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES         |
      | local      | main, production |
      | remote     | main, production |
      | coworker   | main             |


  Scenario: undoing the operation
    When I run `git prune-branches --undo`
    Then I get the error "Nothing to undo"
    And it runs no Git commands
    And I end up on the "main" branch
    Then the existing branches are
      | REPOSITORY | BRANCHES         |
      | local      | main, production |
      | remote     | main, production |
    And I still have my uncommitted file
