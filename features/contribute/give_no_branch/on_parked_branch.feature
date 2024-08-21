Feature: make the current parked branch a contribution branch

  Background:
    Given a Git repo with origin
    And the branch
      | NAME   | TYPE   | PARENT | LOCATIONS     |
      | parked | parked | main   | local, origin |
    And the current branch is "parked"
    When I run "git-town contribute"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "parked" is now a contribution branch
      """
    And branch "parked" is now a contribution branch
    And the current branch is still "parked"
    And there are now no parked branches

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "parked"
    And branch "parked" is now parked
    And there are now no contribution branches
