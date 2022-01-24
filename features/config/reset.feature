Feature: resetting the configuration

  Scenario: everything is configured
    Given the main branch is configured as "main"
    And the perennial branches are configured as "qa" and "staging"
    When I run "git-town config reset"
    Then Git Town is no longer configured for this repo


  Scenario: the main branch is configured but the perennial branches are not
    Given the main branch is configured as "main"
    And the perennial branches are not configured
    When I run "git-town config reset"
    Then Git Town is no longer configured for this repo


  Scenario: the main branch is not configured but the perennial branches are
    Given the main branch name is not configured
    And the perennial branches are configured as "qa"
    When I run "git-town config reset"
    Then Git Town is no longer configured for this repo


  Scenario: nothing is configured yet
    Given I haven't configured Git Town yet
    When I run "git-town config reset"
    Then Git Town is no longer configured for this repo
