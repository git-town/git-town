Feature: require a target branch when staging all open changes

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
    And the current branch is "branch-2"
    And an uncommitted file "changes" with content "my changes"
    When I run "git-town commit --all -m up-commit"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      please provide the branch to commit into
      """
