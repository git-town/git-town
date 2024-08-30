Feature: make another parked branch an observed branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE   | PARENT | LOCATIONS     |
      | parked | parked | main   | local, origin |
    When I run "git-town observe parked"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "parked" is now an observed branch
      """
    And the observed branches are now "parked"
    And there are now no parked branches

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the parked branches are now "parked"
    And there are now no observed branches
