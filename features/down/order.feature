Feature: change the display order

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE    | PARENT  | LOCATIONS |
      | 2024-03    | feature | main    | local     |
      | 2024-04    | feature | main    | local     |
      | 2025-06    | feature | main    | local     |
      | 2025-06-01 | feature | 2025-06 | local     |
    And the current branch is "2025-06"

  Scenario: CLI ascending
    When I run "git-town down --order=asc"
    Then Git Town runs the commands
      | BRANCH  | COMMAND           |
      | 2025-06 | git checkout main |
    And Git Town prints:
      """
      * main
          2024-03
          2024-04
          2025-06
            2025-06-01
      """

  Scenario: CLI descending
    When I run "git-town down --order=desc"
    Then Git Town runs the commands
      | BRANCH  | COMMAND           |
      | 2025-06 | git checkout main |
    And Git Town prints:
      """
      * main
          2025-06
            2025-06-01
          2024-04
          2024-03
      """

  Scenario: global Git descending
    Given global Git setting "git-town.order" is "desc"
    When I run "git-town down"
    Then Git Town runs the commands
      | BRANCH  | COMMAND           |
      | 2025-06 | git checkout main |
    And Git Town prints:
      """
      * main
          2025-06
            2025-06-01
          2024-04
          2024-03
      """

  Scenario: local Git descending
    Given local Git setting "git-town.order" is "desc"
    When I run "git-town down"
    Then Git Town runs the commands
      | BRANCH  | COMMAND           |
      | 2025-06 | git checkout main |
    And Git Town prints:
      """
      * main
          2025-06
            2025-06-01
          2024-04
          2024-03
      """

  Scenario: environment variable descending
    When I run "git-town down" with these environment variables
      | GIT_TOWN_ORDER | desc |
    Then Git Town runs the commands
      | BRANCH  | COMMAND           |
      | 2025-06 | git checkout main |
    And Git Town prints:
      """
      * main
          2025-06
            2025-06-01
          2024-04
          2024-03
      """

  Scenario: config file descending
    Given the configuration file:
      """
      [branches]
      order = "desc"
      """
    When I run "git-town down"
    Then Git Town runs the commands
      | BRANCH  | COMMAND           |
      | 2025-06 | git checkout main |
    And Git Town prints:
      """
      * main
          2025-06
            2025-06-01
          2024-04
          2024-03
      """
