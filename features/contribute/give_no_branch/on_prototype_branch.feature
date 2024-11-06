Feature: make the current prototype branch a contribution branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | prototype | prototype | main   | local, origin |
    And the current branch is "prototype"
    When I run "git-town contribute"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "prototype" is now a contribution branch
      """
    And branch "prototype" is now a contribution branch
    And the current branch is still "prototype"
    And there are now no prototype branches

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "prototype"
    And branch "prototype" is now prototype
    And there are now no contribution branches
