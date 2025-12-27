Feature: switching when the branch has no parent

  Background:
    Given a Git repo with origin
    And the current branch is "main"
    When I run "git-town down"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND |
    And Git Town prints the error:
      """
      branch main has no parent
      """
