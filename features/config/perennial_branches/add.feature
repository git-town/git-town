Feature: add perennial branches

  Background:
    Given the branches "staging" and "qa"

  Scenario: add an existing branch to non-existing configuration
    Given local Git Town setting "perennial-branches" doesn't exist
    When I run "git-town config perennial-branches add staging"
    Then the perennial branches are now "staging"

  Scenario: add an existing branch to existing Git configuration
    Given the perennial branches are "qa"
    When I run "git-town config perennial-branches add staging"
    Then the perennial branches are now "qa" and "staging"

  Scenario: add an existing branch to an empty config file
    Given local Git Town setting "perennial-branches" doesn't exist
    And the configuration file:
      """
      """
    When I run "git-town config perennial-branches add staging"
    And the configuration file is now:
      """
      [branches]
        perennials = ["staging"]
      """
    And local Git Town setting "perennial-branches" still doesn't exist

  Scenario: add an existing branch to existing perennial branches in the config file
    Given local Git Town setting "perennial-branches" doesn't exist
    And the configuration file:
      """
      [branches]
        perennials = ["staging"]
      """
    When I run "git-town config perennial-branches add qa"
    And the configuration file is now:
      """
      [branches]
        perennials = ["qa", "staging"]
      """
    And local Git Town setting "perennial-branches" still doesn't exist

  Scenario: add a non-existing branch
    Given the perennial branches are "qa"
    When I run "git-town config perennial-branches add zonk"
    Then it prints the error:
      """
      branch "zonk" does not exist
      """
    Then the perennial branches are still "qa"
