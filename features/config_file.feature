Feature: Configuration File

  Scenario: Without a configuration file
    Given I don't have a configuration file
    When I run `git ship` and enter "user_main_branch"
    Then a the main branch name is configured as "user_main_branch"


  Scenario: With an old configuration file
    Given I have an old configuration file with main branch: "user_main_branch"
    When I run `git ship` while allowing errors
    Then I end up with a new configuration file with main branch: "user_main_branch"
    And I don't have an old configuration file anymore
