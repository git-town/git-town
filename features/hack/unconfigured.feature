Feature: Ask for missing configuration

  To ensure the hack command finishes
  When configuration information is missing
  I want to have a chance to enter the missing configuration data.

  @skipWindows
  Scenario: running unconfigured
    Given I haven't configured Git Town yet
    When I run "git-town hack feature" and answer the prompts:
      | PROMPT                                     | ANSWER  |
      | Please specify the main development branch | [ENTER] |
    Then it prints the initial configuration prompt
    And it runs the commands
      | BRANCH  | COMMAND                    |
      | main    | git fetch --prune --tags   |
      |         | git rebase origin/main     |
      |         | git branch feature main    |
      |         | git checkout feature       |
    And the main branch is now configured as "main"
    And I am now on the "feature" branch
    And Git Town is now aware of this branch hierarchy
      | BRANCH  | PARENT |
      | feature | main   |
