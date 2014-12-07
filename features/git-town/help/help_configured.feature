Feature: show help screen when Git Town is configured

  Background:
    Given I have configured the main branch name as "main"
    And my non-feature branches are configured as "qa, staging"


  Scenario: git town with no flags
    When I run `git town`
    Then I see "Git Town is a collection of additional Git commands"


  Scenario: git town, configured, with "help" subcommand
    When I run `git town help`
    Then I see "Git Town is a collection of additional Git commands"
