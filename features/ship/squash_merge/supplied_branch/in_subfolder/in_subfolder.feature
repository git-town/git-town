Feature: ship the supplied feature branch from a subfolder

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
      | other   | feature | main   | local     |
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | local    | feature commit |
    And the current branch is "other"
    And Git setting "git-town.ship-strategy" is "squash-merge"
    And a folder "new_folder"
    When I run "git-town ship feature -m 'feature done'" in the "new_folder" folder

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
      | REPOSITORY | BRANCHES    |
      | local      | main, other |
      | origin     | main        |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      |
      | main   | local, origin | feature done |
    And this lineage exists now
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | other  | git checkout main                             |
      | main   | git revert {{ sha 'feature done' }}           |
      |        | git push                                      |
      |        | git branch feature {{ sha 'feature commit' }} |
      |        | git checkout other                            |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE               |
      | main    | local, origin | feature done          |
      |         |               | Revert "feature done" |
      | feature | local         | feature commit        |
    And the initial branches and lineage exist now
