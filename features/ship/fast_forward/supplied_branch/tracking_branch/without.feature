Feature: ship the supplied local feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local         |
      | other   | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION | MESSAGE        | FILE NAME        |
      | feature | local    | feature commit | conflicting_file |
    And the current branch is "other"
    And Git setting "git-town.ship-strategy" is "fast-forward"
    When I run "git-town ship feature"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      | other  | git fetch --prune --tags    |
      |        | git checkout main           |
      | main   | git merge --ff-only feature |
      |        | git push                    |
      |        | git checkout other          |
      | other  | git branch -D feature       |
    And the current branch is now "other"
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, other |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE        |
      | main   | local, origin | feature commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | other  | git branch feature {{ sha 'feature commit' }} |
    And the current branch is now "other"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE        |
      | main   | local, origin | feature commit |
    And the initial branches and lineage exist now
