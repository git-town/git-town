Feature: ask for missing configuration

  Scenario:
    Given Git Town is not configured
    When I run "git-town kill" and enter into the dialog:
      | DIALOG                  | KEYS  |
      | main development branch | enter |
    Then it prints the error:
      """
      you cannot kill the main branch
      """
    And the main branch is now "main"
