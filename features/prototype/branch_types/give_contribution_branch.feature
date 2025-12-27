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
      branch contribution is now a prototype branch
      """
    And branch "contribution" now has type "prototype"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "contribution" now has type "contribution"
