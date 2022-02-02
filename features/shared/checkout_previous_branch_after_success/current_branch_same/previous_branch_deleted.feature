Feature: deleting the current branch makes the main branch the new previous branch

  Scenario: kill
    Given my repo has the feature branches "previous" and "current"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town kill previous"
    Then I am still on the "current" branch
    And the previous Git branch is now "main"

  Scenario: prune-branches
    Given my repo has the feature branches "previous" and "current"
    And the "previous" branch gets deleted on the remote
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town prune-branches"
    Then I am still on the "current" branch
    And the previous Git branch is now "main"

  Scenario: ship
    Given my repo has the feature branches "previous" and "current"
    And my repo contains the commits
      | BRANCH   | LOCATION |
      | previous | local    |
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town ship previous -m 'feature done'"
    Then I am still on the "current" branch
    And the previous Git branch is now "main"
