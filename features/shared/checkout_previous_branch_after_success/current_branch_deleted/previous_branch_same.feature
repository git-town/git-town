Feature: Git checkout history is preserved when deleting the current branch


  Scenario: kill
    Given my repo has the feature branches "previous" and "current"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town kill"
    Then I am now on the "main" branch
    And the previous Git branch is still "previous"

  Scenario: prune-branches
    Given my repo has the feature branches "previous" and "current"
    And the "current" branch gets deleted on the remote
    And the following commits exist in my repo
      | BRANCH   | LOCATION | FILE NAME     | FILE CONTENT     |
      | previous | local    | previous_file | previous content |
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town prune-branches"
    Then I am now on the "main" branch
    And the previous Git branch is still "previous"

  Scenario: ship
    Given my repo has the feature branches "previous" and "current"
    And the following commits exist in my repo
      | BRANCH  | LOCATION | FILE NAME    | FILE CONTENT    |
      | current | remote   | feature_file | feature content |
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town ship -m 'feature done'"
    Then I am now on the "main" branch
    And the previous Git branch is still "previous"
