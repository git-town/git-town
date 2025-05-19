Feature: walk each branch of a stack in dry-run mode

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS |
      | branch-1 | feature | main     | local     |
      | branch-2 | feature | branch-1 | local     |
    And the current branch is "branch-2"
    When I run "git-town walk --all --dry-run"

  Scenario: result
    Then Git Town prints the error:
      """
      Error: cannot dry-run manual iteration through all branches, please call with a shell command to run on each branch
      """
