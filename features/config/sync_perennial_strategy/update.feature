Feature: configure the sync_perennial_strategy

  Scenario Outline:
    When I run "git-town config sync-perennial-strategy <VALUE>"
    Then Git Town setting "sync-perennial-strategy" is now "<VALUE>"

    Examples:
      | VALUE  |
      | rebase |
      | merge  |
