Feature: preserve the previous Git branch

  Background:
    Given my repo has the feature branches "previous" and "current"

  Scenario: current branch gone
    And my repo contains the commits
      | BRANCH  | LOCATION |
      | current | local    |
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town ship -m 'feature done'"
    Then I am now on the "main" branch
    And the previous Git branch is still "previous"

  Scenario: previous branch gone
    Given my repo contains the commits
      | BRANCH   | LOCATION |
      | previous | local    |
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town ship previous -m 'feature done'"
    Then I am still on the "current" branch
    And the previous Git branch is now "main"

  Scenario: both branches exist
    Given my repo has a feature branch "feature"
    Given my repo contains the commits
      | BRANCH  | LOCATION |
      | feature | remote   |
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town ship feature -m "feature done""
    Then I am still on the "current" branch
    And the previous Git branch is still "previous"
