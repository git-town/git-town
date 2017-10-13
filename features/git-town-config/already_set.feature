Feature: listing the configuration

  As a user running the Git Town configuration wizard,
  I want to see the existing configuration values
  So that I can change it more effectively


  Background:
    Given my repository has branches named "production" and "qa"


  Scenario: everything is configured
    Given Git Town's main branch is configured as "main"
    And its perennial branches are configured as "qa"
    When I run `git-town config --setup` and enter "main" and ""
    Then Git Town prints
      """
      Git Town needs to be configured

        1: main
        2: production
        3: qa

      Please specify the main development branch by name or number (current value: main):
      """
    And Git Town prints
      """
      Please specify a perennial branch by name or number. Leave it blank to finish (current value: qa):
      """


  Scenario: empty input
    Given Git Town's main branch is configured as "main"
    And its perennial branches are configured as "qa"
    When I run `git-town config --setup` and enter "", "main" and ""
    Then Git Town prints "A main development branch is required to enable the features provided by Git Town"
    And Git Town's main branch is now configured as "main"
    And my repo is configured with no perennial branches


  Scenario: non-empty input
    Given Git Town's main branch is configured as "main"
    And its perennial branches are configured as "qa"
    When I run `git-town config --setup` and enter:
      | main       |
      | production |
      |            |
    And Git Town's main branch is now configured as "main"
    And Git Town's perennial branches are now configured as "production"
