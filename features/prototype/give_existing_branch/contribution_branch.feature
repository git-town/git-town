Feature: prototype another contribution branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS     |
      | contribution | contribution | main   | local, origin |
    When I run "git-town prototype contribution"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "contribution" is now a prototype branch
      """
    And the prototype branches are now "contribution"
    And there are now no contribution branches

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the contribution branches are now "contribution"
    And there are now no prototype branches
