Feature: handle conflicts between the main branch and its tracking branch

  Background:
    Given my repo has the feature branches "feature" and "other"
    And my repo contains the commits
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting local commit  | conflicting_file | local content   |
      |         | remote   | conflicting remote commit | conflicting_file | remote content  |
      | feature | local    | feature commit            | feature_file     | feature content |
    And I am on the "other" branch
    And my workspace has an uncommitted file
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
    And my repo now has a rebase in progress
    And my uncommitted file is stashed

  Scenario: abort
    When I run "git-town abort"
    Then it runs the commands
      | BRANCH | COMMAND            |
      | main   | git rebase --abort |
      |        | git checkout other |
      | other  | git stash pop      |
    And I am still on the "other" branch
    And my workspace still contains my uncommitted file
    And there is no rebase in progress anymore
    And my repo is left with my original commits
    And Git Town still has the original branch hierarchy

  Scenario: resolve the conflict and continue
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
    And I am now on the "other" branch
    And my workspace still contains my uncommitted file
    And my repo now has the commits
      | BRANCH | LOCATION      | MESSAGE                   | FILE NAME        | FILE CONTENT     |
      | main   | local, remote | conflicting remote commit | conflicting_file | remote content   |
      |        |               | conflicting local commit  | conflicting_file | resolved content |
      |        |               | feature done              | feature_file     | feature content  |
    And the existing branches are
      | REPOSITORY    | BRANCHES    |
      | local, remote | main, other |
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: resolve the conflict, finish the rebase, and continue
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
    And I am now on the "other" branch

  Scenario: resolve the conflict, continue, and undo
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
    And I am now on the "other" branch
    And my repo now has the commits
      | BRANCH  | LOCATION      | MESSAGE                          |
      | main    | local, remote | conflicting remote commit        |
      |         |               | conflicting local commit         |
      |         |               | feature done                     |
      |         |               | Revert "feature done"            |
      | feature | local, remote | feature commit                   |
      |         | remote        | conflicting remote commit        |
      |         |               | conflicting local commit         |
      |         |               | Merge branch 'main' into feature |
    And my repo now has its initial branches and branch hierarchy
