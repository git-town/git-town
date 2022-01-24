Feature: set the new-branch-push-flag

  As a user or tool configuring Git Town
  I want an easy way to specifically set the new branch push flag
  So that I can configure Git Town safely, and the tool does exactly what I want.

  Scenario Outline: update
    When I run "git-town new-branch-push-flag <GIVE>"
    Then the new-branch-push-flag configuration is now <WANT>

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

  Scenario: multiple arguments
    When I run "git-town new-branch-push-flag true false"
    Then it prints the error:
      """
      accepts at most 1 arg(s), received 2
      """
