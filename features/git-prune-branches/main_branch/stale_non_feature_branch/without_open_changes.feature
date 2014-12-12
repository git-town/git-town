Feature: git-prune-branches: does not remove stale non-feature branches when called from the main branch without open changes

  As a developer having empty non-feature branches in my repository
  I want them all to survive a branch pruning
  So that I can keep my repository clean without messing up my deployment infrastructure.


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
