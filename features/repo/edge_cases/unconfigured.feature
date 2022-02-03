@skipWindows
Feature: ask for missing configuration

  Scenario: unconfigured
    Given I haven't configured Git Town yet
    And my repo's origin is "https://github.com/git-town/git-town.git"
    And my computer has the "open" tool installed
    When I run "git-town repo" and answer the prompts:
      | PROMPT                                     | ANSWER  |
      | Please specify the main development branch | [ENTER] |
    And the main branch is now "main"
    And my repo is now has no perennial branches
