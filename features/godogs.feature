Feature: eat lots of godogs

  Scenario: eat 5 out of 12
    Given there are 12 godogs
    When I eat 5
    Then there should be 7 remaining
