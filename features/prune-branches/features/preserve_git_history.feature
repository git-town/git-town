Feature: preserve the previous Git branch

  Background:
    Given my repo has the feature branches "previous" and "current"
    And I am on the "current" branch with "previous" as the previous Git branch

  Scenario: current branch gone, previous branch exists
    And the "current" branch gets deleted at origin
    When I run "git-town prune-branches"
    Then I am now on the "main" branch
    And the previous Git branch is still "previous"

  Scenario: current branch exists, previous branch gone
    Given the "previous" branch gets deleted at origin
    When I run "git-town prune-branches"
    Then I am still on the "current" branch
    And the previous Git branch is now "main"

  Scenario: both branches deleted
    And the "previous" branch gets deleted at origin
    And the "current" branch gets deleted at origin
    When I run "git-town prune-branches"
    Then I am now on the "main" branch
    And the previous Git branch is now "main"

  Scenario: both branches exist
    When I run "git-town prune-branches"
    Then I am still on the "current" branch
    And the previous Git branch is still "previous"
