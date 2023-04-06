Feature: configure the sync-strategy

  Scenario Outline:
    When I run "git-town config sync-strategy <VALUE>"
    Then setting "sync-strategy" is now "<VALUE>"

    Examples:
      | VALUE  |
      | rebase |
      | merge  |
