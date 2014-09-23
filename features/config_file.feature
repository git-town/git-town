Feature: Configuration File Creation

  Scenario: When the config file does not exist
    Given I don't have a configuration file
    When I run `git ship` and enter "user_main_branch"
    Then a the main branch name is configured as "user_main_branch"

