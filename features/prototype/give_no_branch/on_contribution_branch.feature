Feature: prototype the current contribution branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | LOCATIONS |
      | contribution | contribution | local     |
    And the current branch is "contribution"
    When I run "git-town prototype"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "contribution" is now a prototype branch
      """
    And the current branch is still "contribution"
    And branch "contribution" is now prototype
    And there are now no contribution branches

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "contribution"
    And the initial branches and lineage exist
    And branch "contribution" is now a contribution branch
    And there are now no prototype branches
