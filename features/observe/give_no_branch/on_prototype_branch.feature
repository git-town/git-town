Feature: making a prototype branch an observed branch

  Background:
    Given a Git repo with origin
    And the branch
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | prototype | prototype | main   | local, origin |
    And the current branch is "prototype"
    When I run "git-town observe"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "prototype" is now an observed branch
      """
    And the observed branches are now "prototype"
    And there are now no prototype branches

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And there are now no observed branches
    And the prototype branches are now "prototype"
