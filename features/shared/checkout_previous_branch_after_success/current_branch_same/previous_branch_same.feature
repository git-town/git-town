Feature: Git checkout history is preserved when the current and previous branch don't change

  Scenario: kill
    Given my repo has the feature branches "previous" and "current"
    And my repo has a feature branch named "victim"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town kill victim"
    Then I am still on the "current" branch
    And the previous Git branch is still "previous"

  Scenario: new-pull-request
    Given my repo has the feature branches "previous" and "current"
    And my computer has the "open" tool installed
    And my repo's origin is "https://github.com/git-town/git-town.git"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town new-pull-request"
    Then I am still on the "current" branch
    And the previous Git branch is still "previous"

  Scenario: prune-branches
    Given my repo has the feature branches "previous" and "current"
    And the following commits exist in my repo
      | BRANCH   | LOCATION | FILE NAME     | FILE CONTENT     |
      | previous | local    | previous_file | previous content |
      | current  | local    | current_file  | current content  |
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town prune-branches"
    Then I am still on the "current" branch
    And the previous Git branch is still "previous"

  Scenario: repo
    Given my repo has the feature branches "previous" and "current"
    And my computer has the "open" tool installed
    And my repo's origin is "https://github.com/git-town/git-town.git"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town repo"
    Then I am still on the "current" branch
    And the previous Git branch is still "previous"

  Scenario: ship
    Given my repo has the feature branches "previous" and "current"
    And my repo has a feature branch named "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION | FILE NAME    | FILE CONTENT    |
      | feature | remote   | feature_file | feature content |
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town ship feature -m "feature done""
    Then I am still on the "current" branch
    And the previous Git branch is still "previous"

  Scenario: sync
    Given my repo has the feature branches "previous" and "current"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town sync"
    Then I am still on the "current" branch
    And the previous Git branch is still "previous"
