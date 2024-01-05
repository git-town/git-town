Feature: set push-new-branches

  Scenario Outline: configured in local Git metadata
    When I run "git-town config push-new-branches <GIVE>"
    Then local Git Town setting "push-new-branches" is now "<WANT>"

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
    When I run "git-town config push-new-branches zonk"
    Then it prints the error:
      """
      invalid argument: "zonk". Please provide either "yes" or "no"
      """

  Scenario Outline: configured in global Git metadata
    When I run "git-town config push-new-branches --global <GIVE>"
    Then global Git Town setting "push-new-branches" is now "<WANT>"

    Examples:
      | GIVE  | WANT  |
      | true  | true  |
      | false | false |

  Scenario: add to empty config file
    Given the configuration file:
      """
      """
    When I run "git-town config push-new-branches yes"
    Then the configuration file is now:
      """
      push-new-branches = true
      """
    And local Git Town setting "push-new-branches" still doesn't exist

  Scenario: change existing value in config file
    Given the configuration file:
      """
      push-new-branches = true
      """
    When I run "git-town config push-new-branches no"
    Then the configuration file is now:
      """
      push-new-branches = false
      """
    And local Git Town setting "push-new-branches" still doesn't exist
