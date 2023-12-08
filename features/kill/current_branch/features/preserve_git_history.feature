Feature: preserve the previous Git branch

  Background:
    Given the feature branches "previous" and "current"
    And the current branch is "current" and the previous branch is "previous"

  Scenario: previous branch remains
    When I run "git-town kill"
    Then the current branch is now "previous"
    And the previous Git branch is now ""

  Scenario: previous branch is gone
    When I run "git-town kill previous"
    Then the current branch is still "current"
    And the previous Git branch is now ""

  Scenario: current and previous branch remain
    Given a feature branch "victim"
    And the current branch is "current" and the previous branch is "previous"
    When I run "git-town kill victim"
    Then the current branch is still "current"
    And the previous Git branch is still "previous"
