Feature: already existing branch

  Scenario: the branch to create already exists locally
    Given the current branch is a feature branch "old"
    And a local feature branch "existing"
    When I run "git-town prepend existing"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | old    | git fetch --prune --tags |
    And it prints the error:
      """
      there is already a branch "existing"
      """

  Scenario: the branch to create already exists at the origin remote
    Given the current branch is a feature branch "old"
    And a remote feature branch "existing"
    When I run "git-town prepend existing"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | old    | git fetch --prune --tags |
    And it prints the error:
      """
      there is already a branch "existing" at the "origin" remote
      """
