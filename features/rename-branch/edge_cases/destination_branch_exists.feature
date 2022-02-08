Feature: destination branch exists

  Scenario: destination branch exists locally
    Given my repo has the feature branches "alpha" and "beta"
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, origin | alpha commit |
      | beta   | local, origin | beta commit  |
    And I am on the "alpha" branch
    When I run "git-town rename-branch alpha beta"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | alpha  | git fetch --prune --tags |
    And it prints the error:
      """
      a branch named "beta" already exists
      """
    And I am still on the "alpha" branch
    And my repo now has its initial branches and branch hierarchy

  Scenario: destination branch exists in origin
    Given my repo has a feature branch "alpha"
    And the origin has a feature branch "beta"
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, origin | alpha commit |
      | beta   | origin        | beta commit  |
    And I am on the "alpha" branch
    When I run "git-town rename-branch alpha beta"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | alpha  | git fetch --prune --tags |
    And it prints the error:
      """
      a branch named "beta" already exists
      """
    And I am still on the "alpha" branch
    And my repo now has its initial branches and branch hierarchy
