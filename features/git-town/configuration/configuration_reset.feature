Feature: resetting the configuration

  As a user no longer using Git Town on a repository
  I want to be able to cleanly remove all Git Town configuration from my Git repo
  So that my repository is left in a clean state after the uninstallation.


  Scenario: everything is configured
    Given I have configured the main branch name as "main"
    And my non-feature branches are configured as "qa" and "staging"
    When I run `git town config --reset`
    Then Git Town is no longer configured for this repository


  Scenario: the main branch is configured but the non-feature branches are not
    Given I have configured the main branch name as "main"
    And my non-feature branches are not configured
    When I run `git town config --reset`
    Then Git Town is no longer configured for this repository


  Scenario: the main branch is not configured but the non-feature branches are
    Given I don't have a main branch name configured
    And my non-feature branches are configured as "qa"
    When I run `git town config --reset`
    Then Git Town is no longer configured for this repository


  Scenario: nothing is configured yet
    Given I haven't configured Git Town yet
    When I run `git town config --reset`
    Then Git Town is no longer configured for this repository
