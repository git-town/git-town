Feature: does not rename the main branch

  Background:
    Given a Git repo with origin
    And the current branch is "main"

  Scenario: try to rename
    When I run "git-town rename main new"
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And Git Town prints the error:
      """
      the main branch cannot be renamed
      """
    And the current branch is still "main"

  Scenario: try to force rename
    When I run "git-town rename main new --force"
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And Git Town prints the error:
      """
      the main branch cannot be renamed
      """
    And the current branch is still "main"
