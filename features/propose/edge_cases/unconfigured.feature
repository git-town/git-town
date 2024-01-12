Feature: ask for missing configuration

  Scenario: run unconfigured
    Given Git Town is not configured
    And the origin is "https://github.com/git-town/git-town.git"
    And tool "open" is installed
    When I run "git-town propose" and enter into the dialog:
      | DIALOG                  | KEYS  |
      | main development branch | enter |
    And the main branch is now "main"
    And "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town.github.com/compare/feature?expand=1 |
      """
