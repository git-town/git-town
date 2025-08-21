Feature: already existing local branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | local     |
    And the current branch is "main"
    When I run "git-town append existing"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And Git Town prints the error:
      """
      there is already a branch "existing"
      """
# The last "append" command didn't leave runstate that could be undone.
# When the user runs "undo", it would undo the Git Town command that executed before.
