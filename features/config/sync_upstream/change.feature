Feature: set sync-upstream

  Scenario Outline: local setting in Git metadata
    When I run "git-town config sync-upstream <GIVE>"
    Then local Git Town setting "sync-upstream" is now "<WANT>"

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
    When I run "git-town config sync-upstream zonk"
    Then it prints the error:
      """
      invalid argument: "zonk". Please provide either "yes" or "no"
      """

  Scenario Outline: global setting
    When I run "git-town config sync-upstream --global <GIVE>"
    Then global Git Town setting "sync-upstream" is now "<WANT>"

    Examples:
      | GIVE  | WANT  |
      | true  | true  |
      | false | false |

  Scenario: add to empty config file
    Given the configuration file:
      """
      """
    When I run "git-town config sync-upstream yes"
    Then the configuration file is now:
      """
      sync-upstream = true
      """
    And local Git Town setting "sync-upstream" still doesn't exist

  Scenario: change existing value in config file
    Given the configuration file:
      """
      sync-upstream = true
      """
    When I run "git-town config sync-upstream no"
    Then the configuration file is now:
      """
      sync-upstream = false
      """
    And local Git Town setting "sync-upstream" still doesn't exist
