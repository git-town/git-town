Feature: handle conflicts between the main branch and its tracking branch

  Background:
    Given the feature branches "feature" and "other"
    And the commits
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting local commit  | conflicting_file | local content   |
      |         | origin   | conflicting origin commit | conflicting_file | origin content  |
      | feature | local    | feature commit            | feature_file     | feature content |
    And Git Town setting "sync-before-ship" is "true"
    And the current branch is "other"
    And an uncommitted file
    And I run "git-town ship feature -m 'feature done'"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | other  | git fetch --prune --tags |
      |        | git add -A               |
      |        | git stash                |
      |        | git checkout main        |
      | main   | git rebase origin/main   |
    And it prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And it prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      To go back to where you started, run "git town undo".
      """
    And a rebase is now in progress
    And the uncommitted file is stashed

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND            |
      | main   | git rebase --abort |
      |        | git checkout other |
      | other  | git stash pop      |
    And the current branch is still "other"
    And the uncommitted file still exists
    And no rebase is in progress
    And the initial commits exist
    And the initial lineage exists

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and close the editor
    Then it runs the commands
      | BRANCH  | COMMAND                                 |
      | main    | git rebase --continue                   |
      |         | git push                                |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff origin/feature |
      |         | git merge --no-edit --ff main           |
      |         | git checkout main                       |
      | main    | git merge --squash --ff feature         |
      |         | git commit -m "feature done"            |
      |         | git push                                |
      |         | git push origin :feature                |
      |         | git branch -D feature                   |
      |         | git checkout other                      |
      | other   | git stash pop                           |
    And the current branch is now "other"
    And the uncommitted file still exists
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                   | FILE NAME        | FILE CONTENT     |
      | main   | local, origin | conflicting origin commit | conflicting_file | origin content   |
      |        |               | conflicting local commit  | conflicting_file | resolved content |
      |        |               | feature done              | feature_file     | feature content  |
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, other |
    And this lineage exists now
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: resolve, finish the rebase, and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git rebase --continue" and close the editor
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND                                 |
      | main    | git push                                |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff origin/feature |
      |         | git merge --no-edit --ff main           |
      |         | git checkout main                       |
      | main    | git merge --squash --ff feature         |
      |         | git commit -m "feature done"            |
      |         | git push                                |
      |         | git push origin :feature                |
      |         | git branch -D feature                   |
      |         | git checkout other                      |
      | other   | git stash pop                           |
    And the current branch is now "other"

  Scenario: resolve, continue, and undo
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and close the editor
    And I run "git-town undo"
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
      | BRANCH  | LOCATION      | MESSAGE                   |
      | main    | local, origin | conflicting origin commit |
      |         |               | conflicting local commit  |
      |         |               | feature done              |
      |         |               | Revert "feature done"     |
      | feature | local         | feature commit            |
    And the initial branches and lineage exist
