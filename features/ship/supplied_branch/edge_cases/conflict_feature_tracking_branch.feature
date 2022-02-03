Feature: handle conflicts between the supplied feature branch and its tracking branch

  Background:
    Given my repo has the feature branches "feature" and "other-feature"
    And my repo contains the commits
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | feature | local    | local conflicting commit  | conflicting_file | local conflicting content  |
      |         | remote   | remote conflicting commit | conflicting_file | remote conflicting content |
    And I am on the "other-feature" branch
    And my workspace has an uncommitted file
    And I run "git-town ship feature -m 'feature done'"

  Scenario: result
    Then it runs the commands
      | BRANCH        | COMMAND                            |
      | other-feature | git fetch --prune --tags           |
      |               | git add -A                         |
      |               | git stash                          |
      |               | git checkout main                  |
      | main          | git rebase origin/main             |
      |               | git checkout feature               |
      | feature       | git merge --no-edit origin/feature |
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
      | BRANCH        | COMMAND                    |
      | feature       | git merge --abort          |
      |               | git checkout main          |
      | main          | git checkout other-feature |
      | other-feature | git stash pop              |
    And I am now on the "other-feature" branch
    And my workspace still contains my uncommitted file
    And there is no merge in progress
    And my repo is left with my original commits
    And Git Town still has the original branch hierarchy

  Scenario: continue after resolving the conflicts
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH        | COMMAND                      |
      | feature       | git commit --no-edit         |
      |               | git merge --no-edit main     |
      |               | git checkout main            |
      | main          | git merge --squash feature   |
      |               | git commit -m "feature done" |
      |               | git push                     |
      |               | git push origin :feature     |
      |               | git branch -D feature        |
      |               | git checkout other-feature   |
      | other-feature | git stash pop                |
    And I am now on the "other-feature" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY    | BRANCHES            |
      | local, remote | main, other-feature |
    And my repo now has the following commits
      | BRANCH | LOCATION      | MESSAGE      |
      | main   | local, remote | feature done |
    And Git Town is now aware of this branch hierarchy
      | BRANCH        | PARENT |
      | other-feature | main   |

  Scenario: continue after resolving the conflicts and comitting
    When I resolve the conflict in "conflicting_file"
    And I run "git commit --no-edit"
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH        | COMMAND                      |
      | feature       | git merge --no-edit main     |
      |               | git checkout main            |
      | main          | git merge --squash feature   |
      |               | git commit -m "feature done" |
      |               | git push                     |
      |               | git push origin :feature     |
      |               | git branch -D feature        |
      |               | git checkout other-feature   |
      | other-feature | git stash pop                |
    And I am now on the "other-feature" branch
    And my workspace still contains my uncommitted file

  Scenario: undo after continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue"
    And I run "git-town undo"
    Then it runs the commands
      | BRANCH        | COMMAND                                                                                   |
      | other-feature | git add -A                                                                                |
      |               | git stash                                                                                 |
      |               | git checkout main                                                                         |
      | main          | git branch feature {{ sha 'Merge remote-tracking branch 'origin/feature' into feature' }} |
      |               | git push -u origin feature                                                                |
      |               | git revert {{ sha 'feature done' }}                                                       |
      |               | git push                                                                                  |
      |               | git checkout feature                                                                      |
      | feature       | git checkout main                                                                         |
      | main          | git checkout other-feature                                                                |
      | other-feature | git stash pop                                                                             |
    And I am now on the "other-feature" branch
    And my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE                                                    |
      | main    | local, remote | feature done                                               |
      |         |               | Revert "feature done"                                      |
      | feature | local, remote | local conflicting commit                                   |
      |         |               | remote conflicting commit                                  |
      |         |               | Merge remote-tracking branch 'origin/feature' into feature |
    And my repo now has the original branches
    And Git Town now has the original branch hierarchy
