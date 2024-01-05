Feature: set sync-before-ship

  Scenario Outline: local setting in Git metadata
    When I run "git-town config sync-before-ship <GIVE>"
    Then local Git Town setting "sync-before-ship" is now "<WANT>"

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

  Scenario: invalid value in Git metadata
    When I run "git-town config sync-before-ship zonk"
    Then it prints the error:
      """
      invalid argument: "zonk". Please provide either "yes" or "no"
      """

  Scenario Outline: global setting in Git metadata
    When I run "git-town config sync-before-ship --global <GIVE>"
    Then global Git Town setting "sync-before-ship" is now "<WANT>"

    Examples:
      | GIVE  | WANT  |
      | true  | true  |
      | false | false |

  Scenario: add to empty config file
    Given the configuration file:
      """
      """
    When I run "git-town config sync-before-ship yes"
    Then the configuration file is now:
      """
      sync-before-ship = true
      """
    And local Git Town setting "sync-before-ship" still doesn't exist

  Scenario: change existing value in config file
    Given the configuration file:
      """
      sync-before-ship = true
      """
    When I run "git-town config sync-before-ship no"
    Then the configuration file is now:
      """
      sync-before-ship = false
      """
    And local Git Town setting "sync-before-ship" still doesn't exist
