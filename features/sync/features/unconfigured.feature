Feature: ask for missing configuration information

  Scenario: run unconfigured
    Given Git Town is not configured
    When I run "git-town sync" and enter into the dialog:
      | DIALOG                  | KEYS  |
      | main development branch | enter |
    And the main branch is now "main"
