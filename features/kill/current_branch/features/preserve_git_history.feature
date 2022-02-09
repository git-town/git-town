Feature: preserve the previous Git branch

  Background:
    Given the feature branches "previous" and "current"
    And the current branch is "current" and the previous branch is "previous"

  Scenario: previous branch exists
    When I run "git-town kill"
    Then the current branch is now "main"
    And the previous Git branch is still "previous"

  Scenario: previous branch gone
    When I run "git-town kill previous"
    Then the current branch is still "current"
    And the previous Git branch is now "main"

  Scenario: current and previous branch exist
    Given a feature branch "victim"
    And the current branch is "current" and the previous branch is "previous"
    When I run "git-town kill victim"
    Then the current branch is still "current"
    And the previous Git branch is still "previous"
