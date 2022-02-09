Feature: handle conflicts between the main branch and its tracking branch

  Background:
    Given the feature branches "feature" and "other"
    And the commits
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting local commit  | conflicting_file | local content   |
      |         | origin   | conflicting origin commit | conflicting_file | origin content  |
      | feature | local    | feature commit            | feature_file     | feature content |
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
      To abort, run "git-town abort".
      To continue after having resolved conflicts, run "git-town continue".
      """
    And a rebase is now in progress
    And the uncommitted file is stashed

  Scenario: abort
    When I run "git-town abort"
    Then it runs the commands
      | BRANCH | COMMAND            |
      | main   | git rebase --abort |
      |        | git checkout other |
      | other  | git stash pop      |
    And the current branch is still "other"
    And the uncommitted file still exists
    And no rebase is in progress
    And now the initial commits exist
    And the initial branch hierarchy exists

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and close the editor
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | main    | git rebase --continue              |
      |         | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git checkout main                  |
      | main    | git merge --squash feature         |
      |         | git commit -m "feature done"       |
      |         | git push                           |
      |         | git push origin :feature           |
      |         | git branch -D feature              |
      |         | git checkout other                 |
      | other   | git stash pop                      |
    And the current branch is now "other"
    And the uncommitted file still exists
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE                   | FILE NAME        | FILE CONTENT     |
      | main   | local, origin | conflicting origin commit | conflicting_file | origin content   |
      |        |               | conflicting local commit  | conflicting_file | resolved content |
      |        |               | feature done              | feature_file     | feature content  |
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, other |
    And this branch hierarchy exists now
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: resolve, finish the rebase, and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git rebase --continue" and close the editor
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | main    | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git checkout main                  |
      | main    | git merge --squash feature         |
      |         | git commit -m "feature done"       |
      |         | git push                           |
      |         | git push origin :feature           |
      |         | git branch -D feature              |
      |         | git checkout other                 |
      | other   | git stash pop                      |
    And the current branch is now "other"

  Scenario: resolve, continue, and undo
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and close the editor
    And I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                                         |
      | other   | git add -A                                                      |
      |         | git stash                                                       |
      |         | git checkout main                                               |
      | main    | git branch feature {{ sha 'Merge branch 'main' into feature' }} |
      |         | git push -u origin feature                                      |
      |         | git revert {{ sha 'feature done' }}                             |
      |         | git push                                                        |
      |         | git checkout feature                                            |
      | feature | git reset --hard {{ sha 'feature commit' }}                     |
      |         | git checkout main                                               |
      | main    | git checkout other                                              |
      | other   | git stash pop                                                   |
    And the current branch is now "other"
    And now these commits exist
      | BRANCH  | LOCATION      | MESSAGE                          |
      | main    | local, origin | conflicting origin commit        |
      |         |               | conflicting local commit         |
      |         |               | feature done                     |
      |         |               | Revert "feature done"            |
      | feature | local, origin | feature commit                   |
      |         | origin        | conflicting origin commit        |
      |         |               | conflicting local commit         |
      |         |               | Merge branch 'main' into feature |
    And the initial branches and hierarchy exist
