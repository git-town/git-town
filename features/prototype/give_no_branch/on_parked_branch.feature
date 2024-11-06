Feature: prototype the current parked branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE   | PARENT | LOCATIONS |
      | parked | parked | main   | local     |
    And the current branch is "parked"
    When I run "git-town prototype"

  Scenario: result
    Then Git Town runs no commands
    And it prints:
      """
      branch "parked" is now a prototype branch
      """
    And the current branch is still "parked"
    And branch "parked" is now prototype
    And branch "parked" is still parked

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "parked"
    And branch "parked" is still parked
    And there are now no prototype branches
