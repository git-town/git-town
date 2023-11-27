Feature: configure the sync-feature-strategy

  Scenario Outline:
    When I run "git-town config sync-feature-strategy <VALUE>"
    Then Git Town setting "sync-feature-strategy" is now "<VALUE>"

    Examples:
      | VALUE  |
      | rebase |
      | merge  |
