Feature: change offline mode

  Scenario Outline: valid settings
    When I run "git-town config offline <GIVE>"
    Then global setting "offline" is now "<WANT>"

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
    Given global setting "offline" is "false"
    When I run "git-town config offline zonk"
    Then it prints the error:
      """
      invalid value for git-town.offline: "zonk". Please provide either "yes" or "no"
      """
    And global setting "offline" is still "false"
