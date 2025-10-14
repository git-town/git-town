Feature: skip deleting the remote branch when shipping another branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
      | other   | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
      | other   | local         | other commit   |
    And Git setting "git-town.ship-delete-tracking-branch" is "false"
    And Git setting "git-town.ship-strategy" is "squash-merge"
    And the current branch is "other"
    When I run "git-town ship feature -m 'feature done'"
    And origin deletes the "feature" branch

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                         |
      | other  | git fetch --prune --tags        |
      |        | git checkout main               |
      | main   | git merge --squash --ff feature |
      |        | git commit -m "feature done"    |
      |        | git push                        |
      |        | git checkout other              |
      | other  | git branch -D feature           |
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, other |
    And this lineage exists now
      """
      main
        other
      """
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      |
      | main   | local, origin | feature done |
      | other  | local         | other commit |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | other  | git checkout main                             |
      | main   | git revert {{ sha 'feature done' }}           |
      |        | git push                                      |
      |        | git branch feature {{ sha 'feature commit' }} |
      |        | git checkout other                            |
    And the branches are now
      | REPOSITORY | BRANCHES             |
      | local      | main, feature, other |
      | origin     | main, other          |
    And the initial lineage exists now
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE               |
      | main    | local, origin | feature done          |
      |         |               | Revert "feature done" |
      | feature | local         | feature commit        |
      | other   | local         | other commit          |
