Feature: set the global new-branch-push-flag

  Scenario Outline:
    When I run "git-town new-branch-push-flag --global <GIVE>"
    Then the new-branch-push-flag configuration is now <WANT>

    Examples:
      | GIVE  | WANT  |
      | true  | true  |
      | t     | true  |
      | 1     | true  |
      | false | false |
      | f     | false |
      | 0     | false |
