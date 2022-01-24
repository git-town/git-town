Feature: passing an invalid option to the pull strategy configuration


  Scenario: default setting
    When I run "git-town pull-branch-strategy"
    Then it prints:
      """
      rebase
      """

  Scenario: explicit rebase
    Given the pull-branch-strategy configuration is "rebase"
    When I run "git-town pull-branch-strategy"
    Then it prints:
      """
      rebase
      """

  Scenario: explicit merge
    Given the pull-branch-strategy configuration is "merge"
    When I run "git-town pull-branch-strategy"
    Then it prints:
      """
      merge
      """
