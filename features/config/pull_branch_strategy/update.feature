Feature: configure the pull_branch_strategy

  Scenario Outline:
    When I run "git-town config pull-branch-strategy <VALUE>"
    Then Git Town setting "pull-branch-strategy" is now "<VALUE>"

    Examples:
      | VALUE  |
      | rebase |
      | merge  |
