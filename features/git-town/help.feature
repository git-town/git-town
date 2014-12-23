Feature: git town: show help screen

  As a user not knowing certain specifics about Git Town
  I want to be able to get help easily
  So that I can refresh my memory quickly, move on to what I actually wanted to do, and remain efficient.


  Scenario: called with no parameters
    When I run `git town`
    Then I see "Git Town is a collection of additional Git commands"


  Scenario: called with "help" subcommand
    When I run `git town help`
    Then I see "Git Town is a collection of additional Git commands"
