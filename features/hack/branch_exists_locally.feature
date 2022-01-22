Feature: git town-hack: errors when the branch exists locally

  To ensure unique feature branches
  When trying to create a branch with the name of an existing local branch
  I want to see guidance.

  Scenario:
    Given my repo has a feature branch named "existing-feature"
    When I run "git-town hack existing-feature"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      a branch named "existing-feature" already exists
      """
