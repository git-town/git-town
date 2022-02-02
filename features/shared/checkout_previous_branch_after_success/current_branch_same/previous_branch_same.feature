Feature: Git checkout history is preserved when the current and previous branch don't change

  Background:
    Given my repo has the feature branches "previous" and "current"

  Scenario: kill
    Given my repo has a feature branch "victim"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town kill victim"
    Then I am still on the "current" branch
    And the previous Git branch is still "previous"

  Scenario: new-pull-request
    Given my computer has the "open" tool installed
    And my repo's origin is "https://github.com/git-town/git-town.git"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town new-pull-request"
    Then I am still on the "current" branch
    And the previous Git branch is still "previous"

  Scenario: prune-branches
    Given I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town prune-branches"
    Then I am still on the "current" branch
    And the previous Git branch is still "previous"

  Scenario: repo
    Given my computer has the "open" tool installed
    And my repo's origin is "https://github.com/git-town/git-town.git"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town repo"
    Then I am still on the "current" branch
    And the previous Git branch is still "previous"

  Scenario: ship
    Given my repo has a feature branch "feature"
    And my repo contains the commits
      | BRANCH  | LOCATION |
      | feature | remote   |
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town ship feature -m "feature done""
    Then I am still on the "current" branch
    And the previous Git branch is still "previous"

  Scenario: sync
    Given I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town sync"
    Then I am still on the "current" branch
    And the previous Git branch is still "previous"
