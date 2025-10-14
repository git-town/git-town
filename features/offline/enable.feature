Feature: enable offline mode

  Background:
    Given a Git repo with origin

  Scenario Outline: enable via CLI
    When I run "git-town offline <VALUE>"
    Then global Git setting "git-town.offline" is now "true"

    Examples:
      | VALUE |
      | true  |
      | t     |
      | 1     |
      | on    |
      | yes   |

  Scenario: undo
    Given I ran "git-town offline 1"
    When I run "git-town undo"
    Then global Git setting "git-town.offline" now doesn't exist

  Scenario: invalid argument
    When I run "git-town offline zonk"
    Then Git Town prints the error:
      """
      invalid value for git-town.offline: "zonk". Please provide either "yes" or "no"
      """
    And global Git setting "git-town.offline" still doesn't exist
