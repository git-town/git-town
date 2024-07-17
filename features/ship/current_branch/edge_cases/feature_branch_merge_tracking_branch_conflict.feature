Feature: handle conflicts between the shipped branch and its tracking branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    Given the current branch is "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | feature | local    | conflicting local commit  | conflicting_file | local content  |
      |         | origin   | conflicting origin commit | conflicting_file | origin content |
    And Git Town setting "sync-before-ship" is "true"
    When I run "git-town ship -m 'feature done'"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
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
    And the current branch is still "feature"
    And a merge is now in progress

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND           |
      | feature | git merge --abort |
    And the current branch is still "feature"
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
    And the current branch is now "main"
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME        |
      | main   | local, origin | feature done | conflicting_file |
    And no lineage exists now

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
    And the current branch is now "main"
