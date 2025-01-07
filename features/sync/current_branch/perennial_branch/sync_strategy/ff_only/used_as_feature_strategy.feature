Feature: "ff-only" configured as sync-feature-strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION | MESSAGE  |
      | feature | origin   | commit 1 |
    And the current branch is "feature"
    And Git setting "git-town.sync-perennial-strategy" is "ff-only"
    And Git setting "git-town.sync-feature-strategy" is "ff-only"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git checkout main                  |
      | main    | git checkout feature               |
      | feature | git merge --ff-only origin/feature |
    And the current branch is still "feature"
    And these branches exist now
      | REPOSITORY    | BRANCHES      |
      | local, origin | main, feature |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE  |
      | feature | local, origin | commit 1 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                     |
      | feature | git reset --hard {{ sha 'initial commit' }} |
    And the current branch is still "feature"
    And the initial commits exist now
    And the initial branches and lineage exist now
