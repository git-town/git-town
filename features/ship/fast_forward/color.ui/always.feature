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
    And local Git setting "color.ui" is "always"
    And Git setting "git-town.ship-strategy" is "fast-forward"
    And the current branch is "other"
    When I run "git-town ship feature"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      | other  | git fetch --prune --tags    |
      |        | git checkout main           |
      | main   | git merge --ff-only feature |
      |        | git push                    |
      |        | git push origin :feature    |
      |        | git checkout other          |
      | other  | git branch -D feature       |
    And this lineage exists now
      """
      main
        other
      """
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, other |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE        |
      | main   | local, origin | feature commit |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | other  | git branch feature {{ sha 'feature commit' }} |
      |        | git push -u origin feature                    |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE        |
      | main   | local, origin | feature commit |
