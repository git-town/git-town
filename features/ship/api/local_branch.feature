Feature: cannot ship a local branch via API

  Background:
    Given a Git repo with origin
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | local    | feature commit |
    And the current branch is "feature"
    And Git Town setting "ship-strategy" is "api"
    And the origin is "git@github.com:git-town/git-town.git"
    And a proposal for this branch does not exist
    When I run "git-town ship -m done"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
    And it prints the error:
      """
      cannot ship branch "feature" via API because it has no remote branch
      """
    And the initial branches and lineage exist
    And the initial commits exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the initial commits exist
    And the initial branches and lineage exist
