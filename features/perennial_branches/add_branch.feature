Feature: add a branch to the perennial branches configuration

  As a user or tool configuring Git Town's perennial branches
  I want an easy way to add a branch to my set of perennial branches
  So that I can configure Git Town safely, and the tool does exactly what I want.

  Background:
    Given my repo has the branches "staging" and "qa"
    And the perennial branches are configured as "qa"

  @skipWindows
  Scenario: adding a branch
    When I run "git-town perennial-branches update" and answer the prompts:
      | PROMPT                            | ANSWER               |
      | Please specify perennial branches | [DOWN][SPACE][ENTER] |
    Then the perennial branches are now configured as "qa" and "staging"
