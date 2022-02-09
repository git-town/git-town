@skipWindows
Feature: ask for missing configuration

  Scenario: run unconfigured
    Given Git Town is not configured
    And my repo's origin is "https://github.com/git-town/git-town.git"
    And the "open" tool is installed
    When I run "git-town new-pull-request" and answer the prompts:
      | PROMPT                                     | ANSWER  |
      | Please specify the main development branch | [ENTER] |
    And the main branch is now "main"
    And "open" launches a new pull request with this url in my browser:
      """
      https://github.com/git-town/git-town.github.com/compare/feature?expand=1 |
      """
