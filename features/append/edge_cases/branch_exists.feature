Feature: already existing branch

  Scenario: the branch to create already exists locally
    Given a local feature branch "existing"
    When I run "git-town append existing"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      there is already a branch "existing"
      """

  Scenario: the branch to create already exists at the origin remote
    Given a remote feature branch "existing"
    When I run "git-town append existing"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      there is already a branch "existing" at the "origin" remote
      """
