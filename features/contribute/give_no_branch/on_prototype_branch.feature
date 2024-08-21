Feature: make the current prototype branch a contribution branch

  Background:
    Given a Git repo with origin
    And the branch
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | prototype | prototype | main   | local, origin |
    And the current branch is "prototype"
    When I run "git-town contribute"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "prototype" is now a contribution branch
      """
    And branch "prototype" is now a contribution branch
    And the current branch is still "prototype"
    And there are now no prototype branches

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "prototype"
    And branch "prototype" is now prototype
    And there are now no contribution branches
