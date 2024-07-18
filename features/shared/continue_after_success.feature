Feature: continue after successful command

  Scenario Outline:
    Given a Git repo clone
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And I run "git-town <COMMAND>"
    When I run "git-town continue"
    Then it prints:
      """
      nothing to continue
      """

    Examples:
      | COMMAND                 |
      |                         |
      | append new              |
      | completions fish        |
      | config                  |
      | diff-parent             |
      | hack new                |
      | help                    |
      | kill feature            |
      | main_branch             |
      | offline                 |
      | perennial-branches      |
      | prepend new             |
      | propose                 |
      | sync-perennial-strategy |
      | push-new-branches       |
      | rename-branch           |
      | repo                    |
      | ship feature -m done    |
      | sync                    |
      | version                 |
