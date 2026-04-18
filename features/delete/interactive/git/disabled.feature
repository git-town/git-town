@skipWindows
Feature: interactivity disabled, no main branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE   | PARENT | LOCATIONS     |
      | branch-1 | (none) |        | local, origin |
    And the current branch is "branch-1"
    And Git Town is not configured
    And Git setting "git-town.interactive" is "false"
    When I run "git-town delete"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                  |
      | branch-1 | git fetch --prune --tags |
    And Git Town prints the error:
      """
      no main branch configured and interactivity disabled via Git metadata.
      
      To configure:
      git config git-town.main-branch <branch>
      """
