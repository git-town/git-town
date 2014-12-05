Feature: show help screen

  Scenario: git town, with no flags
    When I run `git town`
    Then I see "Git Town is a collection of additional Git commands"


  Scenario: git town, with "help" subcommand
    When I run `git town help`
    Then I see "Git Town is a collection of additional Git commands"
