Feature: handle conflicts between the supplied feature branch and the main branch

  Background:
    Given the feature branches "feature" and "other"
    And the commits
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And I am on the "other" branch
    And my workspace has an uncommitted file
    And I run "git-town ship feature -m 'feature done'"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | other   | git fetch --prune --tags           |
      |         | git add -A                         |
      |         | git stash                          |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
    And it prints the error:
      """
      To abort, run "git-town abort".
      To continue after having resolved conflicts, run "git-town continue".
      """
    And I am now on the "feature" branch
    And my uncommitted file is stashed
    And my repo now has a merge in progress

  Scenario: abort
    When I run "git-town abort"
    Then it runs the commands
      | BRANCH  | COMMAND            |
      | feature | git merge --abort  |
      |         | git checkout main  |
      | main    | git checkout other |
      | other   | git stash pop      |
    And I am now on the "other" branch
    And my workspace still contains my uncommitted file
    And there is no merge in progress
    And now these commits exist
      | BRANCH  | LOCATION      | MESSAGE                    |
      | main    | local, origin | conflicting main commit    |
      | feature | local         | conflicting feature commit |
    And Git Town is still aware of the initial branch hierarchy

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND                      |
      | feature | git commit --no-edit         |
      |         | git checkout main            |
      | main    | git merge --squash feature   |
      |         | git commit -m "feature done" |
      |         | git push                     |
      |         | git push origin :feature     |
      |         | git branch -D feature        |
      |         | git checkout other           |
      | other   | git stash pop                |
    And I am now on the "other" branch
    And my workspace still contains my uncommitted file
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, other |
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE                 | FILE NAME        | FILE CONTENT     |
      | main   | local, origin | conflicting main commit | conflicting_file | main content     |
      |        |               | feature done            | conflicting_file | resolved content |
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: resolve, commit, and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git commit --no-edit"
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND                      |
      | feature | git checkout main            |
      | main    | git merge --squash feature   |
      |         | git commit -m "feature done" |
      |         | git push                     |
      |         | git push origin :feature     |
      |         | git branch -D feature        |
      |         | git checkout other           |
      | other   | git stash pop                |
    And I am now on the "other" branch
    And my workspace still contains my uncommitted file

  Scenario: resolve, continue, and undo
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue"
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
      | feature | git reset --hard {{ sha 'conflicting feature commit' }}         |
      |         | git checkout main                                               |
      | main    | git checkout other                                              |
      | other   | git stash pop                                                   |
    And I am now on the "other" branch
    And now these commits exist
      | BRANCH  | LOCATION      | MESSAGE                          |
      | main    | local, origin | conflicting main commit          |
      |         |               | feature done                     |
      |         |               | Revert "feature done"            |
      | feature | local, origin | conflicting feature commit       |
      |         | origin        | conflicting main commit          |
      |         |               | Merge branch 'main' into feature |
    And the initial branch setup and hierarchy exists
