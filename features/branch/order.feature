Feature: change the display order

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | 2024-03 | feature | main   | local, origin |
      | 2025-06 | feature | main   | local, origin |
    And the current branch is "2025-06"

  Scenario: CLI ascending
    When I run "git-town branch --order=asc"
    Then Git Town runs no commands
    And Git Town prints:
      """
        main
          2024-03
      *   2025-06
      """

  Scenario: CLI descending
    When I run "git-town branch --order=desc"
    Then Git Town runs no commands
    And Git Town prints:
      """
        main
      *   2025-06
          2024-03
      """

  Scenario: global Git descending
    Given global Git setting "git-town.order" is "desc"
    When I run "git-town branch"
    Then Git Town runs no commands
    And Git Town prints:
      """
        main
      *   2025-06
          2024-03
      """

  Scenario: local Git descending
    Given local Git setting "git-town.order" is "desc"
    When I run "git-town branch"
    Then Git Town runs no commands
    And Git Town prints:
      """
        main
      *   2025-06
          2024-03
      """

  Scenario: environment variable descending
    When I run "git-town branch" with these environment variables
      | GIT_TOWN_ORDER | desc |
    Then Git Town runs no commands
    And Git Town prints:
      """
        main
      *   2025-06
          2024-03
      """

  Scenario: config file descending
    Given the committed configuration file:
      """
      [branches]
      order = "desc"
      """
    When I run "git-town branch"
    Then Git Town runs no commands
    And Git Town prints:
      """
        main
      *   2025-06
          2024-03
      """
