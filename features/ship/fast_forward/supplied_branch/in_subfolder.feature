Feature: ship the supplied feature branch from a subfolder using the fast-forward strategy

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
    And an uncommitted file with name "new_folder/other_feature_file" and content "other feature content"
    And Git Town setting "ship-strategy" is "fast-forward"
    When I run "git-town ship feature" in the "new_folder" folder

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                     |
      | other  | git fetch --prune --tags    |
      |        | git add -A                  |
      |        | git stash                   |
      |        | git checkout main           |
      | main   | git merge --ff-only feature |
      |        | git push                    |
      |        | git checkout other          |
      | other  | git branch -D feature       |
      |        | git stash pop               |
    And the current branch is now "other"
    And the uncommitted file still exists
    And the branches are now
      | REPOSITORY | BRANCHES    |
      | local      | main, other |
      | origin     | main        |
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
      |        | git stash pop                                 |
    And the current branch is now "other"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE        |
      | main   | local, origin | feature commit |
    And the initial branches and lineage exist now
