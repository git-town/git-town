Feature: update the push-hook setting

  Scenario Outline: changing the local Git setting
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

  Scenario Outline: changing the global Git setting
    When I run "git-town config push-hook <GIVE> --global"
    Then global Git Town setting "push-hook" is now "<WANT>"

    Examples:
      | GIVE  | WANT  |
      | true  | true  |
      | false | false |

  @this
  Scenario Outline: changing the configuration file
    And the configuration file:
      """
      push-hook = false
      """
    When I run "git-town config push-hook true"
    Then the configuration file is now:
      """
      push-hook = true
      """

  Scenario: setting to an invalid value
    When I run "git-town config push-hook zonk"
    Then it prints the error:
      """
      invalid argument: "zonk". Please provide either "yes" or "no"
      """
