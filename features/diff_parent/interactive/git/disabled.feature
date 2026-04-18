@skipWindows
Feature: no TTY, missing main branch

  Scenario: main branch
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE   | PARENT | LOCATIONS |
      | feature | (none) |        | local     |
    And Git Town is not configured
    And the current branch is "feature"
    And Git setting "git-town.interactive" is "false"
    When I run "git-town diff-parent"
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      Error: no main branch configured and interactivity disabled via Git metadata.
      
      To configure:
      git config git-town.main-branch <branch>
      """
