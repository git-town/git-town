Feature: make a branch non-perennial

  Background:
    Given the branches "staging" and "qa"

  Scenario: existing branch, non-existing configuration
    Given local Git Town setting "perennial-branches" doesn't exist
    When I run "git-town config perennial-branches remove staging"
    Then it prints the error:
      """
      branch "staging" is not perennial
      """
    And local Git Town setting "perennial-branches" still doesn't exist

  Scenario: remove an existing branch from existing Git configuration
    Given the perennial branches are "qa" and "staging"
    When I run "git-town config perennial-branches remove staging"
    Then the perennial branches are now "qa"

  Scenario: remove an existing branch from an empty config file
    Given local Git Town setting "perennial-branches" doesn't exist
    And the configuration file:
      """
      """
    When I run "git-town config perennial-branches remove staging"
    Then it prints the error:
      """
      branch "staging" is not perennial
      """
    And the configuration file is still:
      """
      """
    And local Git Town setting "perennial-branches" still doesn't exist

  Scenario: remove an existing perennial branch from the config file
    Given local Git Town setting "perennial-branches" doesn't exist
    And the configuration file:
      """
      [branches]
        perennials = ["qa", "staging"]
      """
    When I run "git-town config perennial-branches remove qa"
    And the configuration file is now:
      """
      [branches]
        perennials = ["staging"]
      """
    And local Git Town setting "perennial-branches" still doesn't exist

  Scenario: remove a non-existing branch
    Given the perennial branches are "qa"
    When I run "git-town config perennial-branches remove zonk"
    Then it prints the error:
      """
      branch "zonk" does not exist
      """
    And the perennial branches are still "qa"
