Feature: does not compress perennial branches

  Scenario: on perennial branch
    Given a Git repo clone
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | perennial | perennial |        | local, origin |
    And the current branch is "perennial"
    And the commits
      | BRANCH    | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | perennial | local, origin | commit 1 | file_1    | content 1    |
      |           |               | commit 2 | file_2    | content 2    |
    When I run "git-town compress"
    Then it runs the commands
      | BRANCH    | COMMAND                  |
      | perennial | git fetch --prune --tags |
    And it prints the error:
      """
      better not compress perennial branches
      """
    And the current branch is still "perennial"
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: on main branch
    Given a Git repo clone
    Given the commits
      | BRANCH | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | main   | local, origin | commit 1 | file_1    | content 1    |
      |        |               | commit 2 | file_2    | content 2    |
    When I run "git-town compress"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      better not compress perennial branches
      """
    And the current branch is still "main"
    And the initial commits exist
    And the initial branches and lineage exist
