Feature: Ask for missing configuration

  As a user having forgotten to configure Git Town
  I want to be prompted to configure it when I use it the first time
  So that I use a properly configured tool at all times.

  @skipWindows
  Scenario: run unconfigured
    Given I haven't configured Git Town yet
    And my repo's origin is "https://github.com/git-town/git-town.git"
    And my computer has the "open" tool installed
    When I run "git-town new-pull-request" and answer the prompts:
      | PROMPT                                     | ANSWER  |
      | Please specify the main development branch | [ENTER] |
    Then it prints the initial configuration prompt
    And the main branch is now configured as "main"
    And my repo is now configured with no perennial branches
    And "open" launches a new pull request with this url in my browser:
      """
      https://github.com/git-town/git-town.github.com/compare/feature?expand=1 |
      """
