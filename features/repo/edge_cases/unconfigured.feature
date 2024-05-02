@skipWindows
Feature: ask for missing configuration

  @debug @this
  Scenario: unconfigured
    Given Git Town is not configured
    And the origin is "https://github.com/git-town/git-town.git"
    And tool "open" is installed
    And inspect the repo
    When I run "git-town repo" and enter into the dialog:
      | DIALOG            | KEYS  |
      | enter main branch | enter |
    Then "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town
      """
    And the main branch is now "main"
