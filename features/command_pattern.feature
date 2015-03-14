Feature: Commmand pattern architecture

  Scenario: Preconditions are called prior to running steps
    Given I have a command pattern script
    And it has the following preconditions:
      | PRECONDITIONS                       |
      | echo "echo precondition runs first" |
    And it has the following steps:
      | STEPS              |
      | echo "echo step 1" |
      | echo "echo step 2" |
    When I run the command pattern script
    Then I see
      """
      precondition runs first
      step 1
      step 2
      """
