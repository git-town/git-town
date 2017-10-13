Feature: set the hack-push flag

  As a user or tool configuring Git Town
  I want an easy way to specifically set the git-hack push flag
  So that I can configure Git Town safely, and the tool does exactly what I want.


  Scenario: update to "true"
    When I run `git-town hack-push-flag true`
    Then my repo is now configured with "hack-push-flag" set to "true"


  Scenario: update to "false"
    When I run `git-town hack-push-flag false`
    Then my repo is now configured with "hack-push-flag" set to "false"


  Scenario: invalid flag
    When I run `git-town hack-push-flag woof`
    Then Git Town prints the error "Invalid value: 'woof'"
    And it prints the error:
      """
      Usage:
        git-town hack-push-flag [(true | false)] [flags]
      """
