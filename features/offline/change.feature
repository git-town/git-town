Feature: change offline mode

  Background:
    Given a Git repo clone

  Scenario Outline: writing to local Git metadata
    When I run "git-town offline <GIVE>"
    Then global Git Town setting "offline" is now "<WANT>"

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
    And global Git Town setting "offline" is "false"
    When I run "git-town offline zonk"
    Then it prints the error:
      """
      invalid value for git-town.offline: "zonk". Please provide either "yes" or "no"
      """
    And global Git Town setting "offline" is still "false"
