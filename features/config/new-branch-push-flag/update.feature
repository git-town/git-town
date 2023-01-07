Feature: set the new-branch-push-flag

  Scenario Outline: local setting
    When I run "git-town config new-branch-push-flag <GIVE>"
    Then setting "new-branch-push-flag" is now "<WANT>"

    Examples:
      | GIVE  | WANT  |
      | true  | true  |
      | t     | true  |
      | 1     | true  |
      | false | false |
      | f     | false |
      | 0     | false |

  @this
  Scenario: invalid value
    When I run "git-town config new-branch-push-flag zonk"
    Then it prints the error:
      """
      invalid argument: "zonk". Please provide either "yes" or "no"
      """

  Scenario Outline: global setting
    When I run "git-town config new-branch-push-flag --global <GIVE>"
    Then setting "new-branch-push-flag" is now "<WANT>"

    Examples:
      | GIVE  | WANT  |
      | true  | true  |
      | t     | true  |
      | 1     | true  |
      | false | false |
      | f     | false |
      | 0     | false |
