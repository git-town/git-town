Feature: change the perennial branches

  Background:
    Given the branches "staging" and "qa"

  Scenario: add a perennial branch when no configuration exists
    Given local Git Town setting "perennial-branches" doesn't exist
    When I run "git-town config perennial-branches change" and enter into the dialog:
      | DIALOG             | KEYS             |
      | perennial branches | down space enter |
    Then the perennial branches are now "staging"

  Scenario: add a perennial branch to existing local Git configuration
    Given local Git Town setting "perennial-branches" is "staging"
    When I run "git-town config perennial-branches change" and enter into the dialog:
      | DIALOG             | KEYS        |
      | perennial branches | space enter |
    Then the perennial branches are now "qa" and "staging"

  Scenario: remove a perennial branch from existing Git configuration
    Given the perennial branches are "staging" and "qa"
    When I run "git-town config perennial-branches change" and enter into the dialog:
      | DIALOG             | KEYS        |
      | perennial branches | space enter |
    Then the perennial branches are now "staging"

  Scenario: add a perennial branch to an empty config file
    Given local Git Town setting "perennial-branches" doesn't exist
    And the configuration file:
      """
      """
    When I run "git-town config perennial-branches change" and enter into the dialog:
      | DIALOG             | KEYS        |
      | perennial branches | space enter |
    And the configuration file is now:
      """
      [branches]
        perennials = ["qa"]
      """
    And local Git Town setting "perennial-branches" still doesn't exist

  Scenario: add a perennial branch to already existing config file entries
    Given local Git Town setting "perennial-branches" doesn't exist
    And the configuration file:
      """
      [branches]
        perennials = ["staging"]
      """
    When I run "git-town config perennial-branches change" and enter into the dialog:
      | DIALOG             | KEYS        |
      | perennial branches | space enter |
    And the configuration file is now:
      """
      [branches]
        perennials = ["qa", "staging"]
      """
    And local Git Town setting "perennial-branches" still doesn't exist

  Scenario: remove a perennial branch from existing config file entries
    Given local Git Town setting "perennial-branches" doesn't exist
    And the configuration file:
      """
      [branches]
        perennials = ["qa", "staging"]
      """
    When I run "git-town config perennial-branches change" and enter into the dialog:
      | DIALOG             | KEYS        |
      | perennial branches | space enter |
    And the configuration file is now:
      """
      [branches]
        perennials = ["staging"]
      """
    And local Git Town setting "perennial-branches" still doesn't exist
