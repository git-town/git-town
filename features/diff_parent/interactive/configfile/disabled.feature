@skipWindows
Feature: no TTY, missing main branch

  Scenario: main branch
    Given a Git repo with origin
    And the committed configuration file:
      """
      interactive = false
      """
    And the branches
      | NAME    | TYPE   | PARENT | LOCATIONS |
      | feature | (none) |        | local     |
    And Git Town is not configured
    And the current branch is "feature"
    When I run "git-town diff-parent"
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      Error: no main branch configured and interactivity disabled via config file.
      
      To configure:
      git config git-town.main-branch <branch>
      """
