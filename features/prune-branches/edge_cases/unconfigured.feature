@skipWindows
Feature: ask for missing configuration

  Scenario: run unconfigured
    Given I haven't configured Git Town yet
    When I run "git-town prune-branches" and answer the prompts:
      | PROMPT                                     | ANSWER  |
      | Please specify the main development branch | [ENTER] |
    And the main branch is now "main"
    And my repo now has no perennial branches
