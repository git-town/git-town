Feature: Configuration housekeeping

  Scenario: Without a configured main branch name
    Given I don't have a main branch name configured
    When I run `git hack new-feature` and enter "main"
    Then the main branch name is now configured as "main"


  Scenario: Automatic update of old configuration files
    Given I have an old configuration file with main branch: "main"
    When I run `git hack new-feature`
    Then the main branch name is now configured as "main"
    And I don't have an old configuration file anymore
