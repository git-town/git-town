Feature: git town hack: Recognize already existing branches

  Scenario: branch exists locally
    Given my repo has a feature branch named "existing"
    When I run "git-town hack existing"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      a branch named "existing" already exists
      """

  Scenario: branch exists remotely
    Given my coworker has a feature branch named "existing-feature"
    And I am on the "main" branch
    When I run "git-town hack existing-feature"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      a branch named "existing-feature" already exists
      """
