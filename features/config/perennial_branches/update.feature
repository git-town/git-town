@skipWindows
Feature: change the perennial branches

  Background:
    Given the branches "staging" and "qa"

  Scenario: add a perennial branch to existing Git configuration
    Given the perennial branches are "qa"
    When I run "git-town config perennial-branches update" and answer the prompts:
      | PROMPT                            | ANSWER               |
      | Please specify perennial branches | [DOWN][SPACE][ENTER] |
    Then the perennial branches are now "qa" and "staging"

  Scenario: remove a perennial branch from existing Git configuration
    Given the perennial branches are "staging" and "qa"
    When I run "git-town config perennial-branches update" and answer the prompts:
      | PROMPT                            | ANSWER               |
      | Please specify perennial branches | [DOWN][SPACE][ENTER] |
    Then the perennial branches are now "qa"

  Scenario: add perennial branches to an empty config file
    Given local Git Town setting "perennial-branches" doesn't exist
    And the configuration file:
      """
      """
    When I run "git-town config perennial-branches update" and answer the prompts:
      | PROMPT                            | ANSWER               |
      | Please specify perennial branches | [DOWN][SPACE][ENTER] |
    And the configuration file is now:
      """
      [branches]
        perennials = ["staging"]
      """
    And local Git Town setting "perennial-branches" still doesn't exist

  Scenario: add perennial branches to already existing entries in the config file
    Given local Git Town setting "perennial-branches" doesn't exist
    And the configuration file:
      """
      [branches]
        perennials = ["qa"]
      """
    When I run "git-town config perennial-branches update" and answer the prompts:
      | PROMPT                            | ANSWER               |
      | Please specify perennial branches | [DOWN][SPACE][ENTER] |
    And the configuration file is now:
      """
      [branches]
        perennials = ["qa", "staging"]
      """
    And local Git Town setting "perennial-branches" still doesn't exist
