Feature: set the pull branch strategy

  As a user or tool configuring Git Town
  I want an easy way to specifically set the pull branch strategy
  So that I can configure Git Town safely, and the tool does exactly what I want.

  Scenario: update to merge
    When I run "git-town pull-branch-strategy merge"
    Then the pull-branch-strategy configuration is now "merge"

  Scenario: update to rebase
    When I run "git-town pull-branch-strategy rebase"
    Then the pull-branch-strategy configuration is now "rebase"
