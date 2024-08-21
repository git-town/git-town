Feature: make another observed branch a contribution branch

  Background:
    Given a Git repo with origin
    And the branch
      | NAME     | TYPE     | PARENT | LOCATIONS     |
      | observed | observed |        | local, origin |
    And an uncommitted file
    When I run "git-town contribute observed"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "observed" is now a contribution branch
      """
    And the contribution branches are now "observed"
    And there are now no observed branches
    And the current branch is still "main"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND       |
      | main   | git add -A    |
      |        | git stash     |
      |        | git stash pop |
    And the current branch is still "main"
    And the observed branches are now "observed"
    And the uncommitted file still exists
