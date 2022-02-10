Feature: configure the pull_branch_strategy

  Scenario Outline:
    When I run "git-town pull-branch-strategy <VALUE>"
    Then setting "pull-branch-strategy" is now "<VALUE>"

    Examples:
      | VALUE  |
      | rebase |
      | merge  |
