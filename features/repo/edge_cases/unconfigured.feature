@skipWindows
Feature: ask for missing configuration

  Scenario: unconfigured
    Given Git Town is not configured
    And the origin is "https://github.com/git-town/git-town.git"
    And the "open" tool is installed
    When I run "git-town repo" and answer the prompts:
      | PROMPT                                     | ANSWER  |
      | Please specify the main development branch | [ENTER] |
    And the main branch is now "main"
