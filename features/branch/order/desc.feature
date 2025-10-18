Feature: display in reverse order

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | 2024-03 | feature | main   | local, origin |
      | 2025-06 | feature | main   | local, origin |
    And the current branch is "2025-06"
    When I run "git-town branch --order=desc"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
        main
      *   2025-06
          2024-03
      """
