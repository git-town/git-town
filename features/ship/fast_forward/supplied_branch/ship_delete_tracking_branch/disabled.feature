Feature: skip deleting the remote branch when shipping another branch using the fast-forward strategy

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
    And the current branch is "other"
    And Git Town setting "ship-delete-tracking-branch" is "false"
    And Git Town setting "ship-strategy" is "fast-forward"
    When I run "git-town ship feature"
    And origin deletes the "feature" branch

  Scenario: result
    Then it runs the commands
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
      | other  | local         | other commit   |
    And this lineage exists now
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                       |
      | other  | git branch feature {{ sha 'feature commit' }} |
    And the current branch is now "other"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE        |
      | main    | local, origin | feature commit |
      | feature | local         | feature commit |
      | other   | local         | other commit   |
    And these branches exist now
      | REPOSITORY | BRANCHES             |
      | local      | main, feature, other |
      | origin     | main, other          |
    And the initial lineage exists now
