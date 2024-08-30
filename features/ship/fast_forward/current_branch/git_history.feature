Feature: preserve the previous Git branch when shipping using the fast-forward strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | previous | feature | main   | local     |
      | current  | feature | main   | local     |
    And Git Town setting "ship-strategy" is "fast-forward"

  Scenario: current branch gone
    Given the commits
      | BRANCH  | LOCATION |
      | current | local    |
    And the current branch is "current" and the previous branch is "previous"
    When I run "git-town ship"
    Then the current branch is now "main"
    And the previous Git branch is now "previous"

  Scenario: previous branch gone
    Given the commits
      | BRANCH   | LOCATION |
      | previous | local    |
    And the current branch is "current" and the previous branch is "previous"
    When I run "git-town ship previous"
    Then the current branch is still "current"
    And the previous Git branch is now "main"

  Scenario: both branches exist
    Given the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the commits
      | BRANCH  | LOCATION |
      | feature | local    |
    And the current branch is "current" and the previous branch is "previous"
    When I run "git-town ship feature"
    Then the current branch is still "current"
    And the previous Git branch is still "previous"
