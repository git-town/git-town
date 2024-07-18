Feature: does not rename the main branch

  Background:
    Given a Git repo clone
    And the current branch is "main"

  Scenario: try to rename
    When I run "git-town rename-branch main new"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      the main branch cannot be renamed
      """
    And the current branch is still "main"

  Scenario: try to force rename
    When I run "git-town rename-branch main new --force"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      the main branch cannot be renamed
      """
    And the current branch is still "main"
