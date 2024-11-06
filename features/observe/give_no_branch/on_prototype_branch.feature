Feature: observe the current prototype branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | prototype | prototype | main   | local, origin |
    And the current branch is "prototype"
    When I run "git-town observe"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "prototype" is now an observed branch
      """
    And the observed branches are now "prototype"
    And there are now no prototype branches

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And there are now no observed branches
    And the prototype branches are now "prototype"
