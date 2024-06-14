Feature: handle conflicts between the supplied feature branch and its tracking branch

  Background:
    Given the feature branches "feature" and "other"
    And the commits
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | feature | local    | conflicting local commit  | conflicting_file | local content  |
      |         | origin   | conflicting origin commit | conflicting_file | origin content |
    And the current branch is "other"
    And an uncommitted file
    And I run "git-town ship feature -m 'feature done'"

  @this
  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                                 |
      | other   | git fetch --prune --tags                |
      |         | git add -A                              |
      |         | git stash                               |
      |         | git checkout main                       |
      | main    | git rebase origin/main                  |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff origin/feature |
    And it prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And it prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      To go back to where you started, run "git town undo".
      """
    And the current branch is now "feature"
    And the uncommitted file is stashed
    And a merge is now in progress

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND            |
      | feature | git merge --abort  |
      |         | git checkout other |
      | other   | git stash pop      |
    And the current branch is now "other"
    And the uncommitted file still exists
    And no merge is in progress
    And the initial commits exist
    And the initial lineage exists

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND                         |
      | feature | git commit --no-edit            |
      |         | git merge --no-edit --ff main   |
      |         | git checkout main               |
      | main    | git merge --squash --ff feature |
      |         | git commit -m "feature done"    |
      |         | git push                        |
      |         | git push origin :feature        |
      |         | git branch -D feature           |
      |         | git checkout other              |
      | other   | git stash pop                   |
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

  Scenario: resolve, commit, and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git commit --no-edit"
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND                         |
      | feature | git merge --no-edit --ff main   |
      |         | git checkout main               |
      | main    | git merge --squash --ff feature |
      |         | git commit -m "feature done"    |
      |         | git push                        |
      |         | git push origin :feature        |
      |         | git branch -D feature           |
      |         | git checkout other              |
      | other   | git stash pop                   |
    And the current branch is now "other"
    And the uncommitted file still exists

  Scenario: resolve, continue, and undo
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue"
    And I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                                                            |
      | other  | git add -A                                                                         |
      |        | git stash                                                                          |
      |        | git checkout main                                                                  |
      | main   | git revert {{ sha 'feature done' }}                                                |
      |        | git push                                                                           |
      |        | git push origin {{ sha-in-origin 'conflicting origin commit' }}:refs/heads/feature |
      |        | git branch feature {{ sha 'conflicting local commit' }}                            |
      |        | git checkout other                                                                 |
      | other  | git stash pop                                                                      |
    And the current branch is now "other"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                   |
      | main    | local, origin | feature done              |
      |         |               | Revert "feature done"     |
      | feature | local         | conflicting local commit  |
      |         | origin        | conflicting origin commit |
    And the initial branches and lineage exist
