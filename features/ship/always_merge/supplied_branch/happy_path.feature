Feature: ship the supplied feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
      | other   | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        | FILE NAME        |
      | feature | local, origin | feature commit | conflicting_file |
    And the current branch is "other"
    And Git setting "git-town.ship-strategy" is "always-merge"
    When I run "git-town ship feature" and close the editor

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                             |
      | other  | git fetch --prune --tags            |
      |        | git checkout main                   |
      | main   | git merge --no-ff --edit -- feature |
      |        | git push                            |
      |        | git push origin :feature            |
      |        | git checkout other                  |
      | other  | git branch -D feature               |
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, other |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                |
      | main   | local, origin | feature commit         |
      |        |               | Merge branch 'feature' |
    And this lineage exists now
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | other  | git branch feature {{ sha 'feature commit' }} |
      |        | git push -u origin feature                    |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                |
      | main   | local, origin | feature commit         |
      |        |               | Merge branch 'feature' |
    And the initial branches and lineage exist now
