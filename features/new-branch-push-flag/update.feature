Feature: set the new-branch-push-flag

  Scenario Outline: local setting
    When I run "git-town new-branch-push-flag <GIVE>"
    Then the "new-branch-push-flag" setting is now "<WANT>"

    Examples:
      | GIVE  | WANT  |
      | true  | true  |
      | t     | true  |
      | 1     | true  |
      | false | false |
      | f     | false |
      | 0     | false |

  Scenario: invalid value
    When I run "git-town new-branch-push-flag zonk"
    Then it prints the error:
      """
      invalid argument: "zonk". Please provide either "true" or "false"
      """

  Scenario Outline: global setting
    When I run "git-town new-branch-push-flag --global <GIVE>"
    Then the "new-branch-push-flag" setting is now "<WANT>"

    Examples:
      | GIVE  | WANT  |
      | true  | true  |
      | t     | true  |
      | 1     | true  |
      | false | false |
      | f     | false |
      | 0     | false |
