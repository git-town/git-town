Feature: change offline mode

  Scenario Outline: valid settings
    When I run "git-town config offline <GIVE>"
    Then setting "offline" is now "<WANT>"

    Examples:
      | GIVE  | WANT  |
      | true  | true  |
      | t     | true  |
      | 1     | true  |
      | on    | true  |
      | yes   | true  |
      | false | false |
      | f     | false |
      | 0     | false |
      | off   | false |
      | no    | false |

  Scenario: invalid value
    Given setting "offline" is "false"
    When I run "git-town config offline zonk"
    Then it prints the error:
      """
      invalid argument: "zonk". Please provide either "yes" or "no"
      """
    And setting "offline" is still "false"
