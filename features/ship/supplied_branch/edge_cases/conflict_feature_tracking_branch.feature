Feature: handle conflicts between the supplied feature branch and its tracking branch

  Background:
    Given my repo has the feature branches "feature" and "other"
    And the commits
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | feature | local    | conflicting local commit  | conflicting_file | local content  |
      |         | remote   | conflicting remote commit | conflicting_file | remote content |
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
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
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
    And now the initial commits exist
    And Git Town is still aware of the initial branch hierarchy

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND                      |
      | feature | git commit --no-edit         |
      |         | git merge --no-edit main     |
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
    And the existing branches are
      | REPOSITORY    | BRANCHES    |
      | local, remote | main, other |
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE      |
      | main   | local, remote | feature done |
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: resolve, commit, and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git commit --no-edit"
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND                      |
      | feature | git merge --no-edit main     |
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

  Scenario: resolve, continue, and undo
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue"
    And I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                                                                   |
      | other   | git add -A                                                                                |
      |         | git stash                                                                                 |
      |         | git checkout main                                                                         |
      | main    | git branch feature {{ sha 'Merge remote-tracking branch 'origin/feature' into feature' }} |
      |         | git push -u origin feature                                                                |
      |         | git revert {{ sha 'feature done' }}                                                       |
      |         | git push                                                                                  |
      |         | git checkout feature                                                                      |
      | feature | git checkout main                                                                         |
      | main    | git checkout other                                                                        |
      | other   | git stash pop                                                                             |
    And I am now on the "other" branch
    And now these commits exist
      | BRANCH  | LOCATION      | MESSAGE                                                    |
      | main    | local, remote | feature done                                               |
      |         |               | Revert "feature done"                                      |
      | feature | local, remote | conflicting local commit                                   |
      |         |               | conflicting remote commit                                  |
      |         |               | Merge remote-tracking branch 'origin/feature' into feature |
    And my repo now has its initial branches and branch hierarchy
