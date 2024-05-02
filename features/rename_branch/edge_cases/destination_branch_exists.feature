Feature: destination branch exists

  Scenario: destination branch exists locally
    Given the feature branches "alpha" and "beta"
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, origin | alpha commit |
      | beta   | local, origin | beta commit  |
    And the current branch is "alpha"
    When I run "git-town rename-branch alpha beta"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | alpha  | git fetch --prune --tags |
    And it prints the error:
      """
      there is already a branch "beta"
      """
    And the current branch is still "alpha"
    And the initial branches and lineage exist

  Scenario: destination branch exists in origin
    Given the current branch is a feature branch "alpha"
    And a remote branch "beta"
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, origin | alpha commit |
      | beta   | origin        | beta commit  |
    When I run "git-town rename-branch alpha beta"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | alpha  | git fetch --prune --tags |
    And it prints the error:
      """
      there is already a branch "beta" at the "origin" remote
      """
    And the current branch is still "alpha"
    And the initial branches and lineage exist
