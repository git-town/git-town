Feature: set the hack-push strategy

  As a user or tool configuring Git Town
  I want an easy way to specifically set the git-hack push strategy 
  So that I can configure Git Town safely, and the tool does exactly what I want.


  Scenario: update to push
    When I run `git town hack-push-strategy push`
    Then my repo is now configured with "hack-push-strategy" set to "push"


  Scenario: update to local
    When I run `git town hack-push-strategy local`
    Then my repo is now configured with "hack-push-strategy" set to "local"


  Scenario: invalid strategy
    When I run `git town hack-push-strategy woof`
    Then I see
      """
      Invalid git-hack push strategy: 'woof'.
      Valid git-hack push strategies are 'push' and 'local'.
      """
