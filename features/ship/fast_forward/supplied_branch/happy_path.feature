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
    And an uncommitted file with name "conflicting_file" and content "conflicting content"
    And Git Town setting "ship-strategy" is "fast-forward"
    When I run "git-town ship feature"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                     |
      | other  | git fetch --prune --tags    |
      |        | git add -A                  |
      |        | git stash                   |
      |        | git checkout main           |
      | main   | git merge --ff-only feature |
      |        | git push                    |
      |        | git push origin :feature    |
      |        | git checkout other          |
      | other  | git branch -D feature       |
      |        | git stash pop               |
    And the current branch is now "other"
    And the uncommitted file still exists
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
    Then it runs the commands
      | BRANCH | COMMAND                                       |
      | other  | git add -A                                    |
      |        | git stash                                     |
      |        | git branch feature {{ sha 'feature commit' }} |
      |        | git push -u origin feature                    |
      |        | git stash pop                                 |
    And the current branch is now "other"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE        |
      | main   | local, origin | feature commit |
    And the initial branches and lineage exist now
