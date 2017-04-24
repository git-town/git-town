Feature: listing the configuration

  As a user running the Git Town configuration wizard,
  I want to see the existing configuration values
  So that I can change it more effectively


  Background:
    Given I have branches named "production" and "qa"


  Scenario: everything is configured
    Given I have configured the main branch name as "main"
    And my perennial branches are configured as "qa"
    When I run `gt config --setup` and enter "main" and ""
    Then I see
      """
      Git Town needs to be configured

        1: main
        2: production
        3: qa

      Please specify the main development branch by name or number (current value: main):
      """
    And I see
      """
      Please specify a perennial branch by name or number. Leave it blank to finish (current value: qa):
      """


  Scenario: empty input
    Given I have configured the main branch name as "main"
    And my perennial branches are configured as "qa"
    When I run `gt config --setup` and enter "", "main" and ""
    Then I see "A main development branch is required to enable the features provided by Git Town"
    And my repo is configured with the main branch as "main"
    And my repo is configured with no perennial branches


  Scenario: non-empty input
    Given I have configured the main branch name as "main"
    And my perennial branches are configured as "qa"
    When I run `gt config --setup` and enter:
      | main       |
      | production |
      |            |
    And my repo is configured with the main branch as "main"
    And my repo is configured with perennial branches as "production"
