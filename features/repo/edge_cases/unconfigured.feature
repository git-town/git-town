@skipWindows
Feature: ask for missing configuration

  Scenario: unconfigured
    Given Git Town is not configured
    And the origin is "https://github.com/git-town/git-town.git"
    And tool "open" is installed
    When I run "git-town repo" and enter into the dialog:
      | DIALOG                                     | KEYS  |
      | Please specify the main development branch | enter |
    And the main branch is now "main"
