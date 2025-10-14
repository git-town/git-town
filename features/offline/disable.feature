Feature: disable offline mode

  Background:
    Given a Git repo with origin

  Scenario Outline: disable via CLI
    Given offline mode is enabled
    When I run "git-town offline <VALUE>"
    Then global Git setting "git-town.offline" is now "false"

    Examples:
      | VALUE |
      | false |
      | f     |
      | 0     |
      | off   |
      | no    |

  Scenario: undo
    Given global Git setting "git-town.offline" is "true"
    And I ran "git-town offline 0"
    When I run "git-town undo"
    Then global Git setting "git-town.offline" is now "true"

  Scenario: invalid argument
    When I run "git-town offline zonk"
    Then Git Town prints the error:
      """
      invalid value for git-town.offline: "zonk". Please provide either "yes" or "no"
      """
    And global Git setting "git-town.offline" still doesn't exist
