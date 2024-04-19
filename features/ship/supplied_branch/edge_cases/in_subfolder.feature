Feature: ship the supplied feature branch from a subfolder

  Background:
    Given the feature branches "feature" and "other"
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | local    | feature commit |
    And the current branch is "other"
    And an uncommitted file with name "new_folder/other_feature_file" and content "other feature content"
    When I run "git-town ship feature -m 'feature done'" in the "new_folder" folder

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                         |
      | other  | git fetch --prune --tags        |
      |        | git add -A                      |
      |        | git stash                       |
      |        | git checkout main               |
      | main   | git merge --squash --ff feature |
      |        | git commit -m "feature done"    |
      |        | git push                        |
      |        | git push origin :feature        |
      |        | git branch -D feature           |
      |        | git checkout other              |
      | other  | git stash pop                   |
    And the current branch is now "other"
    And the uncommitted file still exists
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, other |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      |
      | main   | local, origin | feature done |
    And this lineage exists now
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                                       |
      | other  | git add -A                                                    |
      |        | git stash                                                     |
      |        | git checkout main                                             |
      | main   | git revert {{ sha 'feature done' }}                           |
      |        | git push                                                      |
      |        | git push origin {{ sha 'initial commit' }}:refs/heads/feature |
      |        | git branch feature {{ sha 'feature commit' }}                 |
      |        | git checkout other                                            |
      | other  | git stash pop                                                 |
    And the current branch is now "other"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE               |
      | main    | local, origin | feature done          |
      |         |               | Revert "feature done" |
      | feature | local         | feature commit        |
    And the initial branches and lineage exist
