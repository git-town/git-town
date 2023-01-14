@this
Feature: update the push-hook setting

  Scenario Outline: local setting
    When I run "git-town config push-hook <GIVE>"
    Then local setting "push-hook" is now "<WANT>"

    Examples:
      | GIVE  | WANT  |
      | true  | true  |
      | yes   | true  |
      | on    | true  |
      | t     | true  |
      | 1     | true  |
      | false | false |
      | no    | false |
      | off   | false |
      | f     | false |
      | 0     | false |

  Scenario: invalid value
    When I run "git-town config push-hook zonk"
    Then it prints the error:
      """
      invalid argument: "zonk". Please provide either "yes" or "no"
      """

  Scenario Outline: global setting
    When I run "git-town config push-hook --global <GIVE>"
    Then global setting "push-hook" is now "<WANT>"

    Examples:
      | GIVE  | WANT  |
      | true  | true  |
      | yes   | true  |
      | on    | true  |
      | t     | true  |
      | 1     | true  |
      | false | false |
      | no    | false |
      | off   | false |
      | f     | false |
      | 0     | false |
