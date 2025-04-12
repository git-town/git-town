Feature: make another contribution branch an observed branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS     |
      | contribution | contribution | main   | local, origin |
    When I run "git-town observe contribution"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "contribution" is now an observed branch
      """
    And the observed branches are now "contribution"
    And branch "contribution" now has type "observed"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the contribution branches are now "contribution"
    And branch "contribution" now has type "contribution"
