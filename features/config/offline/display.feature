Feature: display the current offline status

  Scenario: default value
    When I run "git-town config offline"
    Then it prints:
      """
      no
      """

  Scenario Outline: valid settings
    Given global setting "offline" is "<VALUE>"
    When I run "git-town config offline"
    Then it prints:
      """
      <OUTPUT>
      """
    Examples:
      | VALUE | OUTPUT |
      | yes   | yes    |
      | on    | yes    |
      | true  | yes    |
      | 1     | yes    |
      | t     | yes    |
      | no    | no     |
      | off   | no     |
      | false | no     |
      | f     | no     |
      | 0     | no     |

  @this
  Scenario: invalid value
    Given global setting "offline" is "zonk"
    When I run "git-town config offline"
    Then it prints the error:
      """
      invalid value for git-town.offline: "zonk". Please provide either "true" or "false"
      """
