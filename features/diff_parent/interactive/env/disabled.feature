@skipWindows
Feature: no TTY, missing main branch

  Scenario: main branch
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE   | PARENT | LOCATIONS |
      | feature | (none) |        | local     |
    And Git Town is not configured
    And the current branch is "feature"
    When I run "git-town diff-parent" with these environment variables
      | GIT_TOWN_INTERACTIVE | false |
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      Error: no main branch configured and interactivity disabled via environment variable.
      
      To configure:
      git config git-town.main-branch <branch>
      """
