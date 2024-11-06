Feature: does not delete perennial branches

  Background:
    Given a Git repo with origin

  Scenario: main branch
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    Given the current branch is "feature"
    When I run "git-town delete main"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
    And it prints the error:
      """
      you cannot delete the main branch
      """
    And the current branch is still "feature"
    And the initial branches and lineage exist now

  Scenario: perennial branch
    And the branches
      | NAME | TYPE      | LOCATIONS     |
      | qa   | perennial | local, origin |
    And the current branch is "main"
    When I run "git-town delete qa"
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      you cannot delete perennial branches
      """
