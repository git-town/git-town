Feature: already existing local branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | local     |
    When I run "git-town append existing"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      there is already a branch "existing"
      """

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is now "main"
    And the initial commits exist
    And the initial branches and lineage exist
