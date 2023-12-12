@skipWindows
Feature: ask for missing configuration

  Scenario: run unconfigured
    Given Git Town is not configured
    And the origin is "https://github.com/git-town/git-town.git"
    And tool "open" is installed
    When I run "git-town propose" and answer the prompts:
      | PROMPT                                     | ANSWER  |
      | Please specify the main development branch | [ENTER] |
    And the main branch is now "main"
    And "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town.github.com/compare/feature?expand=1 |
      """
