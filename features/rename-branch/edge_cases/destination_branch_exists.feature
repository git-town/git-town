Feature: destination branch exists

  Scenario: destination branch exists locally
    Given my repo has the feature branches "alpha" and "beta"
    And my repo contains the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, remote | alpha commit |
      | beta   | local, remote | beta commit  |
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

  Scenario: destination branch exists remotely
    Given my repo has a feature branch "alpha"
    And a coworker has a feature branch "beta"
    And my repo contains the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, remote | alpha commit |
      | beta   | remote        | beta commit  |
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
