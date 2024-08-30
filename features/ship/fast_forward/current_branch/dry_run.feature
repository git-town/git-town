Feature: dry-run shipping via the fast-forward strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And the current branch is "feature"
    And Git Town setting "ship-strategy" is "fast-forward"
    When I run "git-town ship --dry-run"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                     |
      | feature | git fetch --prune --tags    |
      |         | git checkout main           |
      | main    | git merge --ff-only feature |
      |         | git push origin :feature    |
      |         | git branch -D feature       |
    And the current branch is still "feature"
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "feature"
    And the initial commits exist
    And the initial branches and lineage exist
