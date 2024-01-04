Feature: update the push-hook setting

  Scenario Outline: changing an existing local setting
    When I run "git-town config push-hook <GIVE>"
    Then local Git Town setting "push-hook" is now "<WANT>"

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

  Scenario Outline: changing the global setting
    When I run "git-town config push-hook <GIVE> --global"
    Then global Git Town setting "push-hook" is now "<WANT>"

    Examples:
      | GIVE | WANT  |
      | yes  | true  |
      | no   | false |

  Scenario: invalid value
    When I run "git-town config push-hook zonk"
    Then it prints the error:
      """
      invalid argument: "zonk". Please provide either "yes" or "no"
      """
