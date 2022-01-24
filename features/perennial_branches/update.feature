@skipWindows
Feature: changing the perennial branches

  Background:
    Given my repo has the branches "staging" and "qa"

  Scenario: add a perennial branch
    Given the perennial branches are configured as "qa"
    When I run "git-town perennial-branches update" and answer the prompts:
      | PROMPT                            | ANSWER               |
      | Please specify perennial branches | [DOWN][SPACE][ENTER] |
    Then the perennial branches are now configured as "qa" and "staging"

  Scenario: remove a perennial branch
    Given the perennial branches are configured as "staging" and "qa"
    When I run "git-town perennial-branches update" and answer the prompts:
      | PROMPT                            | ANSWER               |
      | Please specify perennial branches | [DOWN][SPACE][ENTER] |
    Then the perennial branches are now configured as "qa"
