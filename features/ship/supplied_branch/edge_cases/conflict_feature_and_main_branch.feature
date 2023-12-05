Feature: handle conflicts between the supplied feature branch and the main branch

  Background:
    Given the feature branches "feature" and "other"
    And the commits
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And the current branch is "other"
    And an uncommitted file
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
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And it prints the error:
      """
      To go back to where you started, run "git-town undo".
      To continue after having resolved conflicts, run "git-town continue".
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
    And now these commits exist
      | BRANCH  | LOCATION      | MESSAGE                    |
      | main    | local, origin | conflicting main commit    |
      | feature | local         | conflicting feature commit |
    And the initial branch hierarchy exists

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
    And the current branch is now "other"
    And the uncommitted file still exists
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, other |
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE                 | FILE NAME        | FILE CONTENT     |
      | main   | local, origin | conflicting main commit | conflicting_file | main content     |
      |        |               | feature done            | conflicting_file | resolved content |
    And this branch lineage exists now
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
    And the current branch is now "other"
    And the uncommitted file still exists

  Scenario: resolve, continue, and undo
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue"
    And I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                                       |
      | other  | git add -A                                                    |
      |        | git stash                                                     |
      |        | git checkout main                                             |
      | main   | git revert {{ sha 'feature done' }}                           |
      |        | git push                                                      |
      |        | git push origin {{ sha 'Initial commit' }}:refs/heads/feature |
      |        | git branch feature {{ sha 'conflicting feature commit' }}     |
      |        | git checkout other                                            |
      | other  | git stash pop                                                 |
    And the current branch is now "other"
    And now these commits exist
      | BRANCH  | LOCATION      | MESSAGE                    |
      | main    | local, origin | conflicting main commit    |
      |         |               | feature done               |
      |         |               | Revert "feature done"      |
      | feature | local         | conflicting feature commit |
    And the initial branches and hierarchy exist
