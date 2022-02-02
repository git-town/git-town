Feature: ask for missing configuration

  @skipWindows
  Scenario: run unconfigured
    Given I haven't configured Git Town yet
    And my repo's origin is "https://github.com/git-town/git-town.git"
    And my computer has the "open" tool installed
    When I run "git-town new-pull-request" and answer the prompts:
      | PROMPT                                     | ANSWER  |
      | Please specify the main development branch | [ENTER] |
    Then it prints the initial configuration prompt
    And the main branch is now "main"
    And my repo is now has no perennial branches
    And "open" launches a new pull request with this url in my browser:
      """
      https://github.com/git-town/git-town.github.com/compare/feature?expand=1 |
      """
