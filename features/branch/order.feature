Feature: change the display order

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE    | PARENT  | LOCATIONS |
      | 2024-03    | feature | main    | local     |
      | 2025-03-01 | feature | 2024-03 | local     |
      | 2025-03-02 | feature | 2024-03 | local     |
      | 2025-06    | feature | main    | local     |
      | 2025-06-01 | feature | 2025-06 | local     |
      | 2025-06-02 | feature | 2025-06 | local     |
      | 2025-06-03 | feature | 2025-06 | local     |
    And the current branch is "2025-06"

  Scenario: CLI ascending
    When I run "git-town branch --order=asc"
    Then Git Town runs no commands
    And Git Town prints:
      """
        main
          2024-03
            2025-03-01
            2025-03-02
      *   2025-06
            2025-06-01
            2025-06-02
            2025-06-03
      """

  Scenario: CLI descending
    When I run "git-town branch --order=desc"
    Then Git Town runs no commands
    And Git Town prints:
      """
        main
      *   2025-06
            2025-06-03
            2025-06-02
            2025-06-01
          2024-03
            2025-03-02
            2025-03-01
      """

  Scenario: global Git descending
    Given global Git setting "git-town.order" is "desc"
    When I run "git-town branch"
    Then Git Town runs no commands
    And Git Town prints:
      """
        main
      *   2025-06
            2025-06-03
            2025-06-02
            2025-06-01
          2024-03
            2025-03-02
            2025-03-01
      """

  Scenario: local Git descending
    Given local Git setting "git-town.order" is "desc"
    When I run "git-town branch"
    Then Git Town runs no commands
    And Git Town prints:
      """
        main
      *   2025-06
            2025-06-03
            2025-06-02
            2025-06-01
          2024-03
            2025-03-02
            2025-03-01
      """

  Scenario: environment variable descending
    When I run "git-town branch" with these environment variables
      | GIT_TOWN_ORDER | desc |
    Then Git Town runs no commands
    And Git Town prints:
      """
        main
      *   2025-06
            2025-06-03
            2025-06-02
            2025-06-01
          2024-03
            2025-03-02
            2025-03-01
      """

  Scenario: config file descending
    Given the configuration file:
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
            2025-06-03
            2025-06-02
            2025-06-01
          2024-03
            2025-03-02
            2025-03-01
      """
