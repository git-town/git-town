Feature: resetting the configuration

  As a user no longer using Git Town on a repository
  I want to be able to cleanly remove all Git Town configuration from my git town-repo
  So that my repository is left in a clean state after the uninstallation.


  Scenario: everything is configured
    Given the main branch is configured as "main"
    And the perennial branches are configured as "qa" and "staging"
    When I run `git-town config reset`
    Then Git Town is no longer configured for this repository


  Scenario: the main branch is configured but the perennial branches are not
    Given the main branch is configured as "main"
    And my perennial branches are not configured
    When I run `git-town config reset`
    Then Git Town is no longer configured for this repository


  Scenario: the main branch is not configured but the perennial branches are
    Given I don't have a main branch name configured
    And the perennial branches are configured as "qa"
    When I run `git-town config reset`
    Then Git Town is no longer configured for this repository


  Scenario: nothing is configured yet
    Given I haven't configured Git Town yet
    When I run `git-town config reset`
    Then Git Town is no longer configured for this repository
