Feature: ask for missing configuration information

  Scenario: unconfigured
    Given a Git repo with origin
    And Git Town is not configured
    And Git Town setting "ship-strategy" is "squash-merge"
    When I run "git-town ship" and enter into the dialog:
      | DIALOG      | KEYS  |
      | main branch | enter |
    And the main branch is now "main"
    And it prints the error:
      """
      cannot ship the main branch
      """
