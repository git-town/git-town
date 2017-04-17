Feature: too many arguments

  As a developer providing too many arguments
  I should be reminded of how many arguments the command expects
  So that I can use it correctly without having to look that fact up in the readme.


  Scenario: hack
    When I run `gt hack arg1 arg2`
    Then it runs no commands
    And I get the error "Too many arguments"
    And I get the error
      """
      Usage:
        gt hack <branch> [flags]
      """


  Scenario: hack-push-flag
    When I run `gt hack-push-flag arg1 arg2`
    Then I get the error "Too many arguments"
    And I get the error
      """
      Usage:
        gt hack-push-flag [(true | false)] [flags]
      """


  Scenario: kill
    When I run `gt kill arg1 arg2`
    Then it runs no commands
    And I get the error "Too many arguments"
    And I get the error
      """
      Usage:
        gt kill [<branch>] [flags]
      """


  Scenario: main-branch
    When I run `gt main-branch arg1 arg2`
    Then I get the error "Too many arguments"
    And I get the error
      """
      Usage:
        gt main-branch [<branch>]
      """


  Scenario: new-pull-request
    When I run `gt new-pull-request arg1`
    Then it runs no commands
    And I get the error "Too many arguments"
    And I get the error
      """
      Usage:
        gt new-pull-request [flags]
      """


  Scenario: perennial-branches
    When I run `gt perennial-branches arg1`
    Then I get the error "Too many arguments"
    And I get the error
      """
      Usage:
        gt perennial-branches [flags]
      """


  Scenario: prune-branches
    When I run `gt prune-branches arg1`
    Then it runs no commands
    And I get the error "Too many arguments"
    And I get the error
      """
      Usage:
        gt prune-branches [flags]
      """


  Scenario: pull-branch-strategy
    When I run `gt pull-branch-strategy arg1 arg2`
    Then I get the error "Too many arguments"
    And I get the error
      """
      Usage:
        gt pull-branch-strategy [(rebase | merge)] [flags]
      """


  Scenario: repo
    When I run `gt repo arg1`
    Then it runs no commands
    And I get the error "Too many arguments"
    And I get the error
      """
      Usage:
        gt repo [flags]
      """


  Scenario: sync
    When I run `gt sync arg1`
    Then it runs no commands
    And I get the error "Too many arguments"
    And I get the error
      """
      Usage:
        gt sync [flags]
      """
