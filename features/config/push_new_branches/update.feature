Feature: set push-new-branches

  Scenario Outline: local setting
    When I run "git-town config push-new-branches <GIVE>"
    Then setting "push-new-branches" is now "<WANT>"

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
    When I run "git-town config push-new-branches zonk"
    Then it prints the error:
      """
      invalid argument: "zonk". Please provide either "yes" or "no"
      """

  Scenario Outline: global setting
    When I run "git-town config push-new-branches --global <GIVE>"
    Then setting "push-new-branches" is now "<WANT>"

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
