Feature: set ship-delete-tracking-branch

  Scenario Outline: writing to local Git metadata
    When I run "git-town config ship-delete-tracking-branch <GIVE>"
    Then local Git Town setting "ship-delete-tracking-branch" is now "<WANT>"

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
    When I run "git-town config ship-delete-tracking-branch zonk"
    Then it prints the error:
      """
      invalid argument: "zonk". Please provide either "yes" or "no"
      """

  Scenario Outline: configured in global Git metadata
    When I run "git-town config ship-delete-tracking-branch --global <GIVE>"
    Then global Git Town setting "ship-delete-tracking-branch" is now "<WANT>"

    Examples:
      | GIVE  | WANT  |
      | true  | true  |
      | false | false |

  Scenario: add to empty config file
    Given the configuration file:
      """
      """
    When I run "git-town config ship-delete-tracking-branch yes"
    Then the configuration file is now:
      """
      ship-delete-tracking-branch = true
      """
    And local Git Town setting "ship-delete-tracking-branch" still doesn't exist

  Scenario: change existing value in config file
    Given the configuration file:
      """
      ship-delete-tracking-branch = true
      """
    When I run "git-town config ship-delete-tracking-branch no"
    Then the configuration file is now:
      """
      ship-delete-tracking-branch = false
      """
    And local Git Town setting "ship-delete-tracking-branch" still doesn't exist
