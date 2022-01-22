Feature: git town-hack: errors when the branch exists remotely

  To ensure unique feature branches
  When trying to create a branch with the name of an existing remote branch
  I want to see guidance.

  Scenario:
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
