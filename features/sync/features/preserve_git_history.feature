Feature: preserve the previous Git branch

  Background:
    Given the feature branches "previous" and "current"
    And the current branch is "current" and the previous branch is "previous"
    When I run "git-town sync"

  Scenario: result
    Then the current branch is still "current"
    And the previous Git branch is still "previous"

  Scenario: undo
    When I run "git-town undo"
    Then the current branch is still "current"
    And the previous Git branch is now "main"
