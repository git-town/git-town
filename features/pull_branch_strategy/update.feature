Feature: set the pull branch strategy

  Scenario: update to merge
    When I run "git-town pull-branch-strategy merge"
    Then the pull-branch-strategy configuration is now "merge"

  Scenario: update to rebase
    When I run "git-town pull-branch-strategy rebase"
    Then the pull-branch-strategy configuration is now "rebase"
