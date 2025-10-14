Feature: enable offline mode

  Background:
    Given a Git repo with origin

  Rule: enabling stores the offline status in the global Git settings

    Scenario Outline: enable via CLI
      When I run "git-town offline <GIVE>"
      Then global Git setting "git-town.offline" is now "<WANT>"

      Examples:
        | GIVE | WANT |
        | true | true |
        | t    | true |
        | 1    | true |
        | on   | true |
        | yes  | true |

  Rule: disabling stores the offline status in the global Git settings

    Scenario Outline: disable via CLI
      And offline mode is enabled
      When I run "git-town offline <GIVE>"
      Then global Git setting "git-town.offline" is now "<WANT>"

      Examples:
        | GIVE  | WANT  |
        | false | false |
        | f     | false |
        | 0     | false |
        | off   | false |
        | no    | false |

  Rule: undo removes the config setting

    Scenario: undo
      And I ran "git-town offline"
      When I run "git-town undo"
      Then global Git setting "git-town.offline" now doesn't exist

  Rule: does not accept invalid values

    Scenario: provide invalid value via CLI
      When I run "git-town offline zonk"
      Then Git Town prints the error:
        """
        invalid value for git-town.offline: "zonk". Please provide either "yes" or "no"
        """
      And global Git setting "git-town.offline" still doesn't exist
