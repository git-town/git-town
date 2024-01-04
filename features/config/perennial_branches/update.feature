@skipWindows
Feature: change the perennial branches

  Background:
    Given the branches "staging" and "qa"

  Scenario: add a perennial branch
    Given the perennial branches are "qa"
    When I run "git-town config perennial-branches update" and answer the prompts:
      | PROMPT                            | ANSWER               |
      | Please specify perennial branches | [DOWN][SPACE][ENTER] |
    Then the perennial branches are now "qa" and "staging"

  Scenario: remove a perennial branch
    Given the perennial branches are "staging" and "qa"
    When I run "git-town config perennial-branches update" and answer the prompts:
      | PROMPT                            | ANSWER               |
      | Please specify perennial branches | [DOWN][SPACE][ENTER] |
    Then the perennial branches are now "qa"
