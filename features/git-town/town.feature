Feature: Show correct git town usage

  Scenario: invalid git town command
    When I run `git town invalidcommand`
    Then I see "'invalidcommand' is not a valid Git Town command"
    And I see "usage: git town"
