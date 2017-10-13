Feature: too many arguments

  As a developer providing too many arguments
  I should be reminded of how many arguments the command expects
  So that I can use it correctly without having to look that fact up in the readme.


  Scenario: hack
    When I run `git-town hack arg1 arg2`
    Then Git Town runs no commands
    And it prints the error "Too many arguments"
    And it prints the error:
      """
      Usage:
        git-town hack <branch> [flags]
      """


  Scenario: hack-push-flag
    When I run `git-town hack-push-flag arg1 arg2`
    Then Git Town prints the error "Too many arguments"
    And it prints the error:
      """
      Usage:
        git-town hack-push-flag [(true | false)] [flags]
      """


  Scenario: kill
    When I run `git-town kill arg1 arg2`
    Then Git Town runs no commands
    And it prints the error "Too many arguments"
    And it prints the error:
      """
      Usage:
        git-town kill [<branch>] [flags]
      """


  Scenario: main-branch
    When I run `git-town main-branch arg1 arg2`
    Then Git Town prints the error "Too many arguments"
    And it prints the error:
      """
      Usage:
        git-town main-branch [<branch>]
      """


  Scenario: new-pull-request
    When I run `git-town new-pull-request arg1`
    Then Git Town runs no commands
    And it prints the error "Too many arguments"
    And it prints the error:
      """
      Usage:
        git-town new-pull-request [flags]
      """


  Scenario: perennial-branches
    When I run `git-town perennial-branches arg1`
    Then Git Town prints the error "Too many arguments"
    And it prints the error:
      """
      Usage:
        git-town perennial-branches [flags]
      """


  Scenario: prune-branches
    When I run `git-town prune-branches arg1`
    Then Git Town runs no commands
    And it prints the error "Too many arguments"
    And it prints the error:
      """
      Usage:
        git-town prune-branches [flags]
      """


  Scenario: pull-branch-strategy
    When I run `git-town pull-branch-strategy arg1 arg2`
    Then Git Town prints the error "Too many arguments"
    And it prints the error:
      """
      Usage:
        git-town pull-branch-strategy [(rebase | merge)] [flags]
      """


  Scenario: repo
    When I run `git-town repo arg1`
    Then Git Town runs no commands
    And it prints the error "Too many arguments"
    And it prints the error:
      """
      Usage:
        git-town repo [flags]
      """


  Scenario: sync
    When I run `git-town sync arg1`
    Then Git Town runs no commands
    And it prints the error "Too many arguments"
    And it prints the error:
      """
      Usage:
        git-town sync [flags]
      """
