Feature: show help screen


  Scenario: called with no parameters
    When I run `git town`
    Then I see "Git Town is a collection of additional Git commands"


  Scenario: called with "help" subcommand
    When I run `git town help`
    Then I see "Git Town is a collection of additional Git commands"
