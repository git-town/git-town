Feature: ship an omni-branch via the fast-forward strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And Git setting "git-town.ship-strategy" is "fast-forward"
    And the current branch is "feature"
    When I run "git-town ship"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                     |
      | feature | git fetch --prune --tags    |
      |         | git checkout main           |
      | main    | git merge --ff-only feature |
      |         | git push                    |
      |         | git push origin :feature    |
      |         | git branch -D feature       |
    And no lineage exists now
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE        |
      | main   | local, origin | feature commit |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                               |
      | main   | git branch feature {{ sha-initial 'feature commit' }} |
      |        | git push -u origin feature                            |
      |        | git checkout feature                                  |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE        |
      | main   | local, origin | feature commit |
