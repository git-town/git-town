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

  Scenario: changing an existing entry in the configuration file
    Given the configuration file:
      """
      push-hook = true
      """
    When I run "git-town config push-hook false"
    Then the configuration file is now:
      """
      push-hook = false
      """

  Scenario: creating an entry in the configuration file
    Given the configuration file:
      """
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
