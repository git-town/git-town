Feature: ship the supplied feature branch in a local repo using the always-merge strategy

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
    And an uncommitted file with name "conflicting_file" and content "conflicting content"
    And Git setting "git-town.ship-strategy" is "always-merge"
    When I run "git-town ship feature" and close the editor

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                             |
      | other  | git add -A                          |
      |        | git stash                           |
      |        | git checkout main                   |
      | main   | git merge --no-ff --edit -- feature |
      |        | git checkout other                  |
      | other  | git branch -D feature               |
      |        | git stash pop                       |
    And the current branch is now "other"
    And the uncommitted file still exists
    And the branches are now
      | REPOSITORY | BRANCHES    |
      | local      | main, other |
    And these commits exist now
      | BRANCH | LOCATION | MESSAGE                |
      | main   | local    | feature commit         |
      |        |          | Merge branch 'feature' |
    And this lineage exists now
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | other  | git add -A                                    |
      |        | git stash                                     |
      |        | git checkout main                             |
      | main   | git reset --hard {{ sha 'initial commit' }}   |
      |        | git branch feature {{ sha 'feature commit' }} |
      |        | git checkout other                            |
      | other  | git stash pop                                 |
    And the current branch is now "other"
    And the initial commits exist now
    And the initial branches and lineage exist now