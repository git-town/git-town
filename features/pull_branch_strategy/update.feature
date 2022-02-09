Feature: configure the pull_branch_strategy

  Scenario Outline:
    When I run "git-town pull-branch-strategy <VALUE>"
    Then the "pull-branch-strategy" setting is now "<VALUE>"

    Examples:
      | VALUE  |
      | rebase |
      | merge  |
