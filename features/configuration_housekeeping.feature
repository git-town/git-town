Feature: Configuration housekeeping

  Scenario: Automatic update of old configuration files
    Given I have an old configuration file with main branch: "main"
    When I run `git town-hack new-feature`
    Then my repo is configured with the main branch as "main"
    And I don't have an old configuration file anymore
