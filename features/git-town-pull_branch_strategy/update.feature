Feature: set the pull branch strategy

  As a user or tool configuring Git Town
  I want an easy way to specifically set the pull branch strategy
  So that I can configure Git Town safely, and the tool does exactly what I want.


  Scenario: update to merge
    When I run `git-town pull-branch-strategy merge`
    Then my repo is now configured with "pull-branch-strategy" set to "merge"


  Scenario: update to rebase
    When I run `git-town pull-branch-strategy rebase`
    Then my repo is now configured with "pull-branch-strategy" set to "rebase"


  Scenario: invalid strategy
    When I run `git-town pull-branch-strategy woof`
    Then Git Town prints the error "Invalid value: 'woof'"
    And it prints the error:
      """
      Usage:
        git-town pull-branch-strategy [(rebase | merge)] [flags]
      """
