Feature: Commmand pattern architecture


  Scenario: Preconditions are called prior to running steps
    Given I have a command pattern script
    And it has the following preconditions:
      | PRECONDITIONS                       |
      | echo "echo precondition_runs_first" |
    And it has the following steps:
      | STEPS              |
      | echo "echo step_1" |
      | echo "echo step_2" |
    When I run the command pattern script
    Then I see
      """
      precondition_runs_first
      step_1
      step_2
      """

  Scenario: Steps file is cleaned up
    Given I have a command pattern script
    And it has the following preconditions:
      | PRECONDITIONS                     |
      | echo "echo git-mock precondition" |
    And it has the following steps:
      | STEPS              |
      | echo "echo step_1" |
      | echo "echo step_2" |
    When I run the command pattern script
    Then the steps file for the command pattern script is removed
