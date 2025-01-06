Feature: ship the supplied feature branch in a local repo

  Background:
    Given a local Git repo
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
      | other   | feature | main   | local     |
    And the commits
      | BRANCH  | LOCATION | MESSAGE        | FILE NAME        |
      | feature | local    | feature commit | conflicting_file |
    And the current branch is "other"
    And Git setting "git-town.ship-strategy" is "squash-merge"
    When I run "git-town ship feature -m 'feature done'"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                         |
      | other  | git checkout main               |
      | main   | git merge --squash --ff feature |
      |        | git commit -m "feature done"    |
      |        | git checkout other              |
      | other  | git branch -D feature           |
    And the current branch is now "other"
    And the branches are now
      | REPOSITORY | BRANCHES    |
      | local      | main, other |
    And these commits exist now
      | BRANCH | LOCATION | MESSAGE      |
      | main   | local    | feature done |
    And this lineage exists now
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | other  | git checkout main                             |
      | main   | git reset --hard {{ sha 'initial commit' }}   |
      |        | git branch feature {{ sha 'feature commit' }} |
      |        | git checkout other                            |
    And the current branch is now "other"
    And the initial commits exist now
    And the initial branches and lineage exist now
