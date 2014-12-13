Feature: show help screen when Git Town is not configured

  Background:
    Given I haven't configured Git Town yet


  Scenario: git town with no flags
    When I run `git town`
    Then I see "Git Town is a collection of additional Git commands"


  Scenario: git town, configured, with "help" subcommand
    When I run `git town help`
    Then I see "Git Town is a collection of additional Git commands"
